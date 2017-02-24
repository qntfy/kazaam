package transform

import (
	"fmt"
	"strings"

	simplejson "github.com/bitly/go-simplejson"
)

// Coalesce checks multiple keys and returns the first matching key found.
func Coalesce(spec *Config, data *simplejson.Json) error {
	if spec.Require == true {
		return SpecError("Invalid spec. Coalesce does not support \"require\"")
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
					return ParseError(fmt.Sprintf("Warn: Unable to coerce element to json string: %v", vItem))
				}
				keyList = append(keyList, vItemStr)
			}
		default:
			return ParseError(fmt.Sprintf("Warn: Expected list in message for key: %s", k))
		}

		// iterate over keys to evaluate
		for _, v := range keyList {
			var dataForV *simplejson.Json
			var err error

			// grab the data
			dataForV, err = getJSONPath(data, v, false)
			if err != nil {
				return err
			}
			if dataForV.Interface() != nil {
				data.SetPath(outPath, dataForV.Interface())
				break
			}
		}
	}
	return nil
}
