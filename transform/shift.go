package transform

import (
	"fmt"
	"strings"

	simplejson "github.com/bitly/go-simplejson"
)

// Shift moves values from one provided json path to another.
func Shift(spec *Config, data *simplejson.Json) (*simplejson.Json, error) {
	outData := simplejson.New()
	for k, v := range *spec.Spec {
		array := true
		outPath := strings.Split(k, ".")

		var keyList []string

		// check if `v` is a string or list and build a list of keys to evaluate
		switch v.(type) {
		case string:
			keyList = append(keyList, v.(string))
			array = false
		case []interface{}:
			for _, vItem := range v.([]interface{}) {
				vItemStr, found := vItem.(string)
				if !found {
					return nil, &Error{ErrMsg: fmt.Sprintf("Warn: Unable to coerce element to json string: %v", vItem), ErrType: ParseError}
				}
				keyList = append(keyList, vItemStr)
			}
		default:
			return nil, &Error{ErrMsg: fmt.Sprintf("Warn: Unknown type in message for key: %s", k), ErrType: ParseError}
		}

		// iterate over keys to evaluate
		for _, v := range keyList {
			var dataForV *simplejson.Json
			var err error

			// grab the data
			if v == "$" {
				dataForV = data
			} else {
				dataForV, err = getJSONPath(data, v, spec.Require)
				if err != nil {
					return nil, err
				}
			}

			// if array flag set, encapsulate data
			if array {
				var intSlice = make([]interface{}, 1)
				intSlice[0] = dataForV.Interface()
				dataForV.SetPath(nil, intSlice)
			}

			outData.SetPath(outPath, dataForV.Interface())
		}
	}
	return outData, nil
}
