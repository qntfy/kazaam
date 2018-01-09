package transform

import "testing"

func TestCoalesce(t *testing.T) {
	spec := `{"foo": ["rating.foo", "rating.primary"]}`
	jsonOut := `{"rating":{"example":{"value":3},"primary":{"value":3}},"foo":{"value":3}}`

	cfg := getConfig(spec, false)
	kazaamOut, _ := getTransformTestWrapper(Coalesce, cfg, testJSONInput)
	areEqual, _ := checkJSONBytesEqual(kazaamOut, []byte(jsonOut))

	if !areEqual {
		t.Error("Transformed data does not match expectation.")
		t.Log("Expected:   ", jsonOut)
		t.Log("Actual:     ", kazaamOut)
		t.FailNow()
	}
}

func TestCoalesceWithRequire(t *testing.T) {
	spec := `{"foo": ["rating.foo", "rating.primary"]}`

	cfg := getConfig(spec, true)
	_, err := getTransformTestWrapper(Coalesce, cfg, testJSONInput)

	if err == nil {
		t.Error("Coalesce does not support \"require\" and should throw an error.")
		t.FailNow()
	}
}

func TestCoalesceWithMulti(t *testing.T) {
	spec := `{"foo": ["rating.foo", "rating.primary"], "bar": ["rating.bar", "rating.example.value"]}`
	jsonOut := `{"rating":{"example":{"value":3},"primary":{"value":3}},"foo":{"value":3},"bar":3}`

	cfg := getConfig(spec, false)
	kazaamOut, _ := getTransformTestWrapper(Coalesce, cfg, testJSONInput)
	areEqual, _ := checkJSONBytesEqual(kazaamOut, []byte(jsonOut))

	if !areEqual {
		t.Error("Transformed data does not match expectation.")
		t.Log("Expected:   ", jsonOut)
		t.Log("Actual:     ", kazaamOut)
		t.FailNow()
	}
}

func TestCoalesceWithNotFound(t *testing.T) {
	spec := `{"foo": ["rating.foo", "rating.bar", "ratings"]}`
	jsonOut := `{"rating":{"example":{"value":3},"primary":{"value":3}}}`

	cfg := getConfig(spec, false)
	kazaamOut, _ := getTransformTestWrapper(Coalesce, cfg, testJSONInput)
	areEqual, _ := checkJSONBytesEqual(kazaamOut, []byte(jsonOut))

	if !areEqual {
		t.Error("Transformed data does not match expectation.")
		t.Log("Expected:   ", jsonOut)
		t.Log("Actual:     ", kazaamOut)
		t.FailNow()
	}
}

func TestCoalesceWithIgnore(t *testing.T) {
	jsonInput := `{"rating":{"example":{"value":3},"primary":{"value":3},"foo":5,"bar":""}}`
	spec := `{"foo": ["rating.foo", "rating.primary"], "bar": ["rating.bar", "rating.example.value"], "ignore": [5, ""]}`
	jsonOut := `{"rating":{"example":{"value":3},"primary":{"value":3},"foo":5,"bar":""},"foo":{"value":3},"bar":3}`

	cfg := getConfig(spec, false)
	kazaamOut, _ := getTransformTestWrapper(Coalesce, cfg, jsonInput)
	areEqual, _ := checkJSONBytesEqual(kazaamOut, []byte(jsonOut))

	if !areEqual {
		t.Error("Transformed data does not match expectation.")
		t.Log("Expected:   ", jsonOut)
		t.Log("Actual:     ", string(kazaamOut))
		t.FailNow()
	}
}
