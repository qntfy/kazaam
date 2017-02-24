package transform

import "testing"

func TestDefault(t *testing.T) {
	spec := `{"Range": 5}`
	jsonOut := `{"Range":5,"rating":{"example":{"value":3},"primary":{"value":3}}}`

	cfg := getConfig(spec, false)
	kazaamOut, _ := getTransformTestWrapper(Default, cfg, testJSONInput)

	if kazaamOut != jsonOut {
		t.Error("Transformed data does not match expectation.")
		t.Log("Expected: ", jsonOut)
		t.Log("Actual:   ", kazaamOut)
		t.FailNow()
	}
}
