package ibkr_test

import (
	"context"
	"flag"
	"fmt"
	"strings"
	"testing"
	"time"
	_ "time/tzdata"

	"github.com/mxmauro/ibkr"
	"github.com/mxmauro/ibkr/models"
)

// -----------------------------------------------------------------------------

var liveTrading = flag.Bool("live", false, "run tests against the live IB Gateway on port 4001 instead of paper trading on port 4002")

func TestSuite(t *testing.T) {
	client := connect(t)
	t.Cleanup(func() {
		client.Destroy()
	})

	sw := newStopWatch()

	t.Run("Current-time", func(t *testing.T) {
		swm := newStopWatchMeasure(t, sw)
		defer swm.End()

		testCurrentTime(t, client)
	})
	t.Run("Managed-accounts", func(t *testing.T) {
		t.Parallel()

		swm := newStopWatchMeasure(t, sw)
		defer swm.End()

		testManagedAccounts(t, client)
	})
	t.Run("Historical-data", func(t *testing.T) {
		t.Parallel()

		swm := newStopWatchMeasure(t, sw)
		defer swm.End()

		testHistoricalData(t, client)
	})
	t.Run("Historical-ticks", func(t *testing.T) {
		t.Parallel()

		swm := newStopWatchMeasure(t, sw)
		defer swm.End()

		testHistoricalTicks(t, client)
	})
	t.Run("Contract-details", func(t *testing.T) {
		t.Parallel()

		swm := newStopWatchMeasure(t, sw)
		defer swm.End()

		testContractDetails(t, client)
	})
	t.Run("Matching-symbols", func(t *testing.T) {
		t.Parallel()

		swm := newStopWatchMeasure(t, sw)
		defer swm.End()

		testMatchingSymbols(t, client)
	})
	t.Run("Top-market-data", func(t *testing.T) {
		t.Parallel()

		swm := newStopWatchMeasure(t, sw)
		defer swm.End()

		testTopMarketData(t, client)
	})
	t.Run("Depth-market-data", func(t *testing.T) {
		t.Parallel()

		swm := newStopWatchMeasure(t, sw)
		defer swm.End()

		testDepthMarketData(t, client)
	})
	t.Run("Real-time-bars", func(t *testing.T) {
		swm := newStopWatchMeasure(t, sw)
		defer swm.End()

		testRealTimeBars(t, client)
	})
	t.Run("Tick-by-tick", func(t *testing.T) {
		swm := newStopWatchMeasure(t, sw)
		defer swm.End()

		testTickByTickData(t, client)
	})
	t.Run("Historical-data-keep-up-to-date", func(t *testing.T) {
		swm := newStopWatchMeasure(t, sw)
		defer swm.End()

		testHistoricalDataKeepUpToDate(t, client)
	})
	t.Run("Security-definition-option-parameters", func(t *testing.T) {
		swm := newStopWatchMeasure(t, sw)
		defer swm.End()

		testSecurityDefinitionOptionParameters(t, client)
	})
	t.Run("Paper-orders", func(t *testing.T) {
		swm := newStopWatchMeasure(t, sw)
		defer swm.End()

		testPaperOrders(t)
	})
}

// -----------------------------------------------------------------------------

func connect(t *testing.T) *ibkr.Client {
	return connectWithOrderStore(t, nil)
}

func connectWithOrderStore(t *testing.T, orderStore ibkr.OrderStore) *ibkr.Client {
	address := getTestServerAddress()
	testLogf(t, "Connecting to server (%s)...", address)
	if *liveTrading {
		testLog(t, "Live trading mode enabled")
	} else {
		testLog(t, "Paper trading mode enabled")
	}

	ctx, cancelCtx := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelCtx()

	client, err := ibkr.NewClient(ctx, ibkr.Options{
		Address:    address,
		OrderStore: orderStore,
		EventsHandler: ibkr.NewIncomingMessageLogger(func(msg string) {
			testEventLog(t, msg)
		}),
	})
	if err != nil {
		testFatal(t, err)
	}
	testLog(t, "Server version is:", client.ServerVersion())

	// Done
	return client
}

