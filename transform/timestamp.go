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
	for k, v := range *spec.Spec {
		assertedV, vErr := v.(map[string]interface{})
		if !vErr {
			return nil, SpecError(fmt.Sprintf("Warn: Invalid spec. Unable to get value for key: %s", k))
		}
		inputFormat, inputErr := assertedV["inputFormat"].(string)
		if !inputErr {
			return nil, SpecError(fmt.Sprintf("Warn: Invalid spec. Unable to get \"inputFormat\" for key: %s", k))
		}
		outputFormat, outputErr := assertedV["outputFormat"].(string)
		if !outputErr {
			return nil, SpecError(fmt.Sprintf("Warn: Invalid spec. Unable to get \"outputFormat\" for key: %s", k))
		}
		// check if an array wildcard is present and if it is, treat it the
		// same as a key with an array
		if k[len(k)-2] == '*' {
			k = k[:len(k)-3]
		}
		// grab the data
		dataForV, err := getJSONRaw(data, k, spec.Require)
		if err != nil {
			return nil, err
		}
		// if the key is missing bail and keep iterating
		if bytes.Compare(dataForV, []byte("null")) == 0 {
			continue
		}
		// can only parse and format strings and arrays of strings, check the
		// value type and handle accordingly
		var outData []byte
		switch dataForV[0] {
		case '"':
			formattedItem, err := parseAndFormatValue(inputFormat, outputFormat, string(dataForV[1:len(dataForV)-1]))
			if err != nil {
				return nil, err
			}
			outData = []byte(formattedItem)
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
			outData = bookend([]byte(strings.Join(formattedItems, ",")), '[', ']')
		default:
			return nil, ParseError(fmt.Sprintf("Warn: Unknown type in message for key: %s", v))
		}
		data, err = setPath(data, outData, k)
		if err != nil {
			return nil, err
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

// set k updates the value with properly formatted timestamp(s) and properly
// handles array indexing
func setPath(data, out []byte, k string) ([]byte, error) {
	arrayRefs := jsonPathRe.FindAllStringSubmatch(k, -1)
	var splitRef string
	if arrayRefs != nil && len(arrayRefs) > 0 {
		splitRef = k[len(k)-3:]
		k = k[:len(k)-3]
	}
	splitPath := strings.Split(k, ".")
	if splitRef != "" {
		splitPath = append(splitPath, splitRef)
	}
	return jsonparser.Set(data, out, splitPath...)
}
