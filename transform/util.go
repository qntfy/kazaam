// Package transform package contains canonical implementations of Kazaam transforms.
package transform

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mbordner/kazaam/registry"
	"regexp"
	"strconv"
	"strings"

	"github.com/qntfy/jsonparser"
)

// go test "github.com/qntfy/kazaam/transform" -coverprofile cover.out
// go tool cover -html=cover.out -o cover.html

// ParseError should be thrown when there is an issue with parsing any of the specification or data
type ParseError string

func (p ParseError) Error() string {
	return string(p)
}

// CPathNoDefaultError thrown when a path is missing, and marked conditional, and didn't have a default set
type CPathSkipError string

func (c CPathSkipError) Error() string {
	return string(c)
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

const (
	JSONNull = iota
	JSONString
	JSONInt
	JSONFloat
	JSONBool
)

var (
	ConditionalPathSkip = CPathSkipError("Conditional Path missing and without a default value")
	NonExistentPath     = RequireError("Path does not exist")
	jsonPathRe          = regexp.MustCompile(`([^\[\]]+)\[(.*?)\]`)
	leadingZeroRe       = regexp.MustCompile(`^(-)*(0+)([\.1-9])`)
	// matches converter groups separated by pipe delimiter characters (|), ignoring escaped delimiters
	// (odd number of \ chars proceeding the delimiter)
	// json path, e.g.  path.path?conditional default|converter1 args|converter2 args
	// allows also || to be ignored incase we need to support OR operators
	// phoneNumbers[?(@.type =||= "iPhone")].number ? blah.blah == "  blah:bla \" h  " : "b\"l\tah" |convert||\|\\\|er1 "args args" |converter2 \ args\|\\\| args |converter3 2 2
	pathConverterSplitRe = regexp.MustCompile(`(?:\|)?(?:[^\||\\]*(?:(?:\|\|)|(?:\\(?:\\\\)*.?))*)*[^\|]`) // with forward scanning: (?<=(?<!\\)(?:\\\\)*)\|

	// check if first token is conditional
	conditionalMatchRe = regexp.MustCompile(`(\w.*\w(?:\s*\[\s*\d+\s*\])?)(?:\s*\?\s*)(.*?)\s*$`) //`(\w.*\w)(?:\s*\?\s*)(.*?)\s*$`) // phoneNumbers[?(@.type == "iPhone")].number ? blah.blah==3:"blah"

	// used to parse out the parts of the condition (source, operator, value), and default value if they are provided
	//conditionMatchRe = regexp.MustCompile(`(\w.*\w)(?:\s*)(==|!=|>|<|>=|<=)(?:\s*)(\w*|".*:.*")(?:\s*):(?:\s*)(.*?)(?:\s*)$`) // blah.blah == "  blah:bla \" h  " : "blah"
	conditionMatchRe = regexp.MustCompile(`(?:\s*)(.*?)(?:\s*):(?:\s*)(\w*|".*")(?:\s*)$`) // blah.blah == "  blah:bla \" h  " : "b\"lah"

	// match the slashes, and the next character, to split and unescape characters
	unescapeTokensRe = regexp.MustCompile(`(?:[^\\]+)|\\(.)`)

	// used to parse out converter name and arguments
	converterParsingRe = regexp.MustCompile(`(?:^\|\s*)(\w+)(?:\s*)((?:.*?)(?:\\\s?)*)(?:\s*)$`) // | blah \ blah blah blah \
)

type JSONValue struct {
	valueType         int
	value             interface{}
	num               json.Number
	data              []byte
	floatStrPrecision int
}

func (v *JSONValue) GetType() int {
	return v.valueType
}

func (v *JSONValue) GetValue() interface{} {
	return v.value
}

func (v *JSONValue) GetStringValue() string {
	return v.value.(string)
}

func (v *JSONValue) GetQuotedStringValue() string {
	return strconv.Quote(v.GetStringValue())
}

func (v *JSONValue) GetIntValue() int64 {
	return v.value.(int64)
}

func (v *JSONValue) GetFloatValue() float64 {
	return v.value.(float64)
}

func (v *JSONValue) GetBoolValue() bool {
	return v.value.(bool)
}

func (v *JSONValue) GetNumber() json.Number {
	return v.num
}

func (v *JSONValue) IsNumber() bool {
	return v.valueType == JSONInt || v.valueType == JSONFloat
}

func (v *JSONValue) IsBool() bool {
	return v.valueType == JSONBool
}

func (v *JSONValue) IsString() bool {
	return v.valueType == JSONString
}

func (v *JSONValue) IsNull() bool {
	return v.valueType == JSONNull
}

func (v *JSONValue) GetData() []byte {
	return v.data
}

func (v *JSONValue) SetFloatStringPrecision(p int) {
	v.floatStrPrecision = p
}

func (v *JSONValue) getFloatStringFormat() string {
	if v.floatStrPrecision < 0 {
		return "%f"
	}
	return fmt.Sprintf("%%.%df", v.floatStrPrecision)
}

func (v *JSONValue) String() string {
	switch v.valueType {
	default:
		fallthrough
	case JSONNull:
		return "null"
	case JSONString:
		return strconv.Quote(v.GetStringValue())
	case JSONBool:
		return fmt.Sprintf("%t", v.GetBoolValue())
	case JSONInt:
		return fmt.Sprintf("%d", v.GetIntValue())
	case JSONFloat:
		return fmt.Sprintf(v.getFloatStringFormat(), v.GetFloatValue())
	}
}

// returns a Value (json value type) from the json data byte array, and string path
func GetJsonPathValue(jsonData []byte, path string) (value *JSONValue, err error) {
	var data []byte

	data, err = GetJSONRaw(jsonData, path, true)
	if err != nil {
		return
	}

	value, err = NewJSONValue(data)
	if err != nil {
		return
	}

	return
}

// returns a Value (json value type) from the byte array data for a value
func NewJSONValue(data []byte) (value *JSONValue, err error) {

	// remove leading zeros that are invalid in jaon if it's a number value
	tmp := string(data)
	if m := leadingZeroRe.FindStringSubmatch(tmp); m != nil {
		tmp = tmp[len(m[1])+len(m[2]):]
		if rune(m[3][0]) == '.' {
			m[2] = "0"
		} else {
			m[2] = ""
		}
		tmp = m[1] + m[2] + tmp
		data = []byte(tmp)
	}

	value = new(JSONValue)
	value.data = data
	value.floatStrPrecision = -1

	reader := bytes.NewReader(data)
	decoder := json.NewDecoder(reader)
	decoder.UseNumber()

	if err = decoder.Decode(&value.value); err != nil {
		return
	}

	switch value.value.(type) {
	case string:
		value.valueType = JSONString
	case bool:
		value.valueType = JSONBool
	case nil:
		value.valueType = JSONNull
	case json.Number:
		value.num = value.value.(json.Number)
		if strings.Contains(value.num.String(), ".") {
			value.valueType = JSONFloat
			value.value, err = value.num.Float64()
			if err != nil {
				return
			}
		} else {
			value.valueType = JSONInt
			value.value, err = value.num.Int64()
			if err != nil {
				return
			}
		}
	}

	return
}

// Config contains the options that dictate the behavior of a transform. The internal
// `spec` object can be an arbitrary json configuration for the transform.
type Config struct {
	Spec    *map[string]interface{} `json:"spec"`
	Require bool                    `json:"require,omitempty"`
	InPlace bool                    `json:"inplace,omitempty"`
}

// all it does is remove \ characters for now.
func unescapeString(s string) string {
	if matches := unescapeTokensRe.FindAllStringSubmatch(s, -1); matches != nil {
		tokens := make([]string, len(matches))
		for i, match := range matches {
			if match[0][0] == '\\' {
				tokens[i] = match[1]
			} else {
				tokens[i] = match[0]
			}
		}
		return strings.Join(tokens, "")
	}
	return s
}

type JSONPathConverter struct {
	converter string
	name      string
	arguments string
}

func NewJSONPathConverter(converter string) *JSONPathConverter { // ||stoi \ no \|
	params := new(JSONPathConverter)
	params.converter = converter
	if matches := converterParsingRe.FindStringSubmatch(converter); matches != nil {
		if len(matches) > 1 {
			params.name = matches[1]
			if len(matches) > 2 {
				params.arguments = strconv.Quote(unescapeString(matches[2])) // store as a quoted string, so that it will be unpacked as a string value in Convert
			}
		}
	}
	return params
}

func (converter *JSONPathConverter) isValid() bool {
	return len(converter.name) > 0
}
func (converter *JSONPathConverter) getName() string {
	return converter.name
}
func (converter *JSONPathConverter) getArguments() string {
	return converter.arguments
}

/*
  [jsonPath]?[conditional]|[converter]
	e.g. root.node1[1].property ? root.node0.prop1 == "true" | regex remove_commas

	the path will be broken into 3 components:
	1. the json parser path to the node value
    2. the conditional component
  	3. and a series of converter expressions to pipe the value through.

	The conditional component has 4 forms:

	root.prop1 ?
								// if the value exists, return it, otherwise skip it, regardless on
								// whether paths are required.
	root.prop1 ? defaultVal
								// if value exists, use that, otherwise return the default value. this is JSON syntax
								// so strings, will require double quotes.
	root.prop1 ? root.node1.prop2 == true :
								// if value exists, and the expression is true, return the existing value
								// if the value exists, and the expression is false, skip the existing value
								// note that the : is required here to end the expression.
	root.prop1 ? root.node1.prop2 == true : defaultValue
								// if the value exists, and the expression is true, return the existing value
								// if the value exists and the expression is false,
								// return the default value (JSON syntax)

	Conditional Expressions support  () , <, >, >=, <=, ==, !=, !, true, false, && and ||
	Conditional expressions must evaluate to a boolean true or false value.
	Identifiers in the expression are assumed to be JSON paths within the document, and will evaluate
	to their current value.  Only non composite JSON values are supported: boolean, number,
	string (i.e. not arrays or objects).

	Function calls of the form  <converter name>( json.path, "converter arguments")   are also supported, e.g.
		root.prop1 ? ston(root.node1.prop1) == 3 && regex(root.node2.prop2,"remove_commas) == "1000" :

	The converters component will define a series of value conversions
*/
type JSONPathParameters struct {
	data         []byte
	originalPath string // original path string
	jsonPath     string // parsed out and trimmed json path

	condition           *BasicExpr
	conditional         bool
	defaultValue        []byte
	conditionParseError bool

	converters []*JSONPathConverter
}

func NewJSONPathParameters(data []byte, path string) *JSONPathParameters {
	jsonPathParams := new(JSONPathParameters)

	jsonPathParams.data = data
	jsonPathParams.originalPath = path

	// this is basically parsing out the converter tokens, but the first token will be the
	// path along with any conditional/default value markup
	tokens := pathConverterSplitRe.FindAllString(path, -1)

	// check if the first token contains a conditional declaration,   path.path?conditional default
	// when ? is after the path, it means it's conditional on the path existing, and will be skipped
	// otherwise
	if matches := conditionalMatchRe.FindStringSubmatch(tokens[0]); matches != nil {
		jsonPathParams.conditional = true

		jsonPathParams.jsonPath = matches[1]

		if len(matches[2]) > 0 {
			if conditionMatchRe.MatchString(matches[2]) { // looking for condition : default value

				if matches := conditionMatchRe.FindStringSubmatch(matches[2]); matches != nil {

					jsonPathParams.defaultValue = []byte(unescapeString(matches[2]))

					expr, e := NewBasicExpr(data, matches[1])
					if e == nil {
						jsonPathParams.condition = expr
					} else {
						jsonPathParams.conditionParseError = true
					}
				}

			} else { // or just default value
				jsonPathParams.defaultValue = []byte(unescapeString(matches[2]))
			}
		}
	} else {
		jsonPathParams.jsonPath = strings.Trim(tokens[0], " \t")
	}

	jsonPathParams.converters = make([]*JSONPathConverter, 0, len(tokens)-1)

	// parse out the converter
	if len(tokens) > 1 {
		for _, token := range tokens[1:] {
			converter := NewJSONPathConverter(token)
			if converter.isValid() {
				jsonPathParams.converters = append(jsonPathParams.converters, converter)
			}
		}
	}

	return jsonPathParams
}

func (params *JSONPathParameters) getJsonPath() string {
	return params.jsonPath
}

// means that ? existed in the json path
func (params *JSONPathParameters) isConditional() bool {
	return params.conditional
}

// means that there was an expression in the conditional component of the path, e.g.  path ? <expression> :
func (params *JSONPathParameters) hasConditionExpression() bool {
	return params.condition != nil || params.conditionParseError
}

func (params *JSONPathParameters) evalConditionExpression() (results bool, err error) {
	if params.conditionParseError == false {
		return params.condition.Eval()
	}
	err = errors.New("condition parse error")
	return
}

// has a value after the ? in the conditional component, or a value after the : if there is a conditional expression
func (params *JSONPathParameters) hasDefaultValue() bool {
	if params.defaultValue != nil && len(params.defaultValue) > 0 {
		_, err := NewJSONValue([]byte(params.defaultValue))
		if err == nil {
			return true
		}
	}
	return false
}

func (params *JSONPathParameters) getDefaultValue() ([]byte, error) {
	if params.hasDefaultValue() {
		return params.defaultValue, nil
	}
	return nil, errors.New("invalid or non existent default value")
}

// converters will modify the value, and can be chained
func (params *JSONPathParameters) hasConverters() bool {
	return len(params.converters) > 0
}

func (params *JSONPathParameters) convert(value []byte) (newValue []byte, err error) {
	newValue = value
	for _, c := range params.converters {
		converter := registry.GetConverter(c.getName())
		var args []byte
		if len(c.getArguments()) > 0 {
			args = []byte(c.getArguments())
		}
		newValue, err = converter.Convert(params.data, newValue, args)
		if err != nil {
			return
		}
	}
	return
}

func GetJSONRaw(data []byte, path string, pathRequired bool) ([]byte, error) {
	return getJSONRaw(data, path, pathRequired)
}

// Given a json byte slice `data` and a kazaam `path` string, return the object at the path in data if it exists.
func getJSONRaw(data []byte, path string, pathRequired bool) ([]byte, error) {
	jsonPathParams := NewJSONPathParameters(data, path)
	return getProcessedJSONRaw(data, jsonPathParams.getJsonPath(), pathRequired, jsonPathParams)
}

func getProcessedJSONRaw(data []byte, path string, pathRequired bool, params *JSONPathParameters) ([]byte, error) {
	objectKeys := strings.Split(path, ".")
	numOfInserts := 0
	for element, k := range objectKeys {
		// check the object key to see if it also contains an array reference
		arrayRefs := jsonPathRe.FindAllStringSubmatch(k, -1)
		if arrayRefs != nil && len(arrayRefs) > 0 {
			objKey := arrayRefs[0][1]      // the key
			arrayKeyStr := arrayRefs[0][2] // the array index
			err := validateArrayKeyString(arrayKeyStr)
			if err != nil {
				return nil, err
			}
			// if there's a wildcard array reference
			if arrayKeyStr == "*" {
				// ArrayEach setup
				objectKeys[element+numOfInserts] = objKey
				beforePath := objectKeys[:element+numOfInserts+1]
				newPath := strings.Join(objectKeys[element+numOfInserts+1:], ".")
				var results [][]byte

				// use jsonparser.ArrayEach to copy the array into results
				_, err := jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
					results = append(results, HandleUnquotedStrings(value, dataType))
				}, beforePath...)
				if err == jsonparser.KeyPathNotFoundError {
					if pathRequired {
						return nil, NonExistentPath
					}
				} else if err != nil {
					return nil, err
				}

				// GetJSONRaw() the rest of path for each element in results
				if newPath != "" {
					for i, value := range results {
						intermediate, err := getProcessedJSONRaw(value, newPath, pathRequired, params)
						if err == jsonparser.KeyPathNotFoundError {
							if pathRequired {
								return nil, NonExistentPath
							}
						} else if err != nil {
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
			// separate the array key as a new element in objectKeys
			objectKeys = makePathWithIndex(arrayKeyStr, objKey, objectKeys, element+numOfInserts)
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
		result = HandleUnquotedStrings(result, dataType)
	}
	if len(result) == 0 {
		result = []byte("null")
	}
	if err == jsonparser.KeyPathNotFoundError {
		if params.isConditional() {
			if params.hasDefaultValue() {
				result, _ = params.getDefaultValue()
			} else {
				return nil, ConditionalPathSkip
			}
		} else if pathRequired {
			return nil, NonExistentPath
		}
	} else if params.isConditional() && params.hasConditionExpression() {
		// path exists, but there is a conditional expression
		evaluation, err := params.evalConditionExpression()
		if err == nil {

			if !evaluation {
				if params.hasDefaultValue() {
					result, _ = params.getDefaultValue()
				} else {
					return nil, ConditionalPathSkip
				}
			}

		} else if params.hasDefaultValue() {
			result, _ = params.getDefaultValue()
		} else {
			return nil, ConditionalPathSkip // because there was an error, it should be evaluated as false
		}

	} else if err != nil {
		return nil, err
	}

	if params.hasConverters() {
		result, err = params.convert(result)
	}

	return result, nil
}

func SetJSONRaw(data, out []byte, path string) ([]byte, error) {
	return setJSONRaw(data, out, path)
}

// setJSONRaw sets the value at a key and handles array indexing
func setJSONRaw(data, out []byte, path string) ([]byte, error) {
	var err error
	splitPath := strings.Split(path, ".")
	numOfInserts := 0

	for element, k := range splitPath {
		arrayRefs := jsonPathRe.FindAllStringSubmatch(k, -1)
		if arrayRefs != nil && len(arrayRefs) > 0 {
			objKey := arrayRefs[0][1]      // the key
			arrayKeyStr := arrayRefs[0][2] // the array index
			err = validateArrayKeyString(arrayKeyStr)
			if err != nil {
				return nil, err
			}
			// Note: this branch of the function is not currently used by any
			// existing transforms. It is simpy here to support he generalized
			// form of this operation
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
			splitPath = makePathWithIndex(arrayKeyStr, objKey, splitPath, element+numOfInserts)
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

func DelJSONRaw(data []byte, path string, pathRequired bool) ([]byte, error) {
	return delJSONRaw(data, path, pathRequired)
}

// delJSONRaw deletes the value at a path and handles array indexing
func delJSONRaw(data []byte, path string, pathRequired bool) ([]byte, error) {
	var err error
	splitPath := strings.Split(path, ".")
	numOfInserts := 0

	for element, k := range splitPath {
		arrayRefs := jsonPathRe.FindAllStringSubmatch(k, -1)
		if arrayRefs != nil && len(arrayRefs) > 0 {
			objKey := arrayRefs[0][1]      // the key
			arrayKeyStr := arrayRefs[0][2] // the array index
			err = validateArrayKeyString(arrayKeyStr)
			if err != nil {
				return nil, err
			}

			// not currently supported
			if arrayKeyStr == "*" {
				return nil, SpecError("Array wildcard not supported for this operation.")
			}

			// if not a wildcard then piece that path back together with the
			// array index as an entry in the splitPath slice
			splitPath = makePathWithIndex(arrayKeyStr, objKey, splitPath, element+numOfInserts)
			numOfInserts++
		} else {
			// no array reference, good to go
			continue
		}
	}

	if pathRequired {
		_, _, _, err = jsonparser.Get(data, splitPath...)
		if err == jsonparser.KeyPathNotFoundError {
			return nil, NonExistentPath
		} else if err != nil {
			return nil, err
		}
	}

	data = jsonparser.Delete(data, splitPath...)
	return data, nil
}

// validateArrayKeyString is a helper function to make sure the array index is
// legal
func validateArrayKeyString(arrayKeyStr string) error {
	if arrayKeyStr != "*" && arrayKeyStr != "+" && arrayKeyStr != "-" {
		val, err := strconv.Atoi(arrayKeyStr)
		if val < 0 || err != nil {
			return ParseError(fmt.Sprintf("Warn: Unable to coerce index to integer: %v", arrayKeyStr))
		}
	}
	return nil
}

// makePathWithIndex generats a path slice to pass to jsonparser
func makePathWithIndex(arrayKeyStr, objectKey string, pathSlice []string, pathIndex int) []string {
	arrayKey := string(bookend([]byte(arrayKeyStr), '[', ']'))
	pathSlice[pathIndex] = objectKey
	pathSlice = append(pathSlice, "")
	copy(pathSlice[pathIndex+2:], pathSlice[pathIndex+1:])
	pathSlice[pathIndex+1] = arrayKey
	return pathSlice
}

// add characters at beginning and end of []byte
func bookend(value []byte, bef, aft byte) []byte {
	value = append(value, ' ', aft)
	copy(value[1:], value[:len(value)-2])
	value[0] = bef
	return value
}

// jsonparser strips quotes from returned strings, this adds them back
func HandleUnquotedStrings(value []byte, dt jsonparser.ValueType) []byte {
	if dt == jsonparser.String {
		// bookend() is destructive to underlying slice, need to copy.
		// extra capacity saves an allocation and copy during bookend.
		tmp := make([]byte, len(value), len(value)+2)
		copy(tmp, value)
		value = bookend(tmp, '"', '"')
	}
	return value
}
