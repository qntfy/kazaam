package converter

import (
	"errors"
	"github.com/mbordner/kazaam/transform"
	"strconv"
)

type Ntos struct {
	ConverterBase
}

func (c *Ntos) Convert(jsonData []byte, value []byte, args []byte) (newValue []byte, err error) {
	var jsonValue *transform.JSONValue
	jsonValue, err = c.NewJSONValue(value)
	if err != nil {
		return
	}

	if jsonValue.IsString() {
		newValue = value
	} else if jsonValue.IsNumber() {
		num := jsonValue.GetNumber()
		newValue = []byte(strconv.Quote(num.String()))
	} else {
		err = errors.New("unexpected type")
	}

	return
}
