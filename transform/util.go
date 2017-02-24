// Package transform package contains canonical implementations of Kazaam transforms
package transform

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	simplejson "github.com/bitly/go-simplejson"
)

// Error provids an error mess (ErrMsg) and integer code (ErrType) for
// errors thrown by a transform
type Error struct {
	ErrMsg  string
	ErrType int
}

const (
	// ParseError is thrown when there is a JSON parsing error
	ParseError = iota
	// RequireError is thrown when the JSON path does not exist and is required
	RequireError
	// SpecError is thrown when the kazaam specification is not properly formatted
	SpecError
)

func (e *Error) Error() string {
	switch e.ErrType {
	case ParseError:
		return fmt.Sprintf("ParseError - %s", e.ErrMsg)
	case RequireError:
		return fmt.Sprintf("RequiredError - %s", e.ErrMsg)
	default:
		return fmt.Sprintf("SpecError - %s", e.ErrMsg)
	}
}

// Config contains the options that dictate the behavior of a transform. The internal
// `spec` object can be an arbitrary json configuration for the transform.
type Config struct {
	Spec    *map[string]interface{} `json:"spec"`
	Require bool                    `json:"require,omitempty"`
}

var jsonPathRe = regexp.MustCompile("([^\\[\\]]+)\\[([0-9\\*]+)\\]")

func getJSONPath(j *simplejson.Json, path string, pathRequired bool) (*simplejson.Json, error) {
	jin := j
	objectKeys := strings.Split(path, ".")
	// iterate over each subsequent object key
	for element, k := range objectKeys {
		// check the object key to see if it also contains an array reference
		results := jsonPathRe.FindAllStringSubmatch(k, -1)
		if results != nil {
			objKey := results[0][1]
			arrayKeyStr := results[0][2]
			// if there's a wildcard array reference
			if arrayKeyStr == "*" {
				// get the array
				if pathRequired {
					_, exists := jin.CheckGet(objKey)
					if exists != true {
						return nil, &Error{ErrMsg: fmt.Sprintf("Path does not exist"), ErrType: RequireError}
					}
				}
				jin = jin.Get(objKey)
				arrayLength := len(jin.MustArray())
				// construct the remainder of the jsonPath
				newPath := strings.Join(objectKeys[element+1:], ".")

				// iterate over each entry
				var results []interface{}
				for i := 0; i < arrayLength; i++ {
					if newPath == "" {
						results = append(results, jin.GetIndex(i).Interface())
					} else {
						intermediate, err := getJSONPath(jin.GetIndex(i), newPath, pathRequired)
						if err != nil {
							return nil, err
						}
						results = append(results, intermediate.Interface())
					}
				}

				output := simplejson.New()
				output.SetPath(nil, results)
				return output, nil
			}
			arrayKey, err := strconv.Atoi(arrayKeyStr)
			if err != nil {
				return nil, err
			}
			jin = jin.Get(objKey).GetIndex(arrayKey)
		} else {
			if pathRequired {
				_, exists := jin.CheckGet(k)
				if !exists {
					return nil, &Error{ErrMsg: fmt.Sprintf("Path does not exist"), ErrType: RequireError}
				}
			}
			jin = jin.Get(k)
		}
	}
	return jin, nil
}
