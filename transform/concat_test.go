package transform

import "testing"

func TestConcat(t *testing.T) {
	spec := `{"sources": [{"value": "TEST"}, {"path": "a.timestamp"}], "targetPath": "a.output", "delim": "," }`
	jsonIn := `{"a":{"timestamp":1481305274}}`
	jsonOut := `{"a":{"timestamp":1481305274,"output":"TEST,1481305274"}}`

	cfg := getConfig(spec, false)
	kazaamOut, _ := getTransformTestWrapper(Concat, cfg, jsonIn)

	if kazaamOut != jsonOut {
		t.Error("Transformed data does not match expectation.")
		t.Log("Expected:   ", jsonOut)
		t.Log("Actual:     ", kazaamOut)
		t.FailNow()
	}
}

func TestConcatWithRequireSources(t *testing.T) {
	spec := `{"targetPath": "a.output", "delim": "," }`
	jsonIn := `{"a":{"timestamp":1481305274}}`

	cfg := getConfig(spec, true)
	_, err := getTransformTestWrapper(Concat, cfg, jsonIn)

	if err == nil {
		t.Error("Source field is missing and should throw an error.")
		t.FailNow()
	}
}

func TestConcatWithRequireTargetPath(t *testing.T) {
	spec := `{"sources": [{"value": "TEST"}, {"path": "a.timestamp"}], "delim": "," }`
	jsonIn := `{"a":{"timestamp":1481305274}}`

	cfg := getConfig(spec, true)
	_, err := getTransformTestWrapper(Concat, cfg, jsonIn)

	if err == nil {
		t.Error("targetPath field is missing and should throw an error.")
		t.FailNow()
	}
}

func TestConcatWithRequireSimplePath(t *testing.T) {
	spec := `{"sources": [{"value": "TEST"}, {"path": "not.a.timestamp"}], "targetPath": "a.output", "delim": "," }`
	jsonIn := `{"a":{"timestamp":1481305274}}`

	cfg := getConfig(spec, true)
	_, err := getTransformTestWrapper(Concat, cfg, jsonIn)

	if err == nil {
		t.Error("Transform path does not exist in message and should throw an error")
		t.FailNow()
	}
}

func TestConcatWithReplaceSimplePath(t *testing.T) {
	spec := `{"sources": [{"value": "TEST"}, {"path": "a.timestamp"}], "targetPath": "a.timestamp", "delim": "," }`
	jsonIn := `{"a":{"timestamp":1481305274}}`
	jsonOut := `{"a":{"timestamp":"TEST,1481305274"}}`

	cfg := getConfig(spec, false)
	kazaamOut, _ := getTransformTestWrapper(Concat, cfg, jsonIn)

	if kazaamOut != jsonOut {
		t.Error("Transformed data does not match expectation.")
		t.Log("Expected:   ", jsonOut)
		t.Log("Actual:     ", kazaamOut)
		t.FailNow()
	}
}

func TestConcatWithNoDelimiter(t *testing.T) {
	spec := `{"sources": [{"value": "TEST"}, {"path": "a.timestamp"}], "targetPath": "a.output" }`
	jsonIn := `{"a":{"timestamp":"1481305274"}}`
	jsonOut := `{"a":{"timestamp":"1481305274","output":"TEST1481305274"}}`

	cfg := getConfig(spec, false)
	kazaamOut, _ := getTransformTestWrapper(Concat, cfg, jsonIn)

	if kazaamOut != jsonOut {
		t.Error("Transformed data does not match expectation.")
		t.Log("Expected:   ", jsonOut)
		t.Log("Actual:     ", kazaamOut)
		t.FailNow()
	}
}

func TestConcatWithWildcard(t *testing.T) {
	spec := `{"sources": [{"value": "TEST"}, {"path": "a[*].foo"}], "targetPath": "a.output", "delim": "," }`
	jsonIn := `{"a":[{"foo":0},{"foo":1},{"foo":1},{"foo":2}]}`
	jsonOut := `{"a":{"output":"TEST,0112"}}`

	cfg := getConfig(spec, false)
	kazaamOut, _ := getTransformTestWrapper(Concat, cfg, jsonIn)

	if kazaamOut != jsonOut {
		t.Error("Transformed data does not match expectation.")
		t.Log("Expected:   ", jsonOut)
		t.Log("Actual:     ", kazaamOut)
		t.FailNow()
	}
}

