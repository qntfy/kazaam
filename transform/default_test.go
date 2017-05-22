package transform

import "testing"

func TestDefault(t *testing.T) {
	spec := `{"Range": 5}`
	jsonOut := `{"Range":5,"rating":{"example":{"value":3},"primary":{"value":3}}}`
	// unfortunately raw and simplejson encode differently in this case
	altJsonOut := `{"rating":{"example":{"value":3},"primary":{"value":3}},"Range":5}`

	cfg := getConfig(spec, false)
	kazaamOut, err := getTransformTestWrapper(Default, cfg, testJSONInput)
	kazaamOutRaw, errRaw := getTransformTestWrapperRaw(DefaultRaw, cfg, testJSONInput)

	if err != nil {
		t.Error("Error in transform (simplejson).")
		t.Log("Error: ", err.Error())
		t.FailNow()
	}
	if errRaw != nil {
		t.Error("Error in transform (raw).")
		t.Log("Error: ", errRaw.Error())
		t.FailNow()
	}
	if kazaamOut != jsonOut || (kazaamOutRaw != jsonOut && kazaamOutRaw != altJsonOut) {
		t.Error("Transformed data does not match expectation.")
		t.Log("Expected:   ", jsonOut)
		t.Log("Actual:     ", kazaamOut)
		t.Log("Actual Raw: ", kazaamOutRaw)
		t.FailNow()
	}
}
