package converter

import (
	"bytes"
	"github.com/mbordner/kazaam/transform"
)

type Eqs struct {
	ConverterBase
}

func (c *Eqs) Convert(jsonData []byte, value []byte, args []byte) (newValue []byte, err error) {

	var argsValue *transform.JSONValue
	argsValue, err = c.NewJSONValue(args)
	if err != nil {
		return
	}

	argsBytes := []byte(argsValue.GetStringValue())

	if bytes.Equal(value, argsBytes) == true {
		return []byte("true"), nil
	}

	return []byte("false"), nil

}
