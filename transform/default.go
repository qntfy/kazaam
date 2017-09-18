package transform

import (
	"encoding/json"
	"fmt"
)

// Default sets specific value(s) in output json in raw []byte.
func Default(spec *Config, data []byte) ([]byte, error) {
	for k, v := range *spec.Spec {
		var err error
		dataForV, err := json.Marshal(v)
		if err != nil {
			return nil, ParseError(fmt.Sprintf("Warn: Unable to coerce element to json string: %v", v))
		}
		data, err = setJSONRaw(data, dataForV, k)
		if err != nil {
			return nil, err
		}
	}
	return data, nil
}
