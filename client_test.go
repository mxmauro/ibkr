package ibkr_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/mxmauro/ibkr"
	"github.com/mxmauro/ibkr/models"
)

// -----------------------------------------------------------------------------

type Suite struct {
	t      *testing.T
	client *ibkr.Client
}

func TestSuite(t *testing.T) {
	suite := Suite{
		t: t,
	}

	suite.Setup()

	// -------------------------------------------------------------------------

	suite.RunParallel("Current time", testCurrentTime)
	suite.RunParallel("Managed accounts", testManagedAccounts)
	suite.RunParallel("Historical data", testHistoricalData)
	suite.RunParallel("Historical ticks", testHistoricalTicks)
	suite.RunParallel("Contract details", testContractDetails)
	suite.RunParallel("Matching symbols", testMatchingSymbols)
	suite.Run("Top market data", testTopMarketData)
	suite.Run("Depth market data", testDepthMarketData)
}

// -----------------------------------------------------------------------------

func (suite *Suite) Setup() {
	var err error

	suite.t.Log("Connecting to server...")

	ctx, cancelCtx := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelCtx()

	suite.client, err = ibkr.NewClient(ctx, ibkr.Options{
		Address: "127.0.0.1:4002",
		EventsHandler: ibkr.NewIncomingMessageLogger(func(msg string) {
			suite.t.Log(msg)
		}),
	})
	if err != nil {
		suite.t.Fatal(err)
	}

	suite.t.Cleanup(func() {
		suite.client.Destroy()
	})
}

func (suite *Suite) Run(name string, f func(t *testing.T, client *ibkr.Client)) {
	suite.t.Run(name, func(t *testing.T) {
		f(t, suite.client)
	})
}

func (suite *Suite) RunParallel(name string, f func(t *testing.T, client *ibkr.Client)) {
	suite.t.Run(name, func(t *testing.T) {
		t.Parallel()
		f(t, suite.client)
	})
}

// -----------------------------------------------------------------------------

func testCurrentTime(t *testing.T, client *ibkr.Client) {
	ctx, cancelCtx := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelCtx()

	ts, err := client.RequestCurrentTime(ctx)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("  Time is: " + ts.String())
}

func testManagedAccounts(t *testing.T, client *ibkr.Client) {
	ctx, cancelCtx := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelCtx()

	accounts, err := client.RequestManagedAccounts(ctx)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("  Accounts: " + strings.Join(accounts, ", "))
}

func testHistoricalData(t *testing.T, client *ibkr.Client) {
	ctx, cancelCtx := context.WithTimeout(context.Background(), 10*time.Second)
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
		t.Error(err)
		return
	}
	for _, bar := range histData.Bars {
		t.Log("  " + bar.String())
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
		t.Error(err)
		return
	}
	for _, tick := range histTicks.TicksLast {
		t.Log("  " + tick.String())
	}
}

func testContractDetails(t *testing.T, client *ibkr.Client) {
	ctx, cancelCtx := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelCtx()

	resp, err := client.RequestContractDetails(ctx, models.ContractDetailsRequestOptions{
		Contract: getContract("JPM", "SMART"),
	})
	if err != nil {
		t.Error(err)
		return
	}
	for _, cd := range resp.ContractDetails {
		t.Log("  " + cd.String())
	}
}

func testMatchingSymbols(t *testing.T, client *ibkr.Client) {
	ctx, cancelCtx := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelCtx()

	resp, err := client.RequestMatchingSymbols(ctx, models.MatchingSymbolsRequestOptions{
		Pattern: "AAPL",
	})
	if err != nil {
		t.Error(err)
		return
	}
	for _, cd := range resp.ContractDescriptions {
		t.Log("  " + cd.String())
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
		t.Error(err)
		return
	}
	// ----
	resp, err = client.RequestTopMarketData(ctx, models.TopMarketDataRequestOptions{
		Contract: getContract("VOO", "SMART"),
	})
	if err != nil {
		t.Error(err)
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
					t.Error(resp.Err())
				}
				loop = false
				break
			}
			t.Log("  " + data.String())
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
		t.Error(err)
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
					t.Error(resp.Err())
				}
				loop = false
				break
			}
			data.UpdateBook(&resp.Book)
			t.Log(resp.Book.TableString())
		}
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
