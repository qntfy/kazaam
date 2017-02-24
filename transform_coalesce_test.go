package kazaam_test

import (
	"testing"

	"github.com/qntfy/kazaam"
)

func TestCoalesce(t *testing.T) {
	spec := `[{"operation": "coalesce", "spec": {"foo": ["rating.foo", "rating.primary"]}}]`
	jsonOut := `{"foo":{"value":3},"rating":{"example":{"value":3},"primary":{"value":3}}}`

	kazaamTransform, _ := kazaam.NewKazaam(spec)
	kazaamOut, _ := kazaamTransform.TransformJSONStringToString(testJSONInput)

	if kazaamOut != jsonOut {
		t.Error("Transformed data does not match expectation.")
		t.Log("Expected: ", jsonOut)
		t.Log("Actual:   ", kazaamOut)
		t.FailNow()
	}
}

func TestCoalesceWithRequire(t *testing.T) {
	spec := `[{"operation": "coalesce", "spec": {"foo": ["rating.foo", "rating.primary"]},"require": true}]`

	kazaamTransform, _ := kazaam.NewKazaam(spec)
	_, err := kazaamTransform.TransformJSONStringToString(testJSONInput)

	if err == nil {
		t.Error("Coalesce does not support \"require\" and should throw an error.")
		t.FailNow()
	}
}

func TestCoalesceWithMulti(t *testing.T) {
	spec := `[{"operation": "coalesce", "spec": {"foo": ["rating.foo", "rating.primary"], "bar": ["rating.bar", "rating.example.value"]}}]`
	jsonOut := `{"bar":3,"foo":{"value":3},"rating":{"example":{"value":3},"primary":{"value":3}}}`

	kazaamTransform, _ := kazaam.NewKazaam(spec)
	kazaamOut, _ := kazaamTransform.TransformJSONStringToString(testJSONInput)

	if kazaamOut != jsonOut {
		t.Error("Transformed data does not match expectation.")
		t.Log("Expected: ", jsonOut)
		t.Log("Actual:   ", kazaamOut)
		t.FailNow()
	}
}

func TestCoalesceWithNotFound(t *testing.T) {
	spec := `[{"operation": "coalesce", "spec": {"foo": ["rating.foo", "rating.bar", "ratings"]}}]`
	jsonOut := `{"rating":{"example":{"value":3},"primary":{"value":3}}}`

	kazaamTransform, _ := kazaam.NewKazaam(spec)
	kazaamOut, _ := kazaamTransform.TransformJSONStringToString(testJSONInput)

	if kazaamOut != jsonOut {
		t.Error("Transformed data does not match expectation.")
		t.Log("Expected: ", jsonOut)
		t.Log("Actual:   ", kazaamOut)
		t.FailNow()
	}
}
