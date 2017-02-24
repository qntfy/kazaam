package transform

import (
	"fmt"
	"reflect"
	"strings"

	simplejson "github.com/bitly/go-simplejson"
)

// Concat combines any specified fields and literal strings into a single string value.
func Concat(spec *Config, data *simplejson.Json) (*simplejson.Json, error) {
	sourceList, sourceOk := (*spec.Spec)["sources"]
	if !sourceOk {
		return nil, &Error{ErrMsg: fmt.Sprintf("Unable to get sources"), ErrType: SpecError}
	}
	targetPath, targetOk := (*spec.Spec)["targetPath"]
	if !targetOk {
		return nil, &Error{ErrMsg: fmt.Sprintf("Unable to get targetPath"), ErrType: SpecError}
	}
	delimiter, delimOk := (*spec.Spec)["delim"]
	if !delimOk {
		// missing delimiter.  default to blank
		delimiter = ""
	}

	outString := ""
	applyDelim := false
	for _, vItem := range sourceList.([]interface{}) {
		if applyDelim {
			outString += delimiter.(string)
		}
		value, ok := vItem.(map[string]interface{})["value"]
		if !ok {
			path, ok := vItem.(map[string]interface{})["path"]
			if ok {
				valueNodePtr, err := getJSONPath(data, path.(string), spec.Require)
				switch {
				case err != nil && spec.Require == true:
					return nil, &Error{ErrMsg: fmt.Sprintf("Path does not exist"), ErrType: RequireError}
				case err != nil:
					value = ""
				default:
					zed := (*valueNodePtr).Interface()
					switch zed.(type) {
					case []interface{}:
						temp := ""
						for _, item := range zed.([]interface{}) {
							if item != nil {
								temp += reflect.ValueOf(item).String()
							}
						}
						value = temp
					default:
						value = reflect.ValueOf(zed).String()
					}
				}
			} else {
				return nil, &Error{ErrMsg: fmt.Sprintf("Error processing %v: must have either value or path specified", vItem), ErrType: SpecError}
			}
		}
		outString += value.(string)

		applyDelim = true
	}

	data.SetPath(strings.Split(targetPath.(string), "."), outString)

	return data, nil
}
