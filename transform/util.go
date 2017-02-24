// Package transform package contains canonical implementations of Kazaam transforms.
package transform

import (
	"regexp"
	"strconv"
	"strings"

	simplejson "github.com/bitly/go-simplejson"
)

// ParseError should be thrown when there is an issue with parsing any the specification or data
type ParseError string

func (p ParseError) Error() string {
	return string(p)
}

// RequireError should be thrown if a required key is missing in the data
type RequireError string

func (r RequireError) Error() string {
	return string(r)
}

// SpecError should be thrown if the spec for a transform is malformed
type SpecError string

func (s SpecError) Error() string {
	return string(s)
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
		if results != nil && len(results) > 0 {
			objKey := results[0][1]      // the key
			arrayKeyStr := results[0][2] // the array index
			// if there's a wildcard array reference
			if arrayKeyStr == "*" {
				// get the array
				var exists bool
				jin, exists = jin.CheckGet(objKey)
				if pathRequired && !exists {
					return nil, RequireError("Path does not exist")
					//return nil, &Error{ErrMsg: fmt.Sprintf("Path does not exist"), ErrType: RequireError}
				}
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
			// var exists bool
			// jin, exists = jin.CheckGet(k)
			// if pathRequired && !exists {
			// 	return nil, &Error{ErrMsg: fmt.Sprintf("Path does not exist"), ErrType: RequireError}
			// }
			if pathRequired {
				_, exists := jin.CheckGet(k)
				if !exists {
					return nil, RequireError("Path does not exist")
					// return nil, &Error{ErrMsg: fmt.Sprintf("Path does not exist"), ErrType: RequireError}
				}
			}
			jin = jin.Get(k)
		}
	}
	return jin, nil
}
