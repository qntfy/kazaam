package transform

import (
	"fmt"
	"strings"

	simplejson "github.com/bitly/go-simplejson"
)

// Coalesce checks multiple keys and returns the first matching key found.
func Coalesce(spec *Config, data *simplejson.Json) (*simplejson.Json, error) {
	if spec.Require == true {
		return nil, &Error{ErrMsg: fmt.Sprintf("Invalid spec. Coalesce does not support \"require\""), ErrType: SpecError}
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
					return nil, &Error{ErrMsg: fmt.Sprintf("Warn: Unable to coerce element to json string: %v", vItem), ErrType: ParseError}
				}
				keyList = append(keyList, vItemStr)
			}
		default:
			return nil, &Error{ErrMsg: fmt.Sprintf("Warn: Expected list in message for key: %s", k), ErrType: ParseError}
		}

		// iterate over keys to evaluate
		for _, v := range keyList {
			var dataForV *simplejson.Json
			var err error

			// grab the data
			dataForV, err = getJSONPath(data, v, false)
			if err != nil {
				return nil, err
			}
			if dataForV.Interface() != nil {
				data.SetPath(outPath, dataForV.Interface())
				break
			}
		}
	}
	return data, nil
}
