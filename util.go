package kazaam

import (
	"encoding/json"

	"github.com/buger/jsonparser"
)

// by default, kazaam does not fully validate input data. Use IsJson()
// if you need to confirm input is valid before transforming.
// Note: This operation is very slow and memory/alloc intensive
// relative to most transforms.
func IsJson(s []byte) bool {
	var js map[string]interface{}
	return json.Unmarshal(s, &js) == nil

}

// experimental fast validation with jsonparser
func IsJsonFast(s []byte) bool {
	for _, c := range s {
		switch c {
		case ' ', '\n', '\r', '\t':
			continue
		case '{':
			return isJsonInternal(s, jsonparser.Object)
		case '[':
			return isJsonInternal(s, jsonparser.Array)
		default:
			return false
		}
	}
	return false
}

func isJsonInternal(s []byte, t jsonparser.ValueType) bool {
	valid := true
	if t == jsonparser.Array {
		_, err := jsonparser.ArrayEach(s, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			if valid {
				valid = isJsonInternal(value, dataType)
			}
		})
		if err != nil || !valid {
			return false
		}
	} else if t == jsonparser.Object {
		err := jsonparser.ObjectEach(s, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
			if valid {
				valid = isJsonInternal(value, dataType)
			}
			return nil
		})
		if err != nil || !valid {
			return false
		}
	} else if t == jsonparser.Unknown {
		return false
	}
	return valid
}
