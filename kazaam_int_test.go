package kazaam_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/qntfy/jsonparser"
	"github.com/qntfy/kazaam"
	"github.com/qntfy/kazaam/transform"
)

const testJSONInput = `{"rating":{"example":{"value":3},"primary":{"value":3}}}`

func checkJSONStringsEqual(item1, item2 string) (bool, error) {
	var out1, out2 interface{}

	err := json.Unmarshal([]byte(item1), &out1)
	if err != nil {
		return false, nil
	}

	err = json.Unmarshal([]byte(item2), &out2)
	if err != nil {
		return false, nil
	}

	return reflect.DeepEqual(out1, out2), nil
}

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
	_, err := kazaamTransform.TransformJSONString(`{"data":"foo"}`)
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
	jsonOut2 := `{"rating":{"example":{"value":3},"primary":{"value":3}},"Range":5}`
	spec1 := `[{"operation": "shift", "spec": {"Rating": "rating.primary.value", "example.old": "rating.example"}}]`
	spec2 := `[{"operation": "default", "spec": {"Range": 5}}]`

	transform1, _ := kazaam.NewKazaam(spec1)
	kazaamOut1, _ := transform1.TransformJSONStringToString(testJSONInput)
	areEqual1, _ := checkJSONStringsEqual(kazaamOut1, jsonOut1)

	transform2, _ := kazaam.NewKazaam(spec2)
	kazaamOut2, _ := transform2.TransformJSONStringToString(testJSONInput)
	areEqual2, _ := checkJSONStringsEqual(kazaamOut2, jsonOut2)

	if !areEqual1 {
		t.Error("Transformed data does not match expectation.")
		t.Log("Expected: ", jsonOut1)
		t.Log("Actual:   ", kazaamOut1)
		t.FailNow()
	}

	if !areEqual2 {
		t.Error("Transformed data does not match expectation.")
		t.Log("Expected: ", jsonOut2)
		t.Log("Actual:   ", kazaamOut2)
		t.FailNow()
	}
}

