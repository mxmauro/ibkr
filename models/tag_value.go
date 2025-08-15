package models

import (
	"fmt"

	"github.com/mxmauro/ibkr/utils/encoders/message"
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

func (tvl *TagValueList) EncodeMessage(_ int) ([]byte, error) {
	msgEnc := message.NewRawEncoder()
	for _, tv := range *tvl {
		msgEnc.Raw([]byte(tv.Tag))
		msgEnc.Raw([]byte("="))
		msgEnc.Raw([]byte(tv.Value))
		msgEnc.Raw([]byte(";"))
	}
	msgEnc.AddDelim()
	return msgEnc.Bytes(), msgEnc.Err()
}
