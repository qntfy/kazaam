package converter

import (
	"errors"
	"github.com/mbordner/kazaam/transform"
	"strconv"
	"strings"
)

type Lower struct {
	ConverterBase
}

func (c *Lower) Convert(jsonData []byte, value []byte, args []byte) (newValue []byte, err error) {
	newValue = value

	var jsonValue *transform.JSONValue
	jsonValue, err = c.NewJSONValue(value)
	if err != nil {
		return
	}

	if jsonValue.IsString() == false {
		err = errors.New("value must be string for lower converter")
		return
	}

	origValue := jsonValue.GetStringValue()

	newValue = []byte(strconv.Quote(strings.ToLower(origValue)))

	return
}
