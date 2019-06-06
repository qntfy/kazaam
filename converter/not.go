package converter

import (
	"fmt"
	"github.com/mbordner/kazaam/transform"
)

type Not struct {
	ConverterBase
}

func (c *Not) Convert(jsonData []byte, value []byte, args []byte) (newValue []byte, err error) {

	var v *transform.JSONValue
	v, err = c.NewJSONValue(value)
	if err != nil {
		return
	}

	if v.IsBool() {
		return []byte(fmt.Sprintf("%t", !v.GetBoolValue())), nil
	}

	return []byte("false"), nil

}