// -----------------------------------------------------------------------------

func getTestServerAddress() string {
	port := 4002
	if *liveTrading {
		port = 4001
	}
	return fmt.Sprintf("127.0.0.1:%d", port)
}

// -----------------------------------------------------------------------------

func skipIfLiveTrading(t *testing.T, reason string) {
	t.Helper()

	if *liveTrading {
		testSkip(t, "Skipping in live trading mode: "+reason)
	}
}

// -----------------------------------------------------------------------------

func testCurrentTime(t *testing.T, client *ibkr.Client) {
	ctx, cancelCtx := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelCtx()

	ts, err := client.RequestCurrentTime(ctx)
	if err != nil {
		testError(t, err)
		return
	}
	testLog(t, "  Time is: "+ts.String())
}

func testManagedAccounts(t *testing.T, client *ibkr.Client) {
	ctx, cancelCtx := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelCtx()

	accounts, err := client.RequestManagedAccounts(ctx)
	if err != nil {
		testError(t, err)
		return
	}
	testLog(t, "  Accounts: "+strings.Join(accounts, ", "))
}

func testHistoricalData(t *testing.T, client *ibkr.Client) {
	ctx, cancelCtx := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancelCtx()

	histData, err := client.RequestHistoricalData(ctx, models.HistoricalDataRequestOptions{
		Contract:                getContract("JPM", "SMART"),
		Duration:                1,
		DurationUnit:            models.DurationUnitDays,
		EndDate:                 time.Date(2021, time.May, 12, 0, 0, 0, 0, time.UTC),
		BarSize:                 models.BarSizeOneMinute,
		WhatToShow:              models.WhatToShowTrades,
		OnlyRegularTradingHours: true,
	})
	if err != nil {
		testError(t, err)
		return
	}
	for _, bar := range histData.Bars {
		testLog(t, "  "+bar.String())
	}
}

func testHistoricalTicks(t *testing.T, client *ibkr.Client) {
	ctx, cancelCtx := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelCtx()

	histTicks, err := client.RequestHistoricalTicks(ctx, models.HistoricalTicksRequestOptions{
		Contract:                getContract("JPM", "SMART"),
		StartDate:               time.Date(2025, time.May, 1, 0, 0, 0, 0, time.UTC),
		EndDate:                 time.Date(2025, time.May, 2, 0, 0, 0, 0, time.UTC),
		NumberOfTicks:           1000,
		WhatToShow:              models.WhatToShowTrades,
		OnlyRegularTradingHours: true,
		IgnoreSize:              false,
	})
	if err != nil {
		testError(t, err)
		return
	}
	for _, tick := range histTicks.TicksLast {
		testLog(t, "  "+tick.String())
	}
}

func testContractDetails(t *testing.T, client *ibkr.Client) {
	ctx, cancelCtx := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelCtx()

	resp, err := client.RequestContractDetails(ctx, models.ContractDetailsRequestOptions{
		Contract: getContract("JPM", "SMART"),
	})
	if err != nil {
		testError(t, err)
		return
	}
	for _, cd := range resp.ContractDetails {
		testLog(t, "  "+cd.String())
	}
}

func testMatchingSymbols(t *testing.T, client *ibkr.Client) {
	ctx, cancelCtx := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelCtx()

	resp, err := client.RequestMatchingSymbols(ctx, models.MatchingSymbolsRequestOptions{
		Pattern: "AAPL",
	})
	if err != nil {
		testError(t, err)
		return
	}
	for _, cd := range resp.ContractDescriptions {
		testLog(t, "  "+cd.String())
	}
}

