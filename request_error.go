package ibkr

import (
	"errors"
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

// ErrSubscriptionRequired identifies request failures caused by missing IBKR API data subscriptions.
var ErrSubscriptionRequired = errors.New("subscription required")

// -----------------------------------------------------------------------------

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

// -----------------------------------------------------------------------------

// Is reports whether the request error represents a missing API data subscription.
func (r *RequestError) Is(target error) bool {
	if target != ErrSubscriptionRequired {
		return false
	}

	switch r.Code {
	case 10089, 10189, 2188:
		return true
	default:
		return false
	}
}

// -----------------------------------------------------------------------------

func newRequestError(ts time.Time, code int, message string, advancedOrderRejectJson string) *RequestError {
	return &RequestError{
		Timestamp:               ts,
		Code:                    code,
		Message:                 message,
		AdvancedOrderRejectJson: advancedOrderRejectJson,
	}
}
