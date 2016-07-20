// Package kazaam provides a simple interface for transforming arbitrary JSON in Golang.
package kazaam

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/bitly/go-simplejson"
)

// Kazaam includes internal data required for handling the transformation.
// A Kazaam object must be initialized using the NewKazaam function.
type Kazaam struct {
	spec     string
	specJSON *simplejson.Json
}

// NewKazaam creates a new Kazaam instance by parsing the `spec` argument as JSON and
// returns a pointer to it. The string `spec` must be valid JSON or empty for
// NewKazaam to return a Kazaam object.
//
// If empty, the default Kazaam behavior when the Transform variants are called is to
// return the original data unmodified.
//
// The only validation done at initialization time is to ensure that the `spec` is
// valid JSON. If the JSON is invalid, a nil Kazaam pointer and an explanation of
// the parse error is returned. The contents of the transform specification is
// validated only at Transform time.
func NewKazaam(spec string) (*Kazaam, error) {
	if len(spec) == 0 {
		spec = `[{"operation": "pass"}]`
	}

	json, err := simplejson.NewJson([]byte(spec))
	if err != nil {
		return nil, err
	}
	j := Kazaam{spec: spec, specJSON: json}

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
	for i := range j.specJSON.MustArray() {
		specObj := j.specJSON.GetIndex(i)
		spec := specObj.Get("spec")
		over, overExists := specObj.CheckGet("over")
		if overExists {
			overPath, _ := over.String()
			dataList := data.GetPath(strings.Split(overPath, ".")...).MustArray()

			var transformedDataList []*simplejson.Json
			for _, x := range dataList {
				jsonValue := simplejson.New()
				jsonValue.SetPath(nil, x)
				transformedData, intErr := getTransform(specObj)(spec, jsonValue)
				if intErr != nil {
					return data, err
				}
				transformedDataList = append(transformedDataList, transformedData)
			}
			data.SetPath(strings.Split(overPath, "."), transformedDataList)
		} else {
			data, err = getTransform(specObj)(spec, data)
		}
	}
	return data, err
}

// TransformJSONStringToString loads the JSON string `data`, transforms
// it as per the spec, and returns the transformed JSON string.
func (j *Kazaam) TransformJSONStringToString(data string) (string, error) {
	// read in the arbitrary input data
	d, _ := simplejson.NewJson([]byte(data))
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
func getTransform(specObj *simplejson.Json) func(*simplejson.Json, *simplejson.Json) (*simplejson.Json, error) {
	operation, _ := specObj.Get("operation").String()

	switch operation {
	case "shift":
		return transformShift
	case "default":
		return transformDefault
	default:
		return transformPass
	}
}

// no op transform -- useful for testing/default behavior
func transformPass(spec *simplejson.Json, data *simplejson.Json) (*simplejson.Json, error) {
	return data, nil
}

// simple transform to set default values in output json
func transformDefault(spec *simplejson.Json, data *simplejson.Json) (*simplejson.Json, error) {
	for k, v := range spec.MustMap() {
		data.SetPath(strings.Split(k, "."), v)
	}
	return data, nil
}

func transformShift(spec *simplejson.Json, data *simplejson.Json) (*simplejson.Json, error) {
	outData := simplejson.New()
	for k, v := range spec.MustMap() {
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
					return nil, fmt.Errorf("Warn: Unable to coerce element to json string: %v", vItem)
				}
				keyList = append(keyList, vItemStr)
			}
		default:
			return nil, fmt.Errorf("Warn: Unknown type in message for key: %s", k)
		}

		// iterate over keys to evaluate
		for _, v := range keyList {
			var dataForV *simplejson.Json
			var err error

			// grab the data
			if v == "$" {
				dataForV = data
			} else {
				dataForV, err = getJSONPath(data, v)
				if err != nil {
					return nil, err
				}
			}

			// if array flag set, encapsulate data
			if array {
				tmp, _ := dataForV.MarshalJSON()
				tmpString := "[" + string(tmp) + "]"
				dataForV, _ = simplejson.NewJson([]byte(tmpString))
			}

			outData.SetPath(outPath, dataForV.Interface())
		}
	}
	return outData, nil
}

var jsonPathRe = regexp.MustCompile("([^\\[\\]]+)\\[([0-9\\*]+)\\]")

func getJSONPath(j *simplejson.Json, path string) (*simplejson.Json, error) {
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
						intermediate, err := getJSONPath(jin.GetIndex(i), newPath)
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
			jin = jin.Get(k)
		}
	}
	return jin, nil
}
