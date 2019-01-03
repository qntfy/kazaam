package converter

import (
	"errors"
	"github.com/qntfy/kazaam/transform"
	"go/constant"
	"math"
)

type Ceil struct {
	ConverterBase
}

func (c *Ceil) Convert(jsonData []byte, value []byte, args []byte) (newValue []byte, err error) {

	newValue = value

	var jsonValue *transform.JSONValue
	jsonValue, err = c.NewJSONValue(value)
	if err != nil {
		return
	}
	if jsonValue.IsNumber() == false {
		err = errors.New("invalid value for ceil converter")
		return
	}

	if jsonValue.GetType() == transform.JSONInt {
		return
	}

	val := jsonValue.GetFloatValue()

	val = math.Ceil(val)

	newValue = []byte(constant.MakeInt64(int64(val)).String())

	return
}
