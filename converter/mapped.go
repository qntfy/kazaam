package converter

import (
	"encoding/json"
	"errors"
	"github.com/qntfy/kazaam/transform"
	"strconv"
)

var mappedSpecs map[string]map[string]string

type Mapped struct {
	ConverterBase
}

func (c *Mapped) Init(config []byte) (err error) {
	if err := json.Unmarshal(config, &mappedSpecs); err != nil {
		mappedSpecs = make(map[string]map[string]string)
	}
	return
}

func (c *Mapped) Convert(jsonData []byte, value []byte, args []byte) (newValue []byte, err error) {

	newValue = value

	var jsonValue, argsValue *transform.JSONValue
	jsonValue, err = c.NewJSONValue(value)
	if err != nil {
		return
	}

	argsValue, err = c.NewJSONValue(args)
	if err != nil {
		return
	}

	if jsonValue.IsString() == false || argsValue.IsString() == false {
		err = errors.New("invalid value or arguments for mapped converter")
		return
	}

	mappedCollectionName := argsValue.GetStringValue()
	valueToMap := jsonValue.GetStringValue()

	if group, ok := mappedSpecs[mappedCollectionName]; ok {
		if newValueStr, ok := group[valueToMap]; ok {
			newValue = []byte(strconv.Quote(newValueStr))
		}
	}

	return
}