func testTopMarketData(t *testing.T, client *ibkr.Client) {
	var resp *models.TopMarketDataResponse
	ctx, cancelCtx := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelCtx()

	err := client.RequestMarketDataType(ctx, models.MarketDataTypeRequestOptions{
		Type: models.MarketDataTypeDelayed,
	})
	if err != nil {
		testError(t, err)
		return
	}
	// ----
	resp, err = client.RequestTopMarketData(ctx, models.TopMarketDataRequestOptions{
		Contract: getContract("VOO", "SMART"),
	})
	if err != nil {
		testError(t, err)
		return
	}
	defer resp.Cancel()

	doneCh := time.After(20 * time.Second)
	for loop := true; loop; {
		select {
		case <-doneCh:
			loop = false

		case data, ok := <-resp.Channel:
			if !ok {
				if resp.Err() != nil {
					testError(t, resp.Err())
				}
				loop = false
				break
			}
			testLogf(
				t,
				"  %s, MarketDataType=%s, Delayed=%t, Frozen=%t",
				data.String(),
				data.MarketDataType().String(),
				data.Delayed(),
				data.Frozen(),
			)
		}
	}
}

func testDepthMarketData(t *testing.T, client *ibkr.Client) {
	ctx, cancelCtx := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelCtx()

	resp, err := client.RequestMarketDepthData(ctx, models.MarketDepthDataRequestOptions{
		Contract:  getContract("MSFT", "NASDAQ"),
		RowsCount: 20,
	})
	if err != nil {
		testError(t, err)
		return
	}
	defer resp.Cancel()

	doneCh := time.After(20 * time.Second)
	for loop := true; loop; {
		select {
		case <-doneCh:
			loop = false

		case data, ok := <-resp.Channel:
			if !ok {
				if resp.Err() != nil {
					if reqErr, ok := resp.Err().(*ibkr.RequestError); ok && reqErr.Code == 10089 {
						testLog(t, "  Skipping expected L2 subscription error:", reqErr)
					} else {
						testError(t, resp.Err())
					}
				}
				loop = false
				break
			}
			resp.Book.UpdateBook(data)
			testLog(t, resp.Book.TableString())
		}
	}
}

func testRealTimeBars(t *testing.T, client *ibkr.Client) {
	if !forexRealtimeBarsExpected(time.Now()) {
		testSkip(t, "forex market is closed or in the daily rollover window; realtime bars are not expected promptly")
	}

	ctx, cancelCtx := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancelCtx()

	resp, err := client.RequestRealTimeBars(ctx, models.RealTimeBarsRequestOptions{
		Contract:                getFxContract("EUR", "GBP"),
		BarSize:                 5,
		WhatToShow:              models.WhatToShowMidPoint,
		OnlyRegularTradingHours: true,
	})
	if err != nil {
		testError(t, err)
		return
	}
	defer resp.Cancel()

	select {
	case data, ok := <-resp.Channel:
		if !ok {
			if resp.Err() != nil {
				testError(t, resp.Err())
			} else {
				testError(t, "expected realtime bar, got closed channel")
			}
			return
		}
		testLog(t, "  "+data.String())

	case <-time.After(20 * time.Second):
		testError(t, "timeout waiting for realtime bar")
	}
}

func testTickByTickData(t *testing.T, client *ibkr.Client) {
	if !usStockRegularTradingHoursOpen(time.Now()) {
		testSkip(t, "US regular session is closed; stock tick-by-tick data is not expected to stream")
	}

	ctx, cancelCtx := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancelCtx()

	resp, err := client.RequestTickByTickData(ctx, models.TickByTickRequestOptions{
		Contract:      getContract("SPY", "SMART"),
		TickType:      models.TickByTickTypeLast,
		NumberOfTicks: 10,
	})
	if err != nil {
		testError(t, err)
		return
	}
	defer resp.Cancel()

	select {
	case data, ok := <-resp.Channel:
		if !ok {
			if resp.Err() != nil {
				testError(t, resp.Err())
			} else {
				testError(t, "expected tick-by-tick data, got closed channel")
			}
			return
		}
		testLog(t, "  "+data.String())

	case <-time.After(20 * time.Second):
		testError(t, "timeout waiting for tick-by-tick data")
	}
}

