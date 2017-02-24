package kazaam

import (
	"encoding/json"

	"github.com/qntfy/kazaam/transform"
)

// Spec represents an individual spec element. It describes the name of the operation,
// whether the `over` operator is required, and an operation-specific `Config` that
// describes the configuration of the transform.
type spec struct {
	*transform.Config
	Operation *string `json:"operation"`
	Over      *string `json:"over,omitempty"`
}

// return the transform function based on what's indicated in the operation spec
func (s *spec) getTransform() TransformFunc {
	tform, ok := validSpecTypes[*s.Operation]
	if !ok {
		return transform.Pass
	}
	return tform
}

type specInt spec
type specs []spec

// UnmarshalJSON implements a custon unmarshaller for the Spec type
func (s *spec) UnmarshalJSON(b []byte) (err error) {
	j := specInt{}
	if err = json.Unmarshal(b, &j); err == nil {
		*s = spec(j)
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
