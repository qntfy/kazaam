package transform

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func inArray(val []byte, arr [][]byte) bool {
	for _, arrVal := range arr {
		if bytes.Compare(val, arrVal) == 0 {
			return true
		}
	}
	return false
}

// Coalesce checks multiple keys and returns the first matching key found in raw []byte.
func Coalesce(spec *Config, data []byte) ([]byte, error) {
	if spec.Require == true {
		return nil, SpecError("Invalid spec. Coalesce does not support \"require\"")
	}

	ignoreSlice := [][]byte{[]byte("null")}
	ignoreList, ignoreOk := (*spec.Spec)["ignore"]
	if ignoreOk {
		for _, iItem := range ignoreList.([]interface{}) {
			iByte, err := json.Marshal(iItem)
			if err != nil {
				return nil, SpecError(fmt.Sprintf("Warn: Could not marshal ignore item: %v", iItem))
			}
			ignoreSlice = append(ignoreSlice, iByte)
		}
	}

	for k, v := range *spec.Spec {
		if k == "ignore" {
			continue
		}

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
			if !inArray(dataForV, ignoreSlice) {
				data, err = setJSONRaw(data, dataForV, k)
				if err != nil {
					return nil, err
				}
				break
			}
		}
	}
	return data, nil

}
