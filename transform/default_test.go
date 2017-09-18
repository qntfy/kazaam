package transform

import "testing"

func TestDefault(t *testing.T) {
	spec := `{"Range": 5}`
	jsonOut := `{"rating":{"example":{"value":3},"primary":{"value":3}},"Range":5}`

	cfg := getConfig(spec, false)
	kazaamOut, err := getTransformTestWrapper(Default, cfg, testJSONInput)

	if err != nil {
		t.Error("Error in transform (simplejson).")
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

func TestDefaultArrayAppend(t *testing.T) {
	spec := `{"rating.example[+]": 5}`
	jsonInput := `{"rating":{"example":[3],"primary":{"value":3}}}`
	jsonOut := `{"rating":{"example":[3,5],"primary":{"value":3}}}`

	cfg := getConfig(spec, false)
	kazaamOut, err := getTransformTestWrapper(Default, cfg, jsonInput)

	if err != nil {
		t.Error("Error in transform (simplejson).")
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
