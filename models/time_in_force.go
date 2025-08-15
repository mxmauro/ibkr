package models

// -----------------------------------------------------------------------------

type TimeInForce string

const (
	TimeInForceNone TimeInForce = ""
	TimeInForceDay  TimeInForce = "DAY" // Day order
	TimeInForceGTC  TimeInForce = "GTC" // Good-Til-Canceled
	TimeInForceOPG  TimeInForce = "OPG" // On-Open (MOO/LOO)
	TimeInForceIOC  TimeInForce = "IOC" // Immediate-Or-Cancel
	TimeInForceGTD  TimeInForce = "GTD" // Good-Til-Date
	TimeInForceGTT  TimeInForce = "GTT" // Good-Til-Time
	TimeInForceAUC  TimeInForce = "AUC" // Auction
	TimeInForceFOK  TimeInForce = "FOK" // Fill-Or-Kill
	TimeInForceGTX  TimeInForce = "GTX" // Good-Til-Canceled (variant)
	TimeInForceDTC  TimeInForce = "DTC" // Day-Till-Canceled
)

// -----------------------------------------------------------------------------

func NewTimeInForceFromString(t string) TimeInForce {
	switch t {
	case string(TimeInForceDay):
		return TimeInForceDay
	case string(TimeInForceGTC):
		return TimeInForceGTC
	case string(TimeInForceOPG):
		return TimeInForceOPG
	case string(TimeInForceIOC):
		return TimeInForceIOC
	case string(TimeInForceGTD):
		return TimeInForceGTD
	case string(TimeInForceGTT):
		return TimeInForceGTT
	case string(TimeInForceAUC):
		return TimeInForceAUC
	case string(TimeInForceFOK):
		return TimeInForceFOK
	case string(TimeInForceGTX):
		return TimeInForceGTX
	case string(TimeInForceDTC):
		return TimeInForceDTC
	}
	return TimeInForceNone
}
