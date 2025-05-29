package utils

import (
	"strconv"

	"github.com/mxmauro/ibkr/common"
)

// -----------------------------------------------------------------------------

func IntString(val int64) string {
	return strconv.FormatInt(val, 10)
}

func IntMaxString(val int64) string {
	if val == common.UNSET_INT {
		return ""
	}
	return strconv.FormatInt(val, 10)
}

func FloatString(val float64) string {
	return strconv.FormatFloat(val, 'g', 10, 64)
}

func FloatMaxString(val float64) string {
	if val == common.UNSET_FLOAT {
		return ""
	}
	return strconv.FormatFloat(val, 'g', 10, 64)
}

func BoolString(val bool) string {
	if val {
		return "True"
	}
	return "False"
}
