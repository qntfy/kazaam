// Package kazaam provides a simple interface for transforming arbitrary JSON in Golang.
package kazaam

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/bitly/go-simplejson"
)

type transformFunc func(spec *spec, data *simplejson.Json) (*simplejson.Json, error)

var validSpecTypes map[string]transformFunc

func init() {
	validSpecTypes = map[string]transformFunc{
		"pass":     transformPass,
		"shift":    transformShift,
		"extract":  transformExtract,
		"default":  transformDefault,
		"concat":   transformConcat,
		"coalesce": transformCoalesce,
	}
}

// spec represent each individual spec element
type spec struct {
	Operation *string                 `json:"operation"`
	Spec      *map[string]interface{} `json:"spec"`
	Over      *string                 `json:"over,omitempty"`
	Require   bool                    `json:"require,omitempty"`
}

type KazaamError struct {
	ErrMsg, ErrType string
}

func (e *KazaamError) Error() string {
	return fmt.Sprintf("%s - %s", e.ErrType, e.ErrMsg)
}

type specInt spec
type specs []spec

// UnmarshalJSON implements a custon unmarshaller for the Spec type
func (s *spec) UnmarshalJSON(b []byte) (err error) {
	j := specInt{}
	if err = json.Unmarshal(b, &j); err == nil {
		*s = spec(j)
		if s.Operation == nil {
			err = errors.New("Spec must contain an \"operation\" field")
			return
		}
		if _, ok := validSpecTypes[*s.Operation]; ok == false {
			err = errors.New("Invalid spec operation specified")
		}
		if s.Spec != nil && len(*s.Spec) < 1 {
			err = &KazaamError{ErrMsg: "Spec must contain at least one element", ErrType: "SpecError"}
			return
		}
		return
	}
	return
}

// Kazaam includes internal data required for handling the transformation.
// A Kazaam object must be initialized using the NewKazaam function.
type Kazaam struct {
	spec     string
	specJSON specs
}

// NewKazaam creates a new Kazaam instance by parsing the `spec` argument as JSON and
// returns a pointer to it. The string `spec` must be valid JSON or empty for
// NewKazaam to return a Kazaam object.
//
// If empty, the default Kazaam behavior when the Transform variants are called is to
// return the original data unmodified.
//
// At initialization time, the `spec` is checked to ensure that it is
// valid JSON. Further, it confirms that all individual specs have a properly-specified
// `operation` and details are set if required. If the spec is invalid, a nil Kazaam
// pointer and an explanation of the error is returned. The contents of the transform
// specification is further validated at Transform time.
func NewKazaam(specString string) (*Kazaam, error) {
	if len(specString) == 0 {
		specString = `[{"operation":"pass"}]`
	}
	var specElements specs
	if err := json.Unmarshal([]byte(specString), &specElements); err != nil {
		return nil, err
	}

	j := Kazaam{spec: specString, specJSON: specElements}

	return &j, nil
}

// Transform takes the *simplejson.Json `data`, transforms it according
// to the loaded spec, and returns the modified *simplejson.Json object.
//
// Note: this is a destructive operation: the transformation is done in place.
// You must perform a deep copy of the data prior to calling Transform if
// the original JSON object must be retained.
func (j *Kazaam) Transform(data *simplejson.Json) (*simplejson.Json, error) {
	var err error
	for _, specObj := range j.specJSON {
		//spec := specObj.Get("spec")
		//over, overExists := specObj.CheckGet("over")
		if specObj.Over != nil {
			dataList := data.GetPath(strings.Split(*specObj.Over, ".")...).MustArray()

			var transformedDataList []interface{}
			for _, x := range dataList {
				jsonValue := simplejson.New()
				jsonValue.SetPath(nil, x)
				transformedData, intErr := getTransform(specObj)(&specObj, jsonValue)
				if intErr != nil {
					return data, err
				}
				transformedDataList = append(transformedDataList, transformedData.Interface())
			}
			data.SetPath(strings.Split(*specObj.Over, "."), transformedDataList)

		} else {
			data, err = getTransform(specObj)(&specObj, data)
		}
	}
	return data, err
}

// TransformJSONStringToString loads the JSON string `data`, transforms
// it as per the spec, and returns the transformed JSON string.
func (j *Kazaam) TransformJSONStringToString(data string) (string, error) {
	// read in the arbitrary input data
	d, err := simplejson.NewJson([]byte(data))
	if err != nil {
		return "", err
	}
	outputJSON, err := j.Transform(d)
	if err != nil {
		return "", err
	}
	jsonString, err := outputJSON.MarshalJSON()
	if err != nil {
		return "", err
	}
	return string(jsonString), nil
}

// TransformJSONString loads the JSON string, transforms it as per the
// spec, and returns a pointer to a transformed simplejson.Json.
//
// This function is especially useful when one may need to extract
// multiple fields from the transformed JSON.
func (j *Kazaam) TransformJSONString(data string) (*simplejson.Json, error) {
	// read in the arbitrary input data
	d, err := simplejson.NewJson([]byte(data))
	if err != nil {
		return nil, err
	}
	return j.Transform(d)
}

