package transform

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/buger/jsonparser"
)

// Coalesce checks multiple keys and returns the first matching key found in raw []byte.
func Coalesce(spec *Config, data []byte) ([]byte, error) {
	if spec.Require == true {
		return nil, SpecError("Invalid spec. Coalesce does not support \"require\"")
	}
	for k, v := range *spec.Spec {
		outPath := strings.Split(k, ".")

		var keyList []string

		// check if `v` is a list and build a list of keys to evaluate
		switch v.(type) {
		case []interface{}:
			for _, vItem := range v.([]interface{}) {
				vItemStr, found := vItem.(string)
				if !found {
					return nil, ParseError(fmt.Sprintf("Warn: Unable to coerce element to json string: %v", vItem))
				}
				keyList = append(keyList, vItemStr)
			}
		default:
			return nil, ParseError(fmt.Sprintf("Warn: Expected list in message for key: %s", k))
		}

		// iterate over keys to evaluate
		for _, v := range keyList {
			var dataForV []byte
			var err error

			// grab the data
			dataForV, err = getJSONRaw(data, v, false)
			if err != nil {
				return nil, err
			}
			if bytes.Compare(dataForV, []byte("null")) != 0 {
				data, err = jsonparser.Set(data, dataForV, outPath...)
				if err != nil {
					return nil, err
				}
				break
			}
		}
	}
	return data, nil

}
