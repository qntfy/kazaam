// Package transform package contains canonical implementations of Kazaam transforms.
package transform

import (
	"strconv"
	"strings"

	simplejson "github.com/bitly/go-simplejson"
)

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
				if !exists {
					if pathRequired {
						return jin, RequireError("Path does not exist")
					}
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
