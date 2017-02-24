// Package kazaam provides a simple interface for transforming arbitrary JSON in Golang.
package kazaam

import (
	"encoding/json"
	"fmt"
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

// Error provids an error mess (ErrMsg) and integer code (ErrType) for
// errors thrown by a transform
type Error struct {
	ErrMsg  string
	ErrType int
}

const (
	// ParseError is thrown when there is a JSON parsing error
	ParseError int = iota
	// RequireError is thrown when the JSON path does not exist and is required
	RequireError int = iota
	// SpecError is thrown when the kazaam specification is not properly formatted
	SpecError int = iota
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

type specInt spec
type specs []spec

// UnmarshalJSON implements a custon unmarshaller for the Spec type
func (s *spec) UnmarshalJSON(b []byte) (err error) {
	j := specInt{}
	if err = json.Unmarshal(b, &j); err == nil {
		*s = spec(j)
		if s.Operation == nil {
			err = &Error{ErrMsg: "Spec must contain an \"operation\" field", ErrType: SpecError}
			return
		}
		if _, ok := validSpecTypes[*s.Operation]; ok == false {
			err = &Error{ErrMsg: "Invalid spec operation specified", ErrType: SpecError}
		}
		if s.Spec != nil && len(*s.Spec) < 1 {
			err = &Error{ErrMsg: "Spec must contain at least one element", ErrType: SpecError}
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
	if !ok {
		return transformPass
	}
	return tform
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
