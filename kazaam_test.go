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

func TestKazaamShiftTransform(t *testing.T) {
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

func TestKazaamShiftTransformAndGet(t *testing.T) {
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

func TestKazaamEncapsulateTransform(t *testing.T) {
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

func BenchmarkKazaamEncapsulateTransform(b *testing.B) {
	spec := `[{"operation": "shift", "spec": {"data": ["$"]}}]`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		kazaamTransform, _ := kazaam.NewKazaam(spec)
		kazaamTransform.TransformJSONStringToString(testJSONInput)
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

func TestKazaamDefaultTransform(t *testing.T) {
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

func TestKazaamOverShiftTransform(t *testing.T) {
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

func TestKazaamRadTransform(t *testing.T) {
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

func TestKazaamWildcardExtractTransform(t *testing.T) {
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

func TestKazaamEndArrayAccess(t *testing.T) {
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

func TestKazaamShiftNullArraySpecValue(t *testing.T) {
	spec := `[{"operation": "shift", "spec": {"id": [null, "abc"]}}]`
	jsonIn := `{"data": {"id": true}}`

	kazaamTransform, _ := kazaam.NewKazaam(spec)
	_, err := kazaamTransform.TransformJSONStringToString(jsonIn)

	errMsg := `Warn: Unable to coerce element to json string: <nil>`
	if err.Error() != errMsg {
		t.Error("Error data does not match expectation.")
		t.Log("Expected: ", errMsg)
		t.Log("Actual:   ", err.Error())
		t.FailNow()
	}
}

func TestKazaamShiftNullSpecValue(t *testing.T) {
	spec := `[{"operation": "shift", "spec": {"id": null}}]`
	jsonIn := `{"data": {"id": true}}`

	kazaamTransform, _ := kazaam.NewKazaam(spec)
	_, err := kazaamTransform.TransformJSONStringToString(jsonIn)

	errMsg := `Warn: Unknown type in message for key: id`
	if err.Error() != errMsg {
		t.Error("Error data does not match expectation.")
		t.Log("Expected: ", errMsg)
		t.Log("Actual:   ", err.Error())
		t.FailNow()
	}
}

func TestKazaamExtractTransform(t *testing.T) {
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

func TestKazaamExtractTransformBadPath(t *testing.T) {
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

func TestKazaamConcatTransformSimplePath(t *testing.T) {
	spec := `[{"operation": "concat", "spec": {"sources": [{"value": "TEST"}, {"path": "a.timestamp"}], "targetPath": "a.output", "delim": "," }}]`
	jsonIn := `{"a":{"timestamp": 1481305274}}`

	jsonOut := `{"a":{"output":"TEST,1481305274","timestamp":1481305274}}`

	kazaamTransform, _ := kazaam.NewKazaam(spec)
	kazaamOut, _ := kazaamTransform.TransformJSONStringToString(jsonIn)

	if kazaamOut != jsonOut {
		t.Error("Transformed data does not match expectation.")
		t.Log("Expected: ", jsonOut)
		t.Log("Actual:   ", kazaamOut)
		t.FailNow()
	}
}

func TestKazaamConcatTransformSimplePathReplace(t *testing.T) {
	spec := `[{"operation": "concat", "spec": {"sources": [{"value": "TEST"}, {"path": "a.timestamp"}], "targetPath": "a.timestamp", "delim": "," }}]`
	jsonIn := `{"a":{"timestamp": 1481305274}}`

	jsonOut := `{"a":{"timestamp":"TEST,1481305274"}}`

	kazaamTransform, _ := kazaam.NewKazaam(spec)
	kazaamOut, _ := kazaamTransform.TransformJSONStringToString(jsonIn)

	if kazaamOut != jsonOut {
		t.Error("Transformed data does not match expectation.")
		t.Log("Expected: ", jsonOut)
		t.Log("Actual:   ", kazaamOut)
		t.FailNow()
	}
}

func TestKazaamConcatTransformNoDelim(t *testing.T) {
	spec := `[{"operation": "concat", "spec": {"sources": [{"value": "TEST"}, {"path": "a.timestamp"}], "targetPath": "a.output" }}]`
	jsonIn := `{"a":{"timestamp": "1481305274"}}`

	jsonOut := `{"a":{"output":"TEST1481305274","timestamp":"1481305274"}}`

	kazaamTransform, _ := kazaam.NewKazaam(spec)
	kazaamOut, _ := kazaamTransform.TransformJSONStringToString(jsonIn)

	if kazaamOut != jsonOut {
		t.Error("Transformed data does not match expectation.")
		t.Log("Expected: ", jsonOut)
		t.Log("Actual:   ", kazaamOut)
		t.FailNow()
	}
}

func TestKazaamConcatTransformWildcard(t *testing.T) {

	spec := `[{"operation": "concat", "spec": {"sources": [{"value": "TEST"}, {"path": "a[*].foo"}], "targetPath": "a.output", "delim": "," }}]`
	jsonIn := `{"a":[{"foo": 0}, {"foo": 1}, {"foo": 1}, {"foo": 2}]}`

	jsonOut := `{"a":{"output":"TEST,0112"}}`

	kazaamTransform, _ := kazaam.NewKazaam(spec)
	kazaamOut, _ := kazaamTransform.TransformJSONStringToString(jsonIn)

	if kazaamOut != jsonOut {
		t.Error("Transformed data does not match expectation.")
		t.Log("Expected: ", jsonOut)
		t.Log("Actual:   ", kazaamOut)
		t.FailNow()
	}
}

func TestKazaamConcatTransformWildcardNested(t *testing.T) {

	spec := `[{"operation": "concat", "spec": {"sources": [{"value": "TEST"}, {"path": "a.b[*].foo"}], "targetPath": "a.output", "delim": "," }}]`
	jsonIn := `{"a": {"b": [{"foo": 0}, {"foo": 1}, {"foo": 1}, {"foo": 2}]}}`

	jsonOut := `{"a":{"b":[{"foo":0},{"foo":1},{"foo":1},{"foo":2}],"output":"TEST,0112"}}`

	kazaamTransform, _ := kazaam.NewKazaam(spec)
	kazaamOut, _ := kazaamTransform.TransformJSONStringToString(jsonIn)

	if kazaamOut != jsonOut {
		t.Error("Transformed data does not match expectation.")
		t.Log("Expected: ", jsonOut)
		t.Log("Actual:   ", kazaamOut)
		t.FailNow()
	}
}

func TestKazaamConcatTransformBadPath(t *testing.T) {

	spec := `[{"operation": "concat", "spec": {"sources": [{"value": "TEST"}, {"path": "a[*].bar"}], "targetPath": "a.output", "delim": "," }}]`
	jsonIn := `{"a":[{"foo": 0}, {"foo": 1}, {"foo": 1}, {"foo": 2}]}`

	jsonOut := `{"a":{"output":"TEST,"}}`

	kazaamTransform, _ := kazaam.NewKazaam(spec)
	kazaamOut, _ := kazaamTransform.TransformJSONStringToString(jsonIn)

	if kazaamOut != jsonOut {
		t.Error("Transformed data does not match expectation.")
		t.Log("Expected: ", jsonOut)
		t.Log("Actual:   ", kazaamOut)
		t.FailNow()
	}
}

func TestKazaamConcatTransformBadSpec(t *testing.T) {

	// Bad spec - "Path" should be "path"
	spec := `[{"operation": "concat", "spec": {"sources": [{"value": "TEST"}, {"Path": "a[*].bar"}], "targetPath": "a.timestamp", "delim": "," }}]`
	jsonIn := `{"a":[{"foo": 0}, {"foo": 1}, {"foo": 1}, {"foo": 2}]}`

	// bad path should cause the result to be blank
	jsonOut := ""

	kazaamTransform, _ := kazaam.NewKazaam(spec)
	kazaamOut, _ := kazaamTransform.TransformJSONStringToString(jsonIn)

	if kazaamOut != jsonOut {
		t.Error("Transformed data does not match expectation.")
		t.Log("Expected: ", jsonOut)
		t.Log("Actual:   ", kazaamOut)
		t.FailNow()
	}
}

func TestKazaamConcatTransformMulti(t *testing.T) {
	spec := `[{"operation": "concat", "spec": {"sources": [{"value": "BEGIN"}, {"path": "a[*].foo"}, {"value": "END"}], "targetPath": "a.output", "delim": "," }}]`
	jsonIn := `{"a":[{"foo": 0}, {"foo": 1}, {"foo": 1}, {"foo": 2}]}`

	jsonOut := `{"a":{"output":"BEGIN,0112,END"}}`

	kazaamTransform, _ := kazaam.NewKazaam(spec)
	kazaamOut, _ := kazaamTransform.TransformJSONStringToString(jsonIn)

	if kazaamOut != jsonOut {
		t.Error("Transformed data does not match expectation.")
		t.Log("Expected: ", jsonOut)
		t.Log("Actual:   ", kazaamOut)
		t.FailNow()
	}
}

func TestKazaamConcatTransformSingle(t *testing.T) {
	spec := `[{"operation": "concat", "spec": {"sources": [{"path": "a.timestamp"}], "targetPath": "a.output" }}]`
	jsonIn := `{"a":{"timestamp": 1481305274100000000000000000000}}`

	jsonOut := `{"a":{"output":"1481305274100000000000000000000","timestamp":1481305274100000000000000000000}}`

	kazaamTransform, _ := kazaam.NewKazaam(spec)
	kazaamOut, _ := kazaamTransform.TransformJSONStringToString(jsonIn)

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
