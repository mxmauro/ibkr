package formatter

import (
	"math"
	"strconv"
)

// -----------------------------------------------------------------------------

func Int32String(val int32) string {
	return strconv.FormatInt(int64(val), 10)
}

func Int32MaxString(val *int32) string {
	if val == nil {
		return ""
	}
	return strconv.FormatInt(int64(*val), 10)
}

func Int64String(val int64) string {
	return strconv.FormatInt(val, 10)
}

func Int64MaxString(val *int64) string {
	if val == nil {
		return ""
	}
	return strconv.FormatInt(*val, 10)
}

func FloatString(val float64) string {
	return strconv.FormatFloat(val, 'g', 10, 64)
}

func FloatMaxString(val *float64) string {
	if val == nil || math.IsNaN(*val) || math.IsInf(*val, 0) {
		return ""
	}
	return strconv.FormatFloat(*val, 'g', 10, 64)
}

func BoolString(val bool) string {
	if val {
		return "1"
	}
	return "0"
}
