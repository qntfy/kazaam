package converter

import (
	"errors"
	"github.com/qntfy/kazaam/transform"
	"go/constant"
	"go/token"
)

type Add struct {
	ConverterBase
}

func (c *Add) Convert(jsonData []byte, value []byte, args []byte) (newValue []byte, err error) {

	var jsonValue, argsValue *transform.JSONValue
	jsonValue, err = c.NewJSONValue(value)
	if err != nil {
		return
	}

	argsValue, err = c.NewJSONValue(args)
	if err != nil {
		return
	}

	if jsonValue.IsNumber() == false || argsValue.IsString() == false {
		err = errors.New("invalid value or arguments for add converter")
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
		err = errors.New("arguments should be a number for add converter")
		return
	}

	var left, right constant.Value

	if jsonValue.GetType() == transform.JSONInt {
		left = constant.MakeInt64(jsonValue.GetIntValue())
	} else {
		left = constant.MakeFloat64(jsonValue.GetFloatValue())
	}

	if argsValue.GetType() == transform.JSONInt {
		right = constant.MakeInt64(argsValue.GetIntValue())
	} else {
		right = constant.MakeFloat64(argsValue.GetFloatValue())
	}

	result := constant.BinaryOp(left, token.ADD, right)

	newValue = []byte(result.String())

	return
}
