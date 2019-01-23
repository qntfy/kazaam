package converter

import (
	"fmt"
	"github.com/mbordner/kazaam/transform"
	"github.com/pkg/errors"
)

type Len struct {
	ConverterBase
}

func (c *Len) Convert(jsonData []byte, value []byte, args []byte) (newValue []byte, err error) {
	newValue = value

	var jsonValue *transform.JSONValue
	jsonValue, err = c.NewJSONValue(value)
	if err != nil {
		return
	}

	if jsonValue.IsString() == false {
		err = errors.New("value must be string for len converter")
		return
	}

	origValue := jsonValue.GetStringValue()

	newValue = []byte(fmt.Sprintf("%d", len(origValue)))

	return
}