func TestConcatWithWildcardNested(t *testing.T) {
	spec := `{"sources": [{"value": "TEST"}, {"path": "a.b[*].foo"}], "targetPath": "a.output", "delim": "," }`
	jsonIn := `{"a":{"b":[{"foo":0},{"foo":1},{"foo":1},{"foo":2}]}}`
	jsonOut := `{"a":{"b":[{"foo":0},{"foo":1},{"foo":1},{"foo":2}],"output":"TEST,0112"}}`

	cfg := getConfig(spec, false)
	kazaamOut, _ := getTransformTestWrapper(Concat, cfg, jsonIn)

	if kazaamOut != jsonOut {
		t.Error("Transformed data does not match expectation.")
		t.Log("Expected:   ", jsonOut)
		t.Log("Actual:     ", kazaamOut)
		t.FailNow()
	}
}

func TestConcatWithBadPath(t *testing.T) {
	spec := `{"sources": [{"value": "TEST"}, {"path": "a[*].bar"}], "targetPath": "a.output", "delim": "," }`
	jsonIn := `{"a":[{"foo":0},{"foo":1},{"foo":1},{"foo":2}]}`
	jsonOut := `{"a":{"output":"TEST,"}}`

	cfg := getConfig(spec, false)
	kazaamOut, _ := getTransformTestWrapper(Concat, cfg, jsonIn)

	if kazaamOut != jsonOut {
		t.Error("Transformed data does not match expectation.")
		t.Log("Expected:   ", jsonOut)
		t.Log("Actual:     ", kazaamOut)
		t.FailNow()
	}
}

func TestConcatWithBadSpec(t *testing.T) {
	// Bad spec - "Path" should be "path"
	spec := `{"sources": [{"value": "TEST"}, {"Path": "a[*].bar"}], "targetPath": "a.timestamp", "delim": "," }`
	jsonIn := `{"a":[{"foo":0},{"foo":1},{"foo":1},{"foo":2}]}`
	// bad path should cause the result to be blank
	jsonOut := ""

	cfg := getConfig(spec, false)
	kazaamOut, _ := getTransformTestWrapper(Concat, cfg, jsonIn)

	if kazaamOut != jsonOut {
		t.Error("Transformed data does not match expectation.")
		t.Log("Expected:   ", jsonOut)
		t.Log("Actual:     ", kazaamOut)
		t.FailNow()
	}
}

func TestConcatWithMultiMulti(t *testing.T) {
	spec := `{"sources": [{"value": "BEGIN"}, {"path": "a[*].foo"}, {"value": "END"}], "targetPath": "a.output", "delim": "," }`
	jsonIn := `{"a":[{"foo":0},{"foo":1},{"foo":1},{"foo":2}]}`
	jsonOut := `{"a":{"output":"BEGIN,0112,END"}}`

	cfg := getConfig(spec, false)
	kazaamOut, _ := getTransformTestWrapper(Concat, cfg, jsonIn)

	if kazaamOut != jsonOut {
		t.Error("Transformed data does not match expectation.")
		t.Log("Expected:   ", jsonOut)
		t.Log("Actual:     ", kazaamOut)
		t.FailNow()
	}
}

func TestConcatWithLargeNumbers(t *testing.T) {
	spec := `{"sources": [{"path": "a.timestamp"}], "targetPath": "a.output" }`
	jsonIn := `{"a":{"timestamp":1481305274100000000000000000000}}`
	jsonOut := `{"a":{"timestamp":1481305274100000000000000000000,"output":"1481305274100000000000000000000"}}`

	cfg := getConfig(spec, false)
	kazaamOut, _ := getTransformTestWrapper(Concat, cfg, jsonIn)

	if kazaamOut != jsonOut {
		t.Error("Transformed data does not match expectation.")
		t.Log("Expected:   ", jsonOut)
		t.Log("Actual:     ", kazaamOut)
		t.FailNow()
	}
}
