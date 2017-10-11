package transform

import "testing"

func TestDelete(t *testing.T) {
	spec := `{"paths": ["rating.example"]}`
	jsonOut := `{"rating":{"primary":{"value":3}}}`

	cfg := getConfig(spec, false)
	kazaamOut, err := getTransformTestWrapper(Delete, cfg, testJSONInput)

	if err != nil {
		t.Error("Error in transform (simplejson).")
		t.Log("Error: ", err.Error())
		t.FailNow()
	}

	areEqual, _ := checkJSONBytesEqual(kazaamOut, []byte(jsonOut))
	if !areEqual {
		t.Error("Transformed data does not match expectation.")
		t.Log("Expected:   ", jsonOut)
		t.Log("Actual:     ", string(kazaamOut))
		t.FailNow()
	}
}
