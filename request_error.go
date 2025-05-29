package ibkr

import (
	"strconv"
	"strings"
	"time"
)

// -----------------------------------------------------------------------------

type RequestError struct {
	Timestamp               time.Time
	Code                    int
	Message                 string
	AdvancedOrderRejectJson string
}

// -----------------------------------------------------------------------------

func newRequestError(ts time.Time, code int64, message string, advancedOrderRejectJson string) *RequestError {
	return &RequestError{
		Timestamp:               ts,
		Code:                    int(code),
		Message:                 message,
		AdvancedOrderRejectJson: advancedOrderRejectJson,
	}
}

func (r *RequestError) Error() string {
	sb := strings.Builder{}
	_, _ = sb.WriteString("error ")
	_, _ = sb.WriteString(strconv.Itoa(r.Code))
	if len(r.Message) > 0 {
		_, _ = sb.WriteString(" (")
		_, _ = sb.WriteString(r.Message)
		_, _ = sb.WriteString(")")
	}
	return sb.String()
}
