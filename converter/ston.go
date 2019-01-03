package converter

import (
	"errors"
	"github.com/qntfy/kazaam/transform"
)

type Ston struct {
	ConverterBase
}

func (c *Ston) Convert(jsonData []byte, value []byte, args []byte) (newValue []byte, err error) {
	var jsonValue *transform.JSONValue
	jsonValue, err = c.NewJSONValue(value)
	if err != nil {
		return
	}

	if jsonValue.IsNumber() {
		newValue = value
	} else if jsonValue.IsString() {
		jsonValue, err = c.NewJSONValue([]byte(jsonValue.GetStringValue()))
		if jsonValue.IsNumber() {
			newValue = jsonValue.GetData()
		} else {
			err = errors.New("string doesn't parse to number")
		}
	} else {
		err = errors.New("unexpected type")
	}

	return
}
