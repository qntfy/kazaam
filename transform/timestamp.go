package transform

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/buger/jsonparser"
)

// Timestamp parses and formats timestamp strings using the golang syntax
func Timestamp(spec *Config, data []byte) ([]byte, error) {
	ops, err := (*spec.Spec)["ops"].([]interface{})
	if !err {
		return nil, SpecError("Warn: Invalid spec. Unable to get \"ops\"")
	}
	for idx, v := range ops {
		assertedV := v.(map[string]interface{})
		path, pathErr := assertedV["path"].(string)
		if !pathErr {
			return nil, SpecError(fmt.Sprintf("Warn: Invalid spec. Unable to get \"path\" for item %d", idx))
		}
		inputFormat, inputErr := assertedV["inputFormat"].(string)
		if !inputErr {
			return nil, SpecError(fmt.Sprintf("Warn: Invalid spec. Unable to get \"inputFormat\" for item %d", idx))
		}
		outputFormat, outputErr := assertedV["outputFormat"].(string)
		if !outputErr {
			return nil, SpecError(fmt.Sprintf("Warn: Invalid spec. Unable to get \"outputFormat\" for item %d", idx))
		}
		// check if an array wildcard is present and if it is, treat it the
		// same as a key with an array
		if path[len(path)-2] == '*' {
			path = path[:len(path)-3]
		}
		// grab the data
		dataForV, err := getJSONRaw(data, path, spec.Require)
		if err != nil {
			return nil, err
		}
		// if the key is missing bail and keep iterating
		if bytes.Compare(dataForV, []byte("null")) == 0 {
			continue
		}
		// can only parse and format strings and arrays of strings, check the
		// value type and ahandle accordingly
		switch dataForV[0] {
		case '"':
			formattedItem, err := parseAndFormatValue(inputFormat, outputFormat, string(dataForV[1:len(dataForV)-1]))
			if err != nil {
				return nil, err
			}
			data, err = setPath(data, []byte(formattedItem), path)
			if err != nil {
				return nil, err
			}
		case '[':
			var unformattedItems, formattedItems []string
			_, err = jsonparser.ArrayEach(dataForV, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
				unformattedItems = append(unformattedItems, string(value))
			})
			if err != nil {
				return nil, err
			}
			for _, unformattedItem := range unformattedItems {
				formattedItem, err := parseAndFormatValue(inputFormat, outputFormat, unformattedItem)
				if err != nil {
					return nil, err
				}
				formattedItems = append(formattedItems, formattedItem)
			}
			data, err = setPath(data, bookend([]byte(strings.Join(formattedItems, ",")), '[', ']'), path)
			if err != nil {
				return nil, err
			}
		default:
			return nil, ParseError(fmt.Sprintf("Warn: Unknown type in message for key: %s", v))
		}
	}
	return data, nil
}

// parseAndFormatValue generates a properly formatted timestamp
func parseAndFormatValue(inputFormat, outputFormat, unformattedItem string) (string, error) {
	parsedItem, err := time.Parse(inputFormat, unformattedItem)
	if err != nil {
		return "", err
	}
	formattedItem := "\"" + parsedItem.Format(outputFormat) + "\""
	return formattedItem, nil
}

// set path updates the value with properly formatted timestamp(s) and properly
// handles array indexing
func setPath(data, out []byte, path string) ([]byte, error) {
	arrayRefs := jsonPathRe.FindAllStringSubmatch(path, -1)
	var splitRef string
	if arrayRefs != nil && len(arrayRefs) > 0 {
		splitRef = path[len(path)-3:]
		path = path[:len(path)-3]
	}
	splitPath := strings.Split(path, ".")
	if splitRef != "" {
		splitPath = append(splitPath, splitRef)
	}
	return jsonparser.Set(data, out, splitPath...)
}
