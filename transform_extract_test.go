package kazaam_test

import (
	"testing"

	"github.com/qntfy/kazaam"
)

func TestExtract(t *testing.T) {
	spec := `[{"operation": "extract", "spec": {"path": "_source"}}]`
	jsonIn := `{"data": {"id": true}, "_source": {"a": 123, "b": "str", "c": null, "d": true}}`

	jsonOut := `{"a":123,"b":"str","c":null,"d":true}`

	kazaamTransform, _ := kazaam.NewKazaam(spec)
	kazaamOut, _ := kazaamTransform.TransformJSONStringToString(jsonIn)

	if kazaamOut != jsonOut {
		t.Error("Transformed data does not match expectation.")
		t.Log("Expected: ", jsonOut)
		t.Log("Actual:   ", kazaamOut)
		t.FailNow()
	}
}

func TestExtractWithRequire(t *testing.T) {
	spec := `[{"operation": "extract", "spec": {"path": "not_source"},"require": true}]`
	jsonIn := `{"data": {"id": true}, "_source": {"a": 123, "b": "str", "c": null, "d": true}}`

	kazaamTransform, _ := kazaam.NewKazaam(spec)
	_, err := kazaamTransform.TransformJSONStringToString(jsonIn)

	if err == nil {
		t.Error("Transform path does not exist in message and should throw an error")
		t.FailNow()
	}
}

func TestExtractWithBadPath(t *testing.T) {
	spec := `[{"operation": "extract", "spec": {"path": "test"}}]`
	jsonIn := `{"data": {"id": true}, "_source": {"a": 123, "b": "str", "c": null, "d": true}}`

	jsonOut := "null"

	kazaamTransform, _ := kazaam.NewKazaam(spec)
	kazaamOut, _ := kazaamTransform.TransformJSONStringToString(jsonIn)

	if kazaamOut != jsonOut {
		t.Error("Transformed data does not match expectation.")
		t.Log("Expected: ", jsonOut)
		t.Log("Actual:   ", kazaamOut)
		t.FailNow()
	}
}

func TestExtractWithWildcard(t *testing.T) {
	spec := `[{"operation": "shift", "spec": {"outputArray": "docs[*].data.key"}}]`
	jsonIn := `{"docs": [{"data": {"key": "val1"}},{"data": {"key": "val2"}}]}`
	jsonOut := `{"outputArray":["val1","val2"]}`

	kazaamTransform, _ := kazaam.NewKazaam(spec)
	kazaamOut, _ := kazaamTransform.TransformJSONStringToString(jsonIn)

	if kazaamOut != jsonOut {
		t.Error("Transformed data does not match expectation.")
		t.Log("Expected: ", jsonOut)
		t.Log("Actual:   ", kazaamOut)
		t.FailNow()
	}
}
