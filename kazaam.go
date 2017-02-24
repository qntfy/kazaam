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
// is a `simplejson.Json` object that contains the data to be transformed.
//
// The data object passed in should be modified in-place. Where that is not
// possible, a new `simplejson.Json` object should be created and the pointer
// updated. The function should return an error if necessary.
// Transforms should strive to fail gracefully whenever possible.
type TransformFunc func(spec *transform.Config, data *simplejson.Json) error

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

// Config is used to configure a Kazaam Transformer object. Note: a manually-initialized
// config object (not created with `NewDefaultConfig`) will be UNAWARE of the built-in
// Kazaam transforms. Built-in and third-party Kazaam transforms will have to be
// manually registered for Kazaam to be able to transform data.
type Config struct {
	transforms map[string]TransformFunc
}

// NewDefaultConfig returns a properly initialized Config object that contains
// required mappings for all the built-in transform types.
func NewDefaultConfig() Config {
	// make a copy, otherwise if new transforms are registered, they'll affect the whole package
	specTypes := make(map[string]TransformFunc)
	for k, v := range validSpecTypes {
		specTypes[k] = v
	}
	return Config{transforms: specTypes}
}

// RegisterTransform registers a new transform type that satisfies the TransformFunc
// signature within the Kazaam configuration with the provided name. This function
// enables end-users to create and use custom transforms within Kazaam.
func (c *Config) RegisterTransform(name string, function TransformFunc) error {
	_, ok := c.transforms[name]
	if ok {
		return errors.New("Transform with that name already registered")
	}
	c.transforms[name] = function
	return nil
}

// Kazaam includes internal data required for handling the transformation.
// A Kazaam object must be initialized using the `New` or `NewKazaam` functions.
type Kazaam struct {
	spec     string
	specJSON specs
	config   Config
}

// NewKazaam creates a new Kazaam instance with a default configuration. See
// documentation for `New` for complete details.
func NewKazaam(specString string) (*Kazaam, error) {
	return New(specString, NewDefaultConfig())
}

// New creates a new Kazaam instance by parsing the `spec` argument as JSON and returns a
// pointer to it. Thew string `spec` must be valid JSON or empty for `New` to return
// a Kazaam object. This function also accepts a `Config` object used for modifying the
// behavior of the Kazaam Transformer.
//
// If `spec` is an empty string, the default Kazaam behavior when the Transform variants
// are called is to return the original data unmodified.
//
// At initialization time, the `spec` is checked to ensure that it is
// valid JSON. Further, it confirms that all individual specs have a properly-specified
// `operation` and details are set if required. If the spec is invalid, a nil Kazaam
// pointer and an explanation of the error is returned. The contents of the transform
// specification is further validated at Transform time.
//
// Currently, the Config object allows end users to register additional transform types
// to support performing custom transformations not supported by the canonical set of
// transforms shipped with Kazaam.
func New(specString string, config Config) (*Kazaam, error) {
	if len(specString) == 0 {
		specString = `[{"operation":"pass"}]`
	}
	var specElements specs
	if err := json.Unmarshal([]byte(specString), &specElements); err != nil {
		return nil, err
	}
	// do a check here to ensure all spec types are known
	for _, s := range specElements {
		if _, ok := config.transforms[*s.Operation]; !ok {
			return nil, &transform.Error{ErrMsg: "Invalid spec operation specified", ErrType: transform.SpecError}
		}
	}

	j := Kazaam{spec: specString, specJSON: specElements, config: config}

	return &j, nil
}

// return the transform function based on what's indicated in the operation spec
func (k *Kazaam) getTransform(s *spec) TransformFunc {
	// getting a non-existent transform is checked against before this function is
	// called, hence the _
	tform, _ := k.config.transforms[*s.Operation]
	return tform
}

// Transform takes the *simplejson.Json `data`, transforms it according
// to the loaded spec, and returns the modified *simplejson.Json object.
//
// Note: this is a destructive operation: the transformation is done in place.
// You must perform a deep copy of the data prior to calling Transform if
// the original JSON object must be retained.
func (k *Kazaam) Transform(data *simplejson.Json) (*simplejson.Json, error) {
	var err error
	for _, specObj := range k.specJSON {
		//spec := specObj.Get("spec")
		//over, overExists := specObj.CheckGet("over")
		if specObj.Config != nil && specObj.Over != nil {
			dataList := data.GetPath(strings.Split(*specObj.Over, ".")...).MustArray()

			var transformedDataList []interface{}
			for _, x := range dataList {
				jsonValue := simplejson.New()
				jsonValue.SetPath(nil, x)
				intErr := k.getTransform(&specObj)(specObj.Config, jsonValue)
				if intErr != nil {
					return data, err
				}
				transformedDataList = append(transformedDataList, jsonValue.Interface())
			}
			data.SetPath(strings.Split(*specObj.Over, "."), transformedDataList)

		} else {
			err = k.getTransform(&specObj)(specObj.Config, data)
		}
	}
	return data, err
}

// TransformJSONStringToString loads the JSON string `data`, transforms
// it as per the spec, and returns the transformed JSON string.
func (k *Kazaam) TransformJSONStringToString(data string) (string, error) {
	// read in the arbitrary input data
	d, err := simplejson.NewJson([]byte(data))
	if err != nil {
		return "", err
	}
	outputJSON, err := k.Transform(d)
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
func (k *Kazaam) TransformJSONString(data string) (*simplejson.Json, error) {
	// read in the arbitrary input data
	d, err := simplejson.NewJson([]byte(data))
	if err != nil {
		return nil, err
	}
	return k.Transform(d)
}
