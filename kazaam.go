// Package kazaam provides a simple interface for transforming arbitrary JSON in Golang.
package kazaam

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/bitly/go-simplejson"
	"github.com/qntfy/kazaam/transform"
)

// TransformFunc defines the contract that any Transform function implementation
// must abide by. The transform's first argument is a `kazaam.Spec` object that
// contains any configuration necessary for the transform. The second argument
// is a simplejson object that contains the data to be transformed.
//
// The function should return the transformed data, and an error if necessary.
// Transforms should strive to fail gracefully whenever possible.
type TransformFunc func(spec *transform.Config, data *simplejson.Json) (*simplejson.Json, error)

var validSpecTypes map[string]TransformFunc

func init() {
	validSpecTypes = map[string]TransformFunc{
		"pass":     transform.Pass,
		"shift":    transform.Shift,
		"extract":  transform.Extract,
		"default":  transform.Default,
		"concat":   transform.Concat,
		"coalesce": transform.Coalesce,
	}
}

// Spec represents an individual spec element. It describes the name of the operation,
// whether the `over` operator is required, whether any paths are required, and an
// operation-specific `Spec` that describes the configuration of the operation
type Spec struct {
	*transform.Config
	Operation *string `json:"operation"`
	Over      *string `json:"over,omitempty"`
}

type specInt Spec
type specs []Spec

// UnmarshalJSON implements a custon unmarshaller for the Spec type
func (s *Spec) UnmarshalJSON(b []byte) (err error) {
	j := specInt{}
	if err = json.Unmarshal(b, &j); err == nil {
		*s = Spec(j)
		if s.Operation == nil {
			err = &transform.Error{ErrMsg: "Spec must contain an \"operation\" field", ErrType: transform.SpecError}
			return
		}
		if _, ok := validSpecTypes[*s.Operation]; ok == false {
			err = &transform.Error{ErrMsg: "Invalid spec operation specified", ErrType: transform.SpecError}
		}
		if s.Config != nil && s.Spec != nil && len(*s.Spec) < 1 {
			err = &transform.Error{ErrMsg: "Spec must contain at least one element", ErrType: transform.SpecError}
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
	if j == nil {
		return nil, errors.New("Improperly initialized Kazaam object")
	}
	for _, specObj := range j.specJSON {
		//spec := specObj.Get("spec")
		//over, overExists := specObj.CheckGet("over")
		if specObj.Config != nil && specObj.Over != nil {
			dataList := data.GetPath(strings.Split(*specObj.Over, ".")...).MustArray()

			var transformedDataList []interface{}
			for _, x := range dataList {
				jsonValue := simplejson.New()
				jsonValue.SetPath(nil, x)
				transformedData, intErr := getTransform(specObj)(specObj.Config, jsonValue)
				if intErr != nil {
					return data, err
				}
				transformedDataList = append(transformedDataList, transformedData.Interface())
			}
			data.SetPath(strings.Split(*specObj.Over, "."), transformedDataList)

		} else {
			data, err = getTransform(specObj)(specObj.Config, data)
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
func getTransform(specObj Spec) TransformFunc {
	tform, ok := validSpecTypes[*specObj.Operation]
	if !ok {
		return transform.Pass
	}
	return tform
}
