package models

import (
	"bytes"
	"fmt"

	"github.com/mxmauro/ibkr/utils"
)

// -----------------------------------------------------------------------------

// TagValue maps a tag to a value. Both of them are strings.
// They are used in a slice to convey extra info with the requests.
type TagValue struct {
	Tag   string
	Value string
}

type TagValueList []TagValue

// -----------------------------------------------------------------------------

func NewTagValue() TagValue {
	return TagValue{}
}

func (tv TagValue) String() string {
	return fmt.Sprintf("%s=%s", tv.Tag, tv.Value)
}

func (tvl *TagValueList) EncodeMessage() []byte {
	buf := bytes.Buffer{}
	for _, tv := range *tvl {
		_, _ = buf.WriteString(tv.Tag)
		_, _ = buf.WriteRune('=')
		_, _ = buf.WriteString(tv.Value)
		_, _ = buf.WriteRune(';')
	}
	_ = buf.WriteByte(utils.Delimiter())
	return buf.Bytes()
}
