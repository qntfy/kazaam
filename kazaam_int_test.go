package kazaam_test

import (
	"testing"

	"github.com/qntfy/kazaam"
)

const testJSONInput = `{"rating": {"primary": {"value": 3}, "example": {"value": 3}}}`

func TestKazaamBadInput(t *testing.T) {
	jsonOut := ``
	spec := `[{"operation": "shift","spec": {"Rating": "rating.primary.value","example.old": "rating.example"}}]`

	kazaamTransform, _ := kazaam.NewKazaam(spec)
	kazaamOut, _ := kazaamTransform.TransformJSONStringToString("")

	if kazaamOut != jsonOut {
		t.Error("Transformed data does not match expectation.")
		t.Log("Expected: ", jsonOut)
		t.Log("Actual:   ", kazaamOut)
		t.FailNow()
	}
}

func TestKazaamBadJSONSpecification(t *testing.T) {
	_, err := kazaam.NewKazaam("{spec}")
	if err == nil {
		t.Error("Specification JSON is invalid and should throw an error")
		t.FailNow()
	}
}

func TestKazaamBadJSONTransform(t *testing.T) {
	kazaamTransform, _ := kazaam.NewKazaam(`[{"operation": "shift,"spec": {"data": ["$"]}}]`)
	_, err := kazaamTransform.TransformJSONString(`{"data"}`)
	if err == nil {
		t.Error("Specification JSON is invalid and should throw an error")
		t.FailNow()
	}
}

func TestKazaamBadJSONTransformNoOperation(t *testing.T) {
	_, err := kazaam.NewKazaam(`[{"opeeration": "shift","spec": {"data": ["$"]}}]`)
	if err == nil {
		t.Error("Specification JSON is invalid and should throw an error")
		t.FailNow()
	}
}

func TestKazaamBadJSONTransformBadOperation(t *testing.T) {
	_, err := kazaam.NewKazaam(`[{"operation":"invalid","spec": {"data": ["$"]}}]`)
	if err == nil {
		t.Error("Specification JSON is invalid and should throw an error")
		t.FailNow()
	}
}

func TestKazaamMultipleTransforms(t *testing.T) {
	jsonOut1 := `{"Rating":3,"example":{"old":{"value":3}}}`
	jsonOut2 := `{"Range":5,"rating":{"example":{"value":3},"primary":{"value":3}}}`
	spec1 := `[{"operation": "shift", "spec": {"Rating": "rating.primary.value", "example.old": "rating.example"}}]`
	spec2 := `[{"operation": "default", "spec": {"Range": 5}}]`

	transform1, _ := kazaam.NewKazaam(spec1)
	kazaamOut1, _ := transform1.TransformJSONStringToString(testJSONInput)

	transform2, _ := kazaam.NewKazaam(spec2)
	kazaamOut2, _ := transform2.TransformJSONStringToString(testJSONInput)

	if kazaamOut1 != jsonOut1 {
		t.Error("Transformed data does not match expectation.")
		t.Log("Expected: ", jsonOut1)
		t.Log("Actual:   ", kazaamOut1)
		t.FailNow()
	}

	if kazaamOut2 != jsonOut2 {
		t.Error("Transformed data does not match expectation.")
		t.Log("Expected: ", jsonOut2)
		t.Log("Actual:   ", kazaamOut2)
		t.FailNow()
	}
}

func TestKazaamMultipleTransformsRequire(t *testing.T) {
	jsonOut2 := `{"Range":5,"rating":{"example":{"value":3},"primary":{"value":3}}}`
	spec1 := `[{"operation": "shift", "spec": {"Rating": "rating.primary.no_value", "example.old": "rating.example"}, "require": true}]`
	spec2 := `[{"operation": "default", "spec": {"Range": 5}, "require": true}]`

	transform1, _ := kazaam.NewKazaam(spec1)
	_, out1Err := transform1.TransformJSONStringToString(testJSONInput)

	transform2, _ := kazaam.NewKazaam(spec2)
	kazaamOut2, _ := transform2.TransformJSONStringToString(testJSONInput)

	if out1Err == nil {
		t.Error("Transform path does not exist in message and should throw an error.")
		t.FailNow()
	}

	if kazaamOut2 != jsonOut2 {
		t.Error("Transformed data does not match expectation.")
		t.Log("Expected: ", jsonOut2)
		t.Log("Actual:   ", kazaamOut2)
		t.FailNow()
	}
}

func TestKazaamNoTransform(t *testing.T) {
	jsonOut := `{"rating":{"example":{"value":3},"primary":{"value":3}}}`
	var spec string

	kazaamTransform, _ := kazaam.NewKazaam(spec)
	kazaamOut, _ := kazaamTransform.TransformJSONStringToString(testJSONInput)

	if kazaamOut != jsonOut {
		t.Error("Transformed data does not match expectation.")
		t.Log("Expected: ", jsonOut)
		t.Log("Actual:   ", kazaamOut)
		t.FailNow()
	}
}

func TestKazaamCoalesceTransformAndShift(t *testing.T) {
	spec := `[{
		"operation": "coalesce",
		"spec": {"foo": ["rating.foo", "rating.primary"]}
	}, {
		"operation": "shift",
		"spec": {"rating.foo": "foo", "rating.example.value": "rating.primary.value"}
	}]`
	jsonOut := `{"rating":{"example":{"value":3},"foo":{"value":3}}}`

	kazaamTransform, _ := kazaam.NewKazaam(spec)
	kazaamOut, _ := kazaamTransform.TransformJSONStringToString(testJSONInput)

	if kazaamOut != jsonOut {
		t.Error("Transformed data does not match expectation.")
		t.Log("Expected: ", jsonOut)
		t.Log("Actual:   ", kazaamOut)
		t.FailNow()
	}
}

func TestKazaamTransformMultiOpWithOver(t *testing.T) {
	spec := `[{
		"operation": "concat",
		"over": "a",
		"spec": {"sources": [{"path": "foo"}, {"value": "KEY"}], "targetPath": "url", "delim": ":" }
	}, {
		"operation": "shift",
		"spec": {"urls": "a[*].url" }
	}]`
	jsonIn := `{"a":[{"foo": 0}, {"foo": 1}, {"foo": 2}]}`
	jsonOut := `{"urls":["0:KEY","1:KEY","2:KEY"]}`

	kazaamTransform, _ := kazaam.NewKazaam(spec)
	kazaamOut, _ := kazaamTransform.TransformJSONStringToString(jsonIn)

	if kazaamOut != jsonOut {
		t.Error("Transformed data does not match expectation.")
		t.Log("Expected: ", jsonOut)
		t.Log("Actual:   ", kazaamOut)
		t.FailNow()
	}
}

func BenchmarkKazaamTransformMultiOpWithOver(b *testing.B) {
	spec := `[{
		"operation": "concat",
		"over": "a",
		"spec": {"sources": [{"path": "foo"}, {"value": "KEY"}], "targetPath": "url", "delim": ":" }
	}, {
		"operation": "shift",
		"spec": {"urls": "a[*].url" }
	}]`
	jsonIn := `{"a":[{"foo": 0}, {"foo": 1}, {"foo": 2}]}`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		kazaamTransform, _ := kazaam.NewKazaam(spec)
		kazaamTransform.TransformJSONStringToString(jsonIn)
	}
}

func BenchmarkKazaamEncapsulateTransform(b *testing.B) {
	spec := `[{"operation": "shift", "spec": {"data": ["$"]}}]`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		kazaamTransform, _ := kazaam.NewKazaam(spec)
		kazaamTransform.TransformJSONStringToString(testJSONInput)
	}
}
