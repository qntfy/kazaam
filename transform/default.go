package transform

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/buger/jsonparser"
)

// Default sets specific value(s) in output json in raw []byte.
func Default(spec *Config, data []byte) ([]byte, error) {
	for k, v := range *spec.Spec {
		var err error
		dataForV, err := json.Marshal(v)
		if err != nil {
			return nil, ParseError(fmt.Sprintf("Warn: Unable to coerce element to json string: %v", v))
		}
		data, err = jsonparser.Set(data, dataForV, strings.Split(k, ".")...)
		if err != nil {
			return nil, err
		}
	}
	return data, nil
}
