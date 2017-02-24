package kazaam_test

import (
	"testing"

	"github.com/qntfy/kazaam"
)

func TestShiftWithOver(t *testing.T) {
	jsonIn := `{"rating": {"primary": [{"value": 3}, {"value": 5}], "example": {"value": 3}}}`
	jsonOut := `{"rating":{"example":{"value":3},"primary":[{"new_value":3},{"new_value":5}]}}`
	spec := `[{"operation": "shift", "spec": {"new_value":"value"}, "over":"rating.primary"}]`

	kazaamTransform, _ := kazaam.NewKazaam(spec)
	kazaamOut, _ := kazaamTransform.TransformJSONStringToString(jsonIn)

	if kazaamOut != jsonOut {
		t.Error("Transformed data does not match expectation.")
		t.Log("Expected: ", jsonOut)
		t.Log("Actual:   ", kazaamOut)
		t.FailNow()
	}
}

func TestShiftWithOverAndWildcard(t *testing.T) {
	jsonIn := `{"documents":[{"norm": {"text": "String 1"}}, {"norm": {"text": "String 2"}}]}`
	jsonOut := `{"docs":[{"data":"String 1"},{"data":"String 2"}]}`
	spec := `[{"operation": "shift","spec": {"docs": "documents[*]"}}, {"operation": "shift",  "spec": {"data": "norm.text"}, "over":"docs"}]`

	kazaamTransform, _ := kazaam.NewKazaam(spec)
	kazaamOut, _ := kazaamTransform.TransformJSONStringToString(jsonIn)

	if kazaamOut != jsonOut {
		t.Error("Transformed data does not match expectation.")
		t.Log("Expected: ", jsonOut)
		t.Log("Actual:   ", kazaamOut)
		t.FailNow()
	}
}

func TestShift(t *testing.T) {
	jsonOut := `{"Rating":3,"example":{"old":{"value":3}}}`
	spec := `[{"operation": "shift","spec": {"Rating": "rating.primary.value","example.old": "rating.example"}}]`

	kazaamTransform, _ := kazaam.NewKazaam(spec)
	kazaamOut, _ := kazaamTransform.TransformJSONStringToString(testJSONInput)

	if kazaamOut != jsonOut {
		t.Error("Transformed data does not match expectation.")
		t.Log("Expected: ", jsonOut)
		t.Log("Actual:   ", kazaamOut)
		t.FailNow()
	}
}

func TestShiftAndGet(t *testing.T) {
	jsonOut := 3
	spec := `[{"operation": "shift","spec": {"Rating": "rating.primary.value","example.old": "rating.example"}}]`

	kazaamTransform, _ := kazaam.NewKazaam(spec)
	transformed, err := kazaamTransform.TransformJSONString(testJSONInput)
	if err != nil {
		t.Error("Failed to parse JSON message before transformation")
		t.FailNow()
	}
	kazaamOut, found := transformed.CheckGet("Rating")
	if !found {
		t.Log("Requested key not found")
	}

	if kazaamOut.MustInt() != jsonOut {
		t.Error("Transformed data does not match expectation.")
		t.Log("Expected: ", jsonOut)
		t.Log("Actual:   ", kazaamOut.MustInt())
		t.FailNow()
	}
}

func TestShiftWithMissingKey(t *testing.T) {
	jsonOut := `{"Rating":null,"example":{"old":{"value":3}}}`
	spec := `[{"operation": "shift","spec": {"Rating": "rating.primary.missing_value","example.old": "rating.example"}}]`

	kazaamTransform, _ := kazaam.NewKazaam(spec)
	kazaamOut, _ := kazaamTransform.TransformJSONStringToString(testJSONInput)

	if kazaamOut != jsonOut {
		t.Error("Transformed data does not match expectation.")
		t.Log("Expected: ", jsonOut)
		t.Log("Actual:   ", kazaamOut)
		t.FailNow()
	}
}

func TestShiftDeepExistsRequire(t *testing.T) {
	testJSONInput := `{"rating":{"example":[{"array":[{"value":3}]},{"another":"object"}]}}`
	spec := `[{"operation": "shift", "spec": {"example_res":"rating.example[0].array[*].value"},"require": true}]`
	jsonOut := `{"example_res":[3]}`

	transform, _ := kazaam.NewKazaam(spec)
	kazaamOut, _ := transform.TransformJSONStringToString(testJSONInput)

	if kazaamOut != jsonOut {
		t.Error("Transformed data does not match expectation.")
		t.Log("Expected: ", jsonOut)
		t.Log("Actual:   ", kazaamOut)
		t.FailNow()
	}
}