func TestKazaamMultipleTransformsRequire(t *testing.T) {
	jsonOut2 := `{"rating":{"example":{"value":3},"primary":{"value":3}},"Range":5}`
	spec1 := `[{"operation": "shift", "spec": {"Rating": "rating.primary.no_value", "example.old": "rating.example"}, "require": true}]`
	spec2 := `[{"operation": "default", "spec": {"Range": 5}, "require": true}]`

	transform1, _ := kazaam.NewKazaam(spec1)
	_, out1Err := transform1.TransformJSONStringToString(testJSONInput)

	transform2, _ := kazaam.NewKazaam(spec2)
	kazaamOut2, _ := transform2.TransformJSONStringToString(testJSONInput)
	areEqual2, _ := checkJSONStringsEqual(kazaamOut2, jsonOut2)

	if out1Err == nil {
		t.Error("Transform path does not exist in message and should throw an error.")
		t.FailNow()
	}

	if !areEqual2 {
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
	areEqual, _ := checkJSONStringsEqual(kazaamOut, jsonOut)

	if !areEqual {
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
	jsonOut := `{"rating":{"foo":{"value":3},"example":{"value":3}}}`

	kazaamTransform, _ := kazaam.NewKazaam(spec)
	kazaamOut, _ := kazaamTransform.TransformJSONStringToString(testJSONInput)
	areEqual, _ := checkJSONStringsEqual(kazaamOut, jsonOut)

	if !areEqual {
		t.Error("Transformed data does not match expectation.")
		t.Log("Expected: ", jsonOut)
		t.Log("Actual:   ", kazaamOut)
		t.FailNow()
	}
}

func TestKazaamCoalesceTransformAndShiftWithKeySeparator(t *testing.T) {
	spec := `[{
		"operation": "coalesce",
		"spec": {"foo": ["rating.foo", "rating.primary"]},
		"keySeparator": "."
	}, {
		"operation": "shift",
		"spec": {"rating.foo": "foo", "rating.example.value": "rating.primary.value"},
		"keySeparator": "."
	}]`
	jsonOut := `{"rating":{"foo":{"value":3},"example":{"value":3}}}`

	kazaamTransform, _ := kazaam.NewKazaam(spec)
	kazaamOut, _ := kazaamTransform.TransformJSONStringToString(testJSONInput)
	areEqual, _ := checkJSONStringsEqual(kazaamOut, jsonOut)

	if !areEqual {
		t.Error("Transformed data does not match expectation.")
		t.Log("Expected: ", jsonOut)
		t.Log("Actual:   ", kazaamOut)
		t.FailNow()
	}
}

func TestKazaamCoalesceTransformAndShiftWithKeySeparatorNonDefault(t *testing.T) {
	spec := `[{
		"operation": "coalesce",
		"spec": {"foo": ["rating>foo", "rating>primary"]},
		"keySeparator": ">"
	}, {
		"operation": "shift",
		"spec": {"rating.foo": "foo", "rating.example.value": "rating.primary.value"}
	}]`
	jsonOut := `{"rating":{"foo":{"value":3},"example":{"value":3}}}`

	kazaamTransform, _ := kazaam.NewKazaam(spec)
	kazaamOut, _ := kazaamTransform.TransformJSONStringToString(testJSONInput)
	areEqual, _ := checkJSONStringsEqual(kazaamOut, jsonOut)

	if !areEqual {
		t.Error("Transformed data does not match expectation.")
		t.Log("Expected: ", jsonOut)
		t.Log("Actual:   ", kazaamOut)
		t.FailNow()
	}
}

func TestKazaamShiftTransformWithTimestamp(t *testing.T) {
	spec := `[{
		"operation": "shift",
		"spec": {"newTimestamp":"oldTimestamp","oldTimestamp":"oldTimestamp"}
	}, {
		"operation": "timestamp",
		"spec": {"newTimestamp":{"inputFormat":"Mon Jan _2 15:04:05 -0700 2006","outputFormat":"2006-01-02T15:04:05-0700"}}
	}]`
	jsonIn := `{"oldTimestamp":"Fri Jul 21 08:15:27 +0000 2017"}`
	jsonOut := `{"oldTimestamp":"Fri Jul 21 08:15:27 +0000 2017","newTimestamp":"2017-07-21T08:15:27+0000"}`

	kazaamTransform, _ := kazaam.NewKazaam(spec)
	kazaamOut, _ := kazaamTransform.TransformJSONStringToString(jsonIn)
	areEqual, _ := checkJSONStringsEqual(kazaamOut, jsonOut)

	if !areEqual {
		t.Error("Transformed data does not match expectation.")
		t.Log("Expected: ", jsonOut)
		t.Log("Actual:   ", kazaamOut)
		t.FailNow()
	}
}

func TestShiftWithOverAndWildcard(t *testing.T) {
	spec := `[{"operation": "shift","spec": {"docs": "documents[*]"}}, {"operation": "shift",  "spec": {"data": "norm.text"}, "over":"docs"}]`
	jsonIn := `{"documents":[{"norm":{"text":"String 1"}},{"norm":{"text":"String 2"}}]}`
	jsonOut := `{"docs":[{"data":"String 1"},{"data":"String 2"}]}`

	kazaamTransform, _ := kazaam.NewKazaam(spec)
	kazaamOut, err := kazaamTransform.TransformJSONStringToString(jsonIn)

	if err != nil {
		t.Error("Transform produced error.")
		t.Log("Error: ", err.Error())
		t.Log("Expected: ", jsonOut)
		t.Log("Actual:   ", kazaamOut)
		t.FailNow()
	}

	areEqual, _ := checkJSONStringsEqual(kazaamOut, jsonOut)
	if !areEqual {
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
	jsonIn := `{"a":[{"foo":0},{"foo":1},{"foo":2}]}`
	jsonOut := `{"urls":["0:KEY","1:KEY","2:KEY"]}`

	kazaamTransform, _ := kazaam.NewKazaam(spec)
	kazaamOut, _ := kazaamTransform.TransformJSONStringToString(jsonIn)
	areEqual, _ := checkJSONStringsEqual(kazaamOut, jsonOut)

	if !areEqual {
		t.Error("Transformed data does not match expectation.")
		t.Log("Expected: ", jsonOut)
		t.Log("Actual:   ", kazaamOut)
		t.FailNow()
	}
}

func TestShiftWithOver(t *testing.T) {
	jsonIn := `{"rating":{"primary":[{"value":3},{"value":5}],"example":{"value":3}}}`
	jsonOut := `{"rating":{"primary":[{"new_value":3},{"new_value":5}],"example":{"value":3}}}`
	spec := `[{"operation": "shift", "over": "rating.primary", "spec": {"new_value":"value"}}]`

	kazaamTransform, _ := kazaam.NewKazaam(spec)
	kazaamOut, _ := kazaamTransform.TransformJSONStringToString(jsonIn)
	areEqual, _ := checkJSONStringsEqual(kazaamOut, jsonOut)

	if !areEqual {
		t.Error("Transformed data does not match expectation.")
		t.Log("Expected: ", jsonOut)
		t.Log("Actual:   ", kazaamOut)
		t.FailNow()
	}
}

func TestShiftAndGet(t *testing.T) {
	jsonOut := "3"
	spec := `[{"operation": "shift","spec": {"Rating": "rating.primary.value","example.old": "rating.example"}}]`

	kazaamTransform, _ := kazaam.NewKazaam(spec)
	transformed, err := kazaamTransform.TransformJSONString(testJSONInput)
	if err != nil {
		t.Error("Failed to parse JSON message before transformation")
		t.FailNow()
	}
	kazaamOut, _, _, err := jsonparser.Get(transformed, "Rating")
	if err != nil {
		t.Log("Requested key not found")
	}

	if string(kazaamOut) != jsonOut {
		t.Error("Transformed data does not match expectation.")
		t.Log("Expected: ", jsonOut)
		t.Log("Actual:   ", string(kazaamOut))
		t.FailNow()
	}
}

func TestMissingRequiredField(t *testing.T) {
	jsonIn := `{"meta": {"not_image_cache": null}, "doc": "example"}`
	spec := `[
 		{"operation": "shift", "spec": {"results": "meta.image_cache[0].results[*]"}, "require": true}
	]`

	kazaamTransform, _ := kazaam.NewKazaam(spec)
	k, err := kazaamTransform.TransformJSONStringToString(jsonIn)

	if err == nil {
		t.Error("Should have generated error for null image_cache value")
		t.Error(k)
	}
	e := err.(*kazaam.Error)
	if e.ErrType != kazaam.RequireError {
		t.Error("Unexpected error type")
	}
}

func TestKazaamNoModify(t *testing.T) {
	spec := `[{"operation": "shift","spec": {"Rating": "rating.primary.value","example.old": "rating.example"}}]`
	msgOut := `{"Rating":3,"example":{"old":{"value":3}}}`
	tf, _ := kazaam.NewKazaam(spec)
	data := []byte(testJSONInput)
	jsonOut, _ := tf.Transform(data)

	areEqual1, _ := checkJSONStringsEqual(string(jsonOut), msgOut)

	if !areEqual1 {
		t.Error("Unexpected transformation result")
		t.Error("Actual:", string(jsonOut))
		t.Error("Expected:", msgOut)
	}

	areEqual2, _ := checkJSONStringsEqual(string(data), testJSONInput)

	if !areEqual2 {
		t.Error("Unexpected modification")
		t.Error("Actual:", string(data))
		t.Error("Expected:", testJSONInput)
	}
}

func TestConfigdKazaamGet3rdPartyTransform(t *testing.T) {
	kc := kazaam.NewDefaultConfig()
	kc.RegisterTransform("3rd-party", func(spec *transform.Config, data []byte) ([]byte, error) {
		data, _ = jsonparser.Set(data, []byte(`"does-exist"`), "doesnt-exist")
		return data, nil
	})
	msgOut := `{"test":"data","doesnt-exist":"does-exist"}`

	k, _ := kazaam.New(`[{"operation": "3rd-party"}]`, kc)
	kazaamOut, _ := k.TransformJSONStringToString(`{"test":"data"}`)
	areEqual, _ := checkJSONStringsEqual(kazaamOut, msgOut)
	if !areEqual {
		t.Error("Unexpected transform output")
		t.Log("Actual:   ", kazaamOut)
		t.Log("Expected: ", msgOut)

	}
}

func TestKazaamTransformThreeOpWithOver(t *testing.T) {
	spec := `[{
		"operation": "shift",
		"spec":{"a": "key.array1[0].array2[*]"}
	},
	{
		"operation": "concat",
		"over": "a",
		"spec": {"sources": [{"path": "foo"}, {"value": "KEY"}], "targetPath": "url", "delim": ":" }
	}, {
		"operation": "shift",
		"spec": {"urls": "a[*].url" }
	}]`
	jsonIn := `{"key":{"array1":[{"array2":[{"foo":0},{"foo":1},{"foo":2}]}]}}`
	jsonOut := `{"urls":["0:KEY","1:KEY","2:KEY"]}`

	kazaamTransform, _ := kazaam.NewKazaam(spec)
	kazaamOut, _ := kazaamTransform.TransformJSONStringToString(jsonIn)
	areEqual, _ := checkJSONStringsEqual(kazaamOut, jsonOut)

	if !areEqual {
		t.Error("Transformed data does not match expectation.")
		t.Log("Expected: ", jsonOut)
		t.Log("Actual:   ", kazaamOut)
		t.FailNow()
	}

}

func TestKazaamTransformThreeOpWithOverRequire(t *testing.T) {
	spec := `[{
		"operation": "shift",
		"spec":{"a": "key.array1[0].array2[*]"},
		"require": true
	},
	{
		"operation": "concat",
		"over": "a",
		"spec": {"sources": [{"path": "foo"}, {"value": "KEY"}], "targetPath": "url", "delim": ":" }
	}, {
		"operation": "shift",
		"spec": {"urls": "a[*].url" }
	}]`
	jsonIn := `{"key":{"not_array1":[{"array2":[{"foo": 0}, {"foo": 1}, {"foo": 2}]}]}}`

	kazaamTransform, _ := kazaam.NewKazaam(spec)
	_, err := kazaamTransform.TransformJSONStringToString(jsonIn)
	if err == nil {
		t.Error("Transform path does not exist in message and should throw an error")
		t.FailNow()
	}
}

func TestKazaamTransformTwoOpWithOverRequire(t *testing.T) {
	spec := `[{
		"operation": "shift",
		"spec":{"a": "key.array1[0].array2[*]"},
		"require": true
	},
	{
		"operation": "concat",
		"over": "a",
		"spec": {"sources": [{"path": "foo"}, {"value": "KEY"}], "targetPath": "url", "delim": ":" }
	}]`
	jsonIn := `{"key":{"not_array1":[{"array2":[{"foo": 0}, {"foo": 1}, {"foo": 2}]}]}}`

	kazaamTransform, _ := kazaam.NewKazaam(spec)
	_, err := kazaamTransform.TransformJSONStringToString(jsonIn)
	if err == nil {
		t.Error("Transform path does not exist in message and should throw an error")
		t.FailNow()
	}
}

func TestKazaamTransformDelete(t *testing.T) {
	spec := `[{
		"operation": "delete",
		"spec": {"paths": ["doc.uid", "doc.guidObjects[1]"]}
	}]`
	jsonIn := `{"doc":{"uid":12345,"guid":["guid0","guid2","guid4"],"guidObjects":[{"id":"guid0"},{"id":"guid2"},{"id":"guid4"}]}}`
	jsonOut := `{"doc":{"guid":["guid0","guid2","guid4"],"guidObjects":[{"id":"guid0"},{"id":"guid4"}]}}`

	kazaamTransform, _ := kazaam.NewKazaam(spec)
	kazaamOut, _ := kazaamTransform.TransformJSONStringToString(jsonIn)
	areEqual, _ := checkJSONStringsEqual(kazaamOut, jsonOut)

	if !areEqual {
		t.Error("Transformed data does not match expectation.")
		t.Log("Expected: ", jsonOut)
		t.Log("Actual:   ", kazaamOut)
		t.FailNow()
	}
}

func TestKazaamOverArrayStrings(t *testing.T) {
	spec := `[{
                "operation": "shift",
                "over": "doc.guidObjects",
                "spec": {"raw": "$"}
        }]`
        jsonIn := `{"doc":{"guidObjects":["foo",5,false]}}`
        jsonOut := `{"doc":{"guidObjects":[{"raw":"foo"},{"raw":5},{"raw":false}]}}`

        kazaamTransform, _ := kazaam.NewKazaam(spec)
        kazaamOut, _ := kazaamTransform.TransformJSONStringToString(jsonIn)
        areEqual, _ := checkJSONStringsEqual(kazaamOut, jsonOut)

        if !areEqual {
                t.Error("Transformed data does not match expectation.")
                t.Log("Expected: ", jsonOut)
                t.Log("Actual:   ", kazaamOut)
                t.FailNow()
        }
}
