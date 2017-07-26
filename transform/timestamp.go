package transform

import (
	"bytes"
	"fmt"
	"strconv"
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
		//if k[len(k)-2] == '*' {
		//	k = k[:len(k)-3]
		//}
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
		switch dataForV[0] {
		case '"':
			formattedItem, err := parseAndFormatValue(inputFormat, outputFormat, string(dataForV[1:len(dataForV)-1]))
			if err != nil {
				return nil, err
			}
			data, err = setJSONRaw(data, []byte(formattedItem), k)
			if err != nil {
				return nil, err
			}
		case '[':
			var unformattedItems []string
			_, err = jsonparser.ArrayEach(dataForV, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
				unformattedItems = append(unformattedItems, string(value))
			})
			if err != nil {
				return nil, err
			}
			for idx, unformattedItem := range unformattedItems {
				formattedItem, err := parseAndFormatValue(inputFormat, outputFormat, unformattedItem)
				if err != nil {
					return nil, err
				}
				// replacing the wildcard here feels hacky, but seems to be the
				// quickest way to achieve the outcome we want
				data, err = setJSONRaw(data, []byte(formattedItem), strings.Replace(k, "*", strconv.Itoa(idx), -1))
				if err != nil {
					return nil, err
				}
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
	formattedItem := strings.Join([]string{"\"", parsedItem.Format(outputFormat), "\""}, "")
	return formattedItem, nil
}
