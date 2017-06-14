// Package kazaam provides a simple interface for transforming arbitrary JSON in Golang.
package kazaam

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/buger/jsonparser"
	"github.com/qntfy/kazaam/transform"
)

// TransformFunc defines the contract that any Transform function implementation
// must abide by. The transform's first argument is a `kazaam.Spec` object that
// contains any configuration necessary for the transform. The second argument
// is a `[]byte` object that contains the data to be transformed.
//
// The data object passed in should be modified in-place and returned. Where
// that is not possible, a new `[]byte` object should be created and returned.
// The function should return an error if necessary.
// Transforms should strive to fail gracefully whenever possible.
type TransformFunc func(spec *transform.Config, data []byte) ([]byte, error)

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

// Error provides an error message (ErrMsg) and integer code (ErrType) for
// errors thrown during the execution of a transform
type Error struct {
	ErrMsg  string
	ErrType int
}

const (
	// ParseError is thrown when there is a JSON parsing error
	ParseError = iota
	// RequireError is thrown when the JSON path does not exist and is required
	RequireError
	// SpecError is thrown when the kazaam specification is not properly formatted
	SpecError
)

// Error returns a string representation of the Error
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
			return nil, &Error{ErrMsg: "Invalid spec operation specified", ErrType: SpecError}
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

func transformErrorType(err error) error {
	switch err.(type) {
	case transform.ParseError:
		return &Error{ErrMsg: err.Error(), ErrType: ParseError}
	case transform.RequireError:
		return &Error{ErrMsg: err.Error(), ErrType: RequireError}
	case transform.SpecError:
		return &Error{ErrMsg: err.Error(), ErrType: SpecError}
	default:
		return err
	}
}

// Transform makes a copy of the byte slice `data`, transforms it according
// to the loaded spec, and returns the new, modified byte slice.
func (k *Kazaam) Transform(data []byte) ([]byte, error) {
	d := make([]byte, len(data))
	copy(d, data)
	d, err := k.TransformInPlace(d)
	return d, err
}

// TransformInPlace takes the byte slice `data`, transforms it according
// to the loaded spec, and modifies the byte slice in place.
//
// Note: this is a destructive operation: the transformation is done in place.
// You must perform a deep copy of the data prior to calling TransformInPlace if
// the original JSON object must be retained.
func (k *Kazaam) TransformInPlace(data []byte) ([]byte, error) {
	if k == nil || k.specJSON == nil {
		return data, &Error{ErrMsg: "Kazaam not properly initialized", ErrType: SpecError}
	}
	if len(data) == 0 {
		return data, nil
	}

	var err error
	for _, specObj := range k.specJSON {
		if specObj.Config != nil && specObj.Over != nil {
			var transformedDataList [][]byte
			_, err = jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
				transformedDataList = append(transformedDataList, value)
			}, strings.Split(*specObj.Over, ".")...)
			if err != nil {
				return data, transformErrorType(err)
			}
			for i, value := range transformedDataList {
				x := make([]byte, len(value))
				copy(x, value)
				x, intErr := k.getTransform(&specObj)(specObj.Config, x)
				if intErr != nil {
					return data, transformErrorType(err)
				}
				transformedDataList[i] = x
			}
			// copy into raw []byte format and return
			var buffer bytes.Buffer
			buffer.WriteByte('[')
			for i := 0; i < len(transformedDataList)-1; i++ {
				buffer.Write(transformedDataList[i])
				buffer.WriteByte(',')
			}
			if len(transformedDataList) > 0 {
				buffer.Write(transformedDataList[len(transformedDataList)-1])
			}
			buffer.WriteByte(']')
			data, err = jsonparser.Set(data, buffer.Bytes(), strings.Split(*specObj.Over, ".")...)
			if err != nil {
				return data, transformErrorType(err)
			}

		} else {
			data, err = k.getTransform(&specObj)(specObj.Config, data)
			if err != nil {
				return data, transformErrorType(err)
			}
		}
	}
	return data, transformErrorType(err)
}

// TransformJSONStringToString loads the JSON string `data`, transforms
// it as per the spec, and returns the transformed JSON string.
func (k *Kazaam) TransformJSONStringToString(data string) (string, error) {
	d, err := k.TransformJSONString(data)
	if err != nil {
		return "", err
	}
	return string(d), nil
}

// TransformJSONString loads the JSON string, transforms it as per the
// spec, and returns a pointer to a transformed []byte.
//
// This function is especially useful when one may need to extract
// multiple fields from the transformed JSON.
func (k *Kazaam) TransformJSONString(data string) ([]byte, error) {
	// read in the arbitrary input data
	d := make([]byte, len(data))
	copy(d, []byte(data))
	d, err := k.TransformInPlace(d)
	if err != nil {
		return []byte{}, err
	}
	return d, err
}
