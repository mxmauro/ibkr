package protofmt

import (
	"math"
	"strconv"
)

// -----------------------------------------------------------------------------

func Bool(val bool) *bool {
	return &val
}

func Int32(val int32) *int32 {
	return &val
}

func Int32Max(val *int32) *int32 {
	if val == nil {
		return nil
	}
	copyOfVal := *val
	return &copyOfVal
}

func Int64(val int64) *int64 {
	return &val
}

func Int64Max(val *int64) *int64 {
	if val == nil {
		return nil
	}
	copyOfVal := *val
	return &copyOfVal
}

func Float(val float64) *float64 {
	if math.IsNaN(val) || math.IsInf(val, 0) {
		return nil
	}
	return &val
}

func FloatMax(val *float64) *float64 {
	if val == nil || math.IsNaN(*val) || math.IsInf(*val, 0) {
		return nil
	}
	copyOfVal := *val
	return &copyOfVal
}

func StringFloat(val string) *float64 {
	f, err := strconv.ParseFloat(val, 64)
	if err != nil || math.IsNaN(f) || math.IsInf(f, 0) {
		return nil
	}
	return &f
}

func String(val string) *string {
	if len(val) == 0 {
		return nil
	}
	return &val
}
