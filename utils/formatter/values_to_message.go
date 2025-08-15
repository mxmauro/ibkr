package formatter

import (
	"math"
	"strconv"
)

// -----------------------------------------------------------------------------

func IntMsg(val int) string {
	return strconv.FormatInt(int64(val), 10)
}

func Int32Msg(val int32) string {
	return strconv.FormatInt(int64(val), 10)
}

func Int32MaxMsg(val *int32) string {
	if val == nil {
		return ""
	}
	return strconv.FormatInt(int64(*val), 10)
}

func Int64Msg(val int64) string {
	return strconv.FormatInt(val, 10)
}

func Int64MaxMsg(val *int64) string {
	if val == nil {
		return ""
	}
	return strconv.FormatInt(*val, 10)
}

func FloatMsg(val float64) string {
	return strconv.FormatFloat(val, 'f', -1, 64)
}

func FloatMaxMsg(val *float64) string {
	if val == nil || math.IsNaN(*val) || math.IsInf(*val, 0) {
		return ""
	}
	return strconv.FormatFloat(*val, 'f', -1, 64)
}

func BoolMsg(val bool) string {
	if val {
		return "1"
	}
	return "0"
}
