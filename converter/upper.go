package converter

import (
	"github.com/pkg/errors"
	"github.com/qntfy/kazaam/transform"
	"strconv"
	"strings"
)

type Upper struct {
	ConverterBase
}

func (c *Upper) Convert(jsonData []byte, value []byte, args []byte) (newValue []byte, err error) {
	newValue = value

	var jsonValue *transform.JSONValue
	jsonValue, err = c.NewJSONValue(value)
	if err != nil {
		return
	}

	if jsonValue.IsString() == false {
		err = errors.New("value must be string for upper converter")
		return
	}

	origValue := jsonValue.GetStringValue()

	newValue = []byte(strconv.Quote(strings.ToUpper(origValue)))

	return
}
