package transform

import "testing"

func TestSteps(t *testing.T) {
	jsonOut := `{"example":{"old":{"value":3}},"Rating":3}`
	spec := `{"steps":[{"example.old": "rating.example"},{"Rating": "example.old.value"}]}`

	cfg := getConfig(spec, false)
	kazaamOut, _ := getTransformTestWrapper(Steps, cfg, `{"rating":{"example":{"value":3},"primary":{"value":3}}}`)
	areEqual, _ := checkJSONBytesEqual(kazaamOut, []byte(jsonOut))

	if !areEqual {
		t.Error("Transformed data does not match expectation.")
		t.Log("Expected:   ", jsonOut)
		t.Log("Actual:     ", string(kazaamOut))
		t.FailNow()
	}
}