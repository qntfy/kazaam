package converter

import (
	"errors"
	"github.com/mbordner/kazaam/transform"
	"strconv"
)

type Float struct {
	ConverterBase
}

func (c *Float) Convert(jsonData []byte, value []byte, args []byte) (newValue []byte, err error) {

	var jsonValue,argsValue *transform.JSONValue
	jsonValue, err = c.NewJSONValue(value)
	if err != nil {
		return
	}

	argsValue, err = c.NewJSONValue(args)
	if err != nil {
		return
	}


	if jsonValue.IsNumber() == false || argsValue.IsString() == false {
		err = errors.New("invalid value or arguments for float converter")
		return
	}

	numStrVal := argsValue.GetStringValue()
	if numStrVal[0] == '.' {
		numStrVal = "0" + numStrVal
	}

	// convert the string to number
	argsValue, err = c.NewJSONValue([]byte(numStrVal))
	if err != nil {
		return
	}

	if argsValue.IsNumber() == false {
		err = errors.New("arguments should be a number for float converter")
		return
	}

	precision := argsValue.GetIntValue()

	var val float64

	if jsonValue.GetType() == transform.JSONInt {
		val = float64(jsonValue.GetIntValue())
	} else {
		val = jsonValue.GetFloatValue()
	}

	newValue = []byte(strconv.FormatFloat(val, 'f', int(precision), 64))

	return
}
