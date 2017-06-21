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
	if kazaamOut != jsonOut {
		t.Error("Transformed data does not match expectation.")
		t.Log("Expected:   ", jsonOut)
		t.Log("Actual:     ", kazaamOut)
		t.FailNow()
	}
}
