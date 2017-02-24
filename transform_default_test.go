package kazaam_test

import (
	"testing"

	"github.com/qntfy/kazaam"
)

func TestDefault(t *testing.T) {
	jsonOut := `{"Range":5,"rating":{"example":{"value":3},"primary":{"value":3}}}`
	spec := `[{"operation": "default", "spec": {"Range": 5}}]`

	kazaamTransform, _ := kazaam.NewKazaam(spec)
	kazaamOut, _ := kazaamTransform.TransformJSONStringToString(testJSONInput)

	if kazaamOut != jsonOut {
		t.Error("Transformed data does not match expectation.")
		t.Log("Expected: ", jsonOut)
		t.Log("Actual:   ", kazaamOut)
		t.FailNow()
	}
}
