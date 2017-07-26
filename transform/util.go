// Package transform package contains canonical implementations of Kazaam transforms.
package transform

import (
	"bytes"
	"regexp"
	"strconv"
	"strings"

	"github.com/buger/jsonparser"
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

// Given a json byte slice `data` and a kazaam `path` string, return the object at the path in data if it exists.
func getJSONRaw(data []byte, path string, pathRequired bool) ([]byte, error) {
	objectKeys := strings.Split(path, ".")
	numOfInserts := 0
	for element, k := range objectKeys {
		// check the object key to see if it also contains an array reference
		arrayRefs := jsonPathRe.FindAllStringSubmatch(k, -1)
		if arrayRefs != nil && len(arrayRefs) > 0 {
			objKey := arrayRefs[0][1]      // the key
			arrayKeyStr := arrayRefs[0][2] // the array index
			// if there's a wildcard array reference
			if arrayKeyStr == "*" {
				// ArrayEach setup
				objectKeys[element+numOfInserts] = objKey
				beforePath := objectKeys[:element+numOfInserts+1]
				newPath := strings.Join(objectKeys[element+numOfInserts+1:], ".")
				var results [][]byte

				// use jsonparser.ArrayEach to copy the array into results
				_, err := jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
					results = append(results, value)
				}, beforePath...)
				if err == jsonparser.KeyPathNotFoundError {
					if pathRequired {
						return nil, RequireError("Path does not exist")
					}
				} else if err != nil {
					return nil, err
				}

				// GetJSONRaw() the rest of path for each element in results
				if newPath != "" {
					for i, value := range results {
						intermediate, err := getJSONRaw(value, newPath, pathRequired)
						if err != nil {
							return nil, err
						}
						results[i] = intermediate
					}
				}

				// copy into raw []byte format and return
				var buffer bytes.Buffer
				buffer.WriteByte('[')
				for i := 0; i < len(results)-1; i++ {
					buffer.Write(results[i])
					buffer.WriteByte(',')
				}
				if len(results) > 0 {
					buffer.Write(results[len(results)-1])
				}
				buffer.WriteByte(']')
				return buffer.Bytes(), nil
			}
			_, err := strconv.Atoi(arrayKeyStr)
			if err != nil {
				return nil, err
			}
			// separate the array key as a new element in objectKeys
			objectKeys = makePathWithIndex(arrayKeyStr, objKey, objectKeys, element, numOfInserts)
			numOfInserts++
		} else {
			// no array reference, good to go
			continue
		}
	}
	result, dataType, _, err := jsonparser.Get(data, objectKeys...)

	// jsonparser strips quotes from Strings
	if dataType == jsonparser.String {
		// bookend() is destructive to underlying slice, need to copy.
		// extra capacity saves an allocation and copy during bookend.
		tmp := make([]byte, len(result), len(result)+2)
		copy(tmp, result)
		result = bookend(tmp, '"', '"')
	}
	if len(result) == 0 {
		result = []byte("null")
	}
	if err == jsonparser.KeyPathNotFoundError {
		if pathRequired {
			return nil, RequireError("Path does not exist")
		}
	} else if err != nil {
		return nil, err
	}
	return result, nil
}

// setPath updates the value with properly formatted timestamp(s) and properly
// handles array indexing
func setJSONRaw(data, out []byte, path string) ([]byte, error) {
	var err error
	splitPath := strings.Split(path, ".")
	numOfInserts := 0

	for element, k := range splitPath {
		arrayRefs := jsonPathRe.FindAllStringSubmatch(k, -1)
		if arrayRefs != nil && len(arrayRefs) > 0 {
			objKey := arrayRefs[0][1]      // the key
			arrayKeyStr := arrayRefs[0][2] // the array index
			if arrayKeyStr == "*" {
				// ArrayEach setup
				splitPath[element+numOfInserts] = objKey
				beforePath := splitPath[:element+numOfInserts+1]
				afterPath := strings.Join(splitPath[element+numOfInserts+1:], ".")
				// use jsonparser.ArrayEach to count the number of items in the
				// array
				var arraySize int
				_, err = jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
					arraySize++
				}, beforePath...)
				if err != nil {
					return nil, err
				}

				// setJSONRaw() the rest of path for each element in results
				for i := 0; i < arraySize; i++ {
					var newPath string
					// iterate through each item in the array by replacing the
					// wildcard with an int and joining the path back together
					newArrayKey := strings.Join([]string{"[", strconv.Itoa(i), "]"}, "")
					beforePathStr := strings.Join(beforePath, ".")
					beforePathArrayKeyStr := strings.Join([]string{beforePathStr, newArrayKey}, "")
					// if there's nothing that comes after the array index,
					// don't join so that we avoid trailing cruft
					if len(afterPath) > 0 {
						newPath = strings.Join([]string{beforePathArrayKeyStr, afterPath}, ".")
					} else {
						newPath = beforePathArrayKeyStr
					}
					// now call the function, but this time with an array index
					// instead of a wildcard
					data, err = setJSONRaw(data, out, newPath)
					if err != nil {
						return nil, err
					}
				}
				return data, nil
			}
			// if not a wildcard then piece that path back together with the
			// array index as an entry in the splitPath slice
			splitPath = makePathWithIndex(arrayKeyStr, objKey, splitPath, element, numOfInserts)
			numOfInserts++
		} else {
			continue
		}
	}
	data, err = jsonparser.Set(data, out, splitPath...)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// makePathWithIndex generats a path slice to pass to jsonparser
func makePathWithIndex(arrayKeyStr, objectKey string, pathSlice []string, element, numOfInserts int) []string {
	arrayKey := string(bookend([]byte(arrayKeyStr), '[', ']'))
	pathSlice[element+numOfInserts] = objectKey
	pathSlice = append(pathSlice, "")
	copy(pathSlice[element+numOfInserts+2:], pathSlice[element+numOfInserts+1:])
	pathSlice[element+numOfInserts+1] = arrayKey
	return pathSlice
}

// add characters at beginning and end of []byte
func bookend(value []byte, bef, aft byte) []byte {
	value = append(value, ' ', aft)
	copy(value[1:], value[:len(value)-2])
	value[0] = bef
	return value
}
