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
		if s.Config != nil && s.Spec != nil && len(*s.Spec) < 1 {
			err = &Error{ErrMsg: "Spec must contain at least one element", ErrType: SpecError}
			return
		}
		return
	}
	return
}
