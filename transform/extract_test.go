package transform

import "testing"

func TestExtract(t *testing.T) {
	spec := `{"path": "_source"}`
	jsonIn := `{"data": {"id": true}, "_source": {"a": 123, "b": "str", "c": null, "d": true}}`
	jsonOut := `{"a":123,"b":"str","c":null,"d":true}`

	cfg := getConfig(spec, "", false)
	kazaamOut, _ := getTransformTestWrapper(Extract, cfg, jsonIn)

	if kazaamOut != jsonOut {
		t.Error("Transformed data does not match expectation.")
		t.Log("Expected: ", jsonOut)
		t.Log("Actual:   ", kazaamOut)
		t.FailNow()
	}
}

func TestExtractWithRequire(t *testing.T) {
	spec := `{"path": "not_source"}`
	jsonIn := `{"data": {"id": true}, "_source": {"a": 123, "b": "str", "c": null, "d": true}}`

	cfg := getConfig(spec, "", true)
	_, err := getTransformTestWrapper(Extract, cfg, jsonIn)

	if err == nil {
		t.Error("Transform path does not exist in message and should throw an error")
		t.FailNow()
	}
}

func TestExtractWithBadPath(t *testing.T) {
	spec := `{"path": "test"}`
	jsonIn := `{"data": {"id": true}, "_source": {"a": 123, "b": "str", "c": null, "d": true}}`
	jsonOut := "null"

	cfg := getConfig(spec, "", false)
	kazaamOut, _ := getTransformTestWrapper(Extract, cfg, jsonIn)

	if kazaamOut != jsonOut {
		t.Error("Transformed data does not match expectation.")
		t.Log("Expected: ", jsonOut)
		t.Log("Actual:   ", kazaamOut)
		t.FailNow()
	}
}
