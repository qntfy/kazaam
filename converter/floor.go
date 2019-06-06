package converter

import (
	"errors"
	"github.com/mbordner/kazaam/transform"
	"go/constant"
	"math"
)

type Floor struct {
	ConverterBase
}

func (c *Floor) Convert(jsonData []byte, value []byte, args []byte) (newValue []byte, err error) {

	newValue = value

	var jsonValue *transform.JSONValue
	jsonValue, err = c.NewJSONValue(value)
	if err != nil {
		return
	}
	if jsonValue.IsNumber() == false {
		err = errors.New("invalid value for floor converter")
		return
	}

	if jsonValue.GetType() == transform.JSONInt {
		return
	}

	val := jsonValue.GetFloatValue()

	val = math.Floor(val)

	newValue = []byte(constant.MakeInt64(int64(val)).String())

	return
}