func testHistoricalDataKeepUpToDate(t *testing.T, client *ibkr.Client) {
	onlyRegularTradingHours := usStockRegularTradingHoursOpen(time.Now())
	if !onlyRegularTradingHours {
		testSkip(t, "US regular session is closed; keep-up-to-date stock historical data is not expected to stream promptly")
	}

	ctx, cancelCtx := context.WithTimeout(context.Background(), 50*time.Second)
	defer cancelCtx()

	resp, err := client.RequestHistoricalData(ctx, models.HistoricalDataRequestOptions{
		Contract:                getContract("AAPL", "SMART"),
		Duration:                1,
		DurationUnit:            models.DurationUnitDays,
		BarSize:                 models.BarSizeFiveSeconds,
		WhatToShow:              models.WhatToShowTrades,
		OnlyRegularTradingHours: onlyRegularTradingHours,
		KeepUpToDate:            true,
	})
	if err != nil {
		testError(t, err)
		return
	}

	select {
	case data, ok := <-resp.Channel:
		if !ok {
			if resp.Err() != nil {
				testError(t, resp.Err())
			} else {
				testError(t, "expected keep-up-to-date historical data, got closed channel")
			}
			return
		}
		testLog(t, "  "+data.String())

	case <-time.After(45 * time.Second):
		testError(t, "timeout waiting for keep-up-to-date historical data")
	}

	resp.Cancel()

	select {
	case _, ok := <-resp.Channel:
		if ok {
			testLog(t, "  received buffered historical item after cancel")
		}
	case <-time.After(5 * time.Second):
		testError(t, "timeout waiting for historical keep-up-to-date stream to close after cancel")
	}
}

func testSecurityDefinitionOptionParameters(t *testing.T, client *ibkr.Client) {
	ctx, cancelCtx := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancelCtx()

	resp, err := client.RequestSecurityDefinitionOptionParameters(ctx, models.SecurityDefinitionOptionParametersRequestOptions{
		UnderlyingSymbol:       "IBM",
		UnderlyingSecurityType: models.SecurityTypeStock,
		UnderlyingContractID:   8314,
	})
	if err != nil {
		testError(t, err)
		return
	}
	if len(resp.Parameters) == 0 {
		testError(t, "expected security definition option parameters")
		return
	}
	for _, param := range resp.Parameters {
		testLogf(t, "  exchange=%s class=%s multiplier=%s expirations=%d strikes=%d",
			param.Exchange, param.TradingClass, param.Multiplier, len(param.Expirations), len(param.Strikes))
	}
}

func getContract(symbol string, exchange string) *models.Contract {
	contract := models.NewContract()
	contract.Symbol = symbol
	contract.SecType = models.SecurityTypeStock
	contract.Currency = "USD"
	contract.Exchange = exchange
	return contract
}

func getFxContract(symbol string, currency string) *models.Contract {
	contract := models.NewContract()
	contract.Symbol = symbol
	contract.SecType = models.SecurityTypeForex
	contract.Currency = currency
	contract.Exchange = "IDEALPRO"
	return contract
}

func forexMarketOpen(now time.Time) bool {
	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		return false
	}

	now = now.In(loc)
	switch now.Weekday() {
	case time.Friday:
		hour, minute, _ := now.Clock()
		return hour < 17 || (hour == 17 && minute == 0)
	case time.Saturday:
		return false
	case time.Sunday:
		hour, minute, _ := now.Clock()
		return hour > 17 || (hour == 17 && minute >= 0)
	default:
		return true
	}
}

func forexRealtimeBarsExpected(now time.Time) bool {
	if !forexMarketOpen(now) {
		return false
	}

	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		return false
	}

	now = now.In(loc)
	hour, minute, _ := now.Clock()
	if hour == 17 && minute < 15 {
		return false
	}

	return true
}

func usStockRegularTradingHoursOpen(now time.Time) bool {
	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		return false
	}

	now = now.In(loc)
	if now.Weekday() == time.Saturday || now.Weekday() == time.Sunday {
		return false
	}

	hour, minute, _ := now.Clock()
	currentMinute := hour*60 + minute

	return currentMinute >= 9*60+30 && currentMinute < 16*60
}