func TestShiftShallowExistsRequire(t *testing.T) {
	spec := `[{"operation": "shift","spec": {"Rating": "not_a_field"},"require": true}]`

	kazaamTransform, _ := kazaam.NewKazaam(spec)
	_, err := kazaamTransform.TransformJSONStringToString(testJSONInput)
	if err == nil {
		t.Error("Transform path does not exist in message and should throw an error")
		t.FailNow()
	}
}

func TestShiftDeepArraysRequire(t *testing.T) {
	spec := `[{"operation": "shift","spec": {"Rating": "rating.does[0].not[*].exist"}, "require": true}]`

	kazaamTransform, _ := kazaam.NewKazaam(spec)
	_, err := kazaamTransform.TransformJSONStringToString(testJSONInput)
	if err == nil {
		t.Error("Transform path does not exist in message and should throw an error")
		t.FailNow()
	}
}

func TestShiftDeepNoArraysRequire(t *testing.T) {
	spec := `[{"operation": "shift","spec": {"Rating": "rating.does.not.exist"}, "require": true}]`

	kazaamTransform, _ := kazaam.NewKazaam(spec)
	_, err := kazaamTransform.TransformJSONStringToString(testJSONInput)
	if err == nil {
		t.Error("Transform path does not exist in message and should throw an error")
		t.FailNow()
	}
}

func TestShiftWithEncapsulate(t *testing.T) {
	jsonOut := `{"data":[{"rating":{"example":{"value":3},"primary":{"value":3}}}]}`
	spec := `[{"operation": "shift", "spec": {"data": ["$"]}}]`

	kazaamTransform, _ := kazaam.NewKazaam(spec)
	kazaamOut, _ := kazaamTransform.TransformJSONStringToString(testJSONInput)

	if kazaamOut != jsonOut {
		t.Error("Transformed data does not match expectation.")
		t.Log("Expected: ", jsonOut)
		t.Log("Actual:   ", kazaamOut)
		t.FailNow()
	}
}

func TestShiftWithNullSpecValue(t *testing.T) {
	spec := `[{"operation": "shift", "spec": {"id": null}}]`
	jsonIn := `{"data": {"id": true}}`

	kazaamTransform, _ := kazaam.NewKazaam(spec)
	_, err := kazaamTransform.TransformJSONStringToString(jsonIn)

	errMsg := `ParseError - Warn: Unknown type in message for key: id`
	if err.Error() != errMsg {
		t.Error("Error data does not match expectation.")
		t.Log("Expected: ", errMsg)
		t.Log("Actual:   ", err.Error())
		t.FailNow()
	}
}

func TestShiftWithNullArraySpecValue(t *testing.T) {
	spec := `[{"operation": "shift", "spec": {"id": [null, "abc"]}}]`
	jsonIn := `{"data": {"id": true}}`

	kazaamTransform, _ := kazaam.NewKazaam(spec)
	_, err := kazaamTransform.TransformJSONStringToString(jsonIn)

	errMsg := `ParseError - Warn: Unable to coerce element to json string: <nil>`

	if err.Error() != errMsg {
		t.Error("Error data does not match expectation.")
		t.Log("Expected: ", errMsg)
		t.Log("Actual:   ", err.Error())
		t.FailNow()
	}
}

func TestShiftWithEndArrayAccess(t *testing.T) {
	spec := `[{"operation": "shift", "spec": {"id": "docs[1].data[0]"}}]`
	jsonIn := `{"docs": [{"data": ["abc", "def"]},{"data": ["ghi", "jkl"]}]}`
	jsonOut := `{"id":"ghi"}`

	kazaamTransform, _ := kazaam.NewKazaam(spec)
	kazaamOut, _ := kazaamTransform.TransformJSONStringToString(jsonIn)

	if kazaamOut != jsonOut {
		t.Error("Transformed data does not match expectation.")
		t.Log("Expected: ", jsonOut)
		t.Log("Actual:   ", kazaamOut)
		t.FailNow()
	}
}