// return the transform function based on what's indicated in the operation spec
func getTransform(specObj spec) transformFunc {
	tform, ok := validSpecTypes[*specObj.Operation]
	if ok == false {
		return transformPass
	}
	return tform
}

// no op transform -- useful for testing/default behavior
func transformPass(spec *spec, data *simplejson.Json) (*simplejson.Json, error) {
	return data, nil
}

// simple transform to set default values in output json
func transformDefault(spec *spec, data *simplejson.Json) (*simplejson.Json, error) {
	for k, v := range *spec.Spec {
		data.SetPath(strings.Split(k, "."), v)
	}
	return data, nil
}

func transformExtract(spec *spec, data *simplejson.Json) (*simplejson.Json, error) {
	outPath, ok := (*spec.Spec)["path"]
	if !ok {
		return nil, &KazaamError{ErrMsg: fmt.Sprintf("Unable to get path"), ErrType: "SpecError"}
	}
	outData, err := getJSONPath(data, outPath.(string), spec.Require)
	return outData, err
}

func transformShift(spec *spec, data *simplejson.Json) (*simplejson.Json, error) {
	outData := simplejson.New()
	for k, v := range *spec.Spec {
		array := true
		outPath := strings.Split(k, ".")

		var keyList []string

		// check if `v` is a string or list and build a list of keys to evaluate
		switch v.(type) {
		case string:
			keyList = append(keyList, v.(string))
			array = false
		case []interface{}:
			for _, vItem := range v.([]interface{}) {
				vItemStr, found := vItem.(string)
				if !found {
					return nil, &KazaamError{ErrMsg: fmt.Sprintf("Warn: Unable to coerce element to json string: %v", vItem), ErrType: "ParseError"}
				}
				keyList = append(keyList, vItemStr)
			}
		default:
			return nil, &KazaamError{ErrMsg: fmt.Sprintf("Warn: Unknown type in message for key: %s", k), ErrType: "ParseError"}
		}

		// iterate over keys to evaluate
		for _, v := range keyList {
			var dataForV *simplejson.Json
			var err error

			// grab the data
			if v == "$" {
				dataForV = data
			} else {
				dataForV, err = getJSONPath(data, v, spec.Require)
				if err != nil {
					return nil, err
				}
			}

			// if array flag set, encapsulate data
			if array {
				var intSlice = make([]interface{}, 1)
				intSlice[0] = dataForV.Interface()
				dataForV.SetPath(nil, intSlice)
			}

			outData.SetPath(outPath, dataForV.Interface())
		}
	}
	return outData, nil
}

func transformConcat(spec *spec, data *simplejson.Json) (*simplejson.Json, error) {
	sourceList, sourceOk := (*spec.Spec)["sources"]
	if !sourceOk {
		return nil, &KazaamError{ErrMsg: fmt.Sprintf("Unable to get sources"), ErrType: "SpecError"}
	}
	targetPath, targetOk := (*spec.Spec)["targetPath"]
	if !targetOk {
		return nil, &KazaamError{ErrMsg: fmt.Sprintf("Unable to get targetPath"), ErrType: "SpecError"}
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
					return nil, &KazaamError{ErrMsg: fmt.Sprintf("Path does not exist"), ErrType: "RequireError"}
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
				return nil, &KazaamError{ErrMsg: fmt.Sprintf("Error processing %v: must have either value or path specified", vItem), ErrType: "SpecError"}
			}
		}
		outString += value.(string)

		applyDelim = true
	}

	data.SetPath(strings.Split(targetPath.(string), "."), outString)

	return data, nil
}

func transformCoalesce(spec *spec, data *simplejson.Json) (*simplejson.Json, error) {
	if spec.Require == true {
		return nil, &KazaamError{ErrMsg: fmt.Sprintf("Invalid spec. Coalesce does not support \"require\""), ErrType: "SpecError"}
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
					return nil, &KazaamError{ErrMsg: fmt.Sprintf("Warn: Unable to coerce element to json string: %v", vItem), ErrType: "ParseError"}
				}
				keyList = append(keyList, vItemStr)
			}
		default:
			return nil, &KazaamError{ErrMsg: fmt.Sprintf("Warn: Expected list in message for key: %s", k), ErrType: "ParseError"}
		}

		// iterate over keys to evaluate
		for _, v := range keyList {
			var dataForV *simplejson.Json
			var err error

			// grab the data
			dataForV, err = getJSONPath(data, v, false)
			if err != nil {
				return nil, err
			}
			if dataForV.Interface() != nil {
				data.SetPath(outPath, dataForV.Interface())
				break
			}
		}
	}
	return data, nil
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
				if pathRequired == true {
					jin, exists := jin.CheckGet(objKey)
					if exists != true {
						return jin, &KazaamError{ErrMsg: fmt.Sprintf("Path does not exist"), ErrType: "RequireError"}
					}
				} else {
					jin = jin.Get(objKey)
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
			if pathRequired == true {
				_, exists := jin.CheckGet(k)
				if exists != true {
					return nil, &KazaamError{ErrMsg: fmt.Sprintf("Path does not exist"), ErrType: "RequireError"}
				}
			}
			jin = jin.Get(k)
		}
	}
	return jin, nil
}
