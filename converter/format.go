package converter

import (
	"errors"
	"fmt"
	"github.com/mbordner/kazaam/transform"
	"strconv"
)

type Format struct {
	ConverterBase
}

func (c *Format) Convert(jsonData []byte, value []byte, args []byte) (newValue []byte, err error) {
	newValue = value

	var jsonValue, argsValue *transform.JSONValue
	jsonValue, err = c.NewJSONValue(value)
	if err != nil {
		return
	}

	argsValue, err = transform.NewJSONValue(args)
	if err != nil {
		return
	}

	if argsValue.IsString() == false {
		err = errors.New("invalid value or arguments for substr converter")
		return
	}

	newValue = []byte(strconv.Quote(fmt.Sprintf(argsValue.GetStringValue(), jsonValue.GetValue())))

	return
}
