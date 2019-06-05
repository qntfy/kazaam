package converter

import (
	"github.com/mbordner/kazaam/transform"
)

type ConverterBase struct{}

func (c *ConverterBase) Init(config []byte) (err error) {
	return
}

func (c *ConverterBase) Convert(jsonData []byte, value []byte, args []byte) (newValue []byte, err error) {
	newValue = value
	return
}

func (c *ConverterBase) GetJsonPathValue(jsonData []byte, path string) (value *transform.JSONValue, err error) {
	return transform.GetJsonPathValue(jsonData, path)
}

func (c *ConverterBase) NewJSONValue(data []byte) (value *transform.JSONValue, err error) {
	return transform.NewJSONValue(data)
}
