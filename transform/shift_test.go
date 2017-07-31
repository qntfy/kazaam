package transform

import "testing"

func TestShift(t *testing.T) {
	jsonOut := `{"Rating":3,"example":{"old":{"value":3}}}`
	spec := `{"Rating": "rating.primary.value","example.old": "rating.example"}`

	cfg := getConfig(spec, false)
	kazaamOut, _ := getTransformTestWrapper(Shift, cfg, testJSONInput)
	areEqual, _ := checkJSONBytesEqual(kazaamOut, []byte(jsonOut))

	if !areEqual {
		t.Error("Transformed data does not match expectation.")
		t.Log("Expected:   ", jsonOut)
		t.Log("Actual:     ", kazaamOut)
		t.FailNow()
	}
}

func TestShiftWithWildcard(t *testing.T) {
	spec := `{"outputArray": "docs[*].data.key"}`
	jsonIn := `{"docs": [{"data": {"key": "val1"}},{"data": {"key": "val2"}}]}`
	jsonOut := `{"outputArray":["val1","val2"]}`

	cfg := getConfig(spec, false)
	kazaamOut, _ := getTransformTestWrapper(Shift, cfg, jsonIn)
	areEqual, _ := checkJSONBytesEqual(kazaamOut, []byte(jsonOut))

	if !areEqual {
		t.Error("Transformed data does not match expectation.")
		t.Log("Expected:   ", jsonOut)
		t.Log("Actual:     ", kazaamOut)
		t.FailNow()
	}
}

func TestShiftWithWildcardEmptySlice(t *testing.T) {
	spec := `{"outputArray": "docs[*].data.key"}`
	jsonIn := `{"docs": []}`
	jsonOut := `{"outputArray":[]}`

	cfg := getConfig(spec, false)
	kazaamOut, _ := getTransformTestWrapper(Shift, cfg, jsonIn)
	areEqual, _ := checkJSONBytesEqual(kazaamOut, []byte(jsonOut))

	if !areEqual {
		t.Error("Transformed data does not match expectation.")
		t.Log("Expected:   ", jsonOut)
		t.Log("Actual:     ", kazaamOut)
		t.FailNow()
	}
}

func TestShiftWithMissingKey(t *testing.T) {
	spec := `{"Rating": "rating.primary.missing_value","example.old": "rating.example"}`
	jsonOut := `{"Rating":null,"example":{"old":{"value":3}}}`

	cfg := getConfig(spec, false)
	kazaamOut, _ := getTransformTestWrapper(Shift, cfg, testJSONInput)
	areEqual, _ := checkJSONBytesEqual(kazaamOut, []byte(jsonOut))

	if !areEqual {
		t.Error("Transformed data does not match expectation.")
		t.Log("Expected:   ", jsonOut)
		t.Log("Actual:     ", kazaamOut)
		t.FailNow()
	}
}

func TestShiftDeepExistsRequire(t *testing.T) {
	testJSONInput := `{"rating":{"example":[{"array":[{"value":3}]},{"another":"object"}]}}`
	spec := `{"example_res":"rating.example[0].array[*].value"}`
	jsonOut := `{"example_res":[3]}`

	cfg := getConfig(spec, true)
	kazaamOut, _ := getTransformTestWrapper(Shift, cfg, testJSONInput)
	areEqual, _ := checkJSONBytesEqual(kazaamOut, []byte(jsonOut))

	if !areEqual {
		t.Error("Transformed data does not match expectation.")
		t.Log("Expected:   ", jsonOut)
		t.Log("Actual:     ", kazaamOut)
		t.FailNow()
	}
}

func TestShiftShallowExistsRequire(t *testing.T) {
	spec := `{"Rating": "not_a_field"}`

	cfg := getConfig(spec, true)
	_, err := getTransformTestWrapper(Shift, cfg, testJSONInput)

	if err == nil {
		t.Error("Transform path does not exist in message and should throw an error")
		t.FailNow()
	}
}

func TestShiftDeepArraysRequire(t *testing.T) {
	spec := `{"Rating": "rating.does[0].not[*].exist"}`

	cfg := getConfig(spec, true)
	_, err := getTransformTestWrapper(Shift, cfg, testJSONInput)

	if err == nil {
		t.Error("Transform path does not exist in message and should throw an error")
		t.FailNow()
	}
}

func TestShiftDeepNoArraysRequire(t *testing.T) {
	spec := `{"Rating": "rating.does.not.exist"}`

	cfg := getConfig(spec, true)
	_, err := getTransformTestWrapper(Shift, cfg, testJSONInput)

	if err == nil {
		t.Error("Transform path does not exist in message and should throw an error")
		t.FailNow()
	}
}

func TestShiftWithEncapsulate(t *testing.T) {
	jsonOut := `{"data":[{"rating":{"example":{"value":3},"primary":{"value":3}}}]}`
	spec := `{"data": ["$"]}`

	cfg := getConfig(spec, false)
	kazaamOut, _ := getTransformTestWrapper(Shift, cfg, testJSONInput)
	areEqual, _ := checkJSONBytesEqual(kazaamOut, []byte(jsonOut))

	if !areEqual {
		t.Error("Transformed data does not match expectation.")
		t.Log("Expected:   ", jsonOut)
		t.Log("Actual:     ", kazaamOut)
		t.FailNow()
	}
}

func TestShiftWithNullSpecValue(t *testing.T) {
	spec := `{"id": null}`
	jsonIn := `{"data": {"id": true}}`

	cfg := getConfig(spec, false)
	_, err := getTransformTestWrapper(Shift, cfg, jsonIn)

	errMsg := `Warn: Unknown type in message for key: id`
	if err.Error() != errMsg {
		t.Error("Error data does not match expectation.")
		t.Log("Expected:   ", errMsg)
		t.Log("Actual:     ", err.Error())
		t.FailNow()
	}
}

func TestShiftWithNullArraySpecValue(t *testing.T) {
	spec := `{"id": [null, "abc"]}`
	jsonIn := `{"data": {"id": true}}`

	cfg := getConfig(spec, false)
	_, err := getTransformTestWrapper(Shift, cfg, jsonIn)

	errMsg := `Warn: Unable to coerce element to json string: <nil>`
	if err.Error() != errMsg {
		t.Error("Error data does not match expectation.")
		t.Log("Expected:   ", errMsg)
		t.Log("Actual:     ", err.Error())
		t.FailNow()
	}
}

func TestShiftWithEndArrayAccess(t *testing.T) {
	spec := `{"id": "docs[1].data[0]"}`
	jsonIn := `{"docs": [{"data": ["abc", "def"]},{"data": ["ghi", "jkl"]}]}`
	jsonOut := `{"id":"ghi"}`

	cfg := getConfig(spec, false)
	kazaamOut, err := getTransformTestWrapper(Shift, cfg, jsonIn)

	if err != nil {
		t.Error("Error on transform.")
		t.Log("Expected: ", jsonOut)
		t.Log("Error: ", err.Error())
		t.FailNow()
	}

	areEqual, _ := checkJSONBytesEqual(kazaamOut, []byte(jsonOut))
	if !areEqual {
		t.Error("Transformed data does not match expectation.")
		t.Log("Expected:   ", jsonOut)
		t.Log("Actual:     ", kazaamOut)
		t.FailNow()
	}
}
