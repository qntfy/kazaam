package transform

import "testing"

func TestCoalesce(t *testing.T) {
	spec := `{"foo": ["rating.foo", "rating.primary"]}`
	jsonOut := `{"foo":{"value":3},"rating":{"example":{"value":3},"primary":{"value":3}}}`
	// unfortunately the libs encode in different order in this case
	altJsonOut := `{"rating":{"example":{"value":3},"primary":{"value":3}},"foo":{"value":3}}`

	cfg := getConfig(spec, false)
	kazaamOut, _ := getTransformTestWrapper(Coalesce, cfg, testJSONInput)
	kazaamOutRaw, _ := getTransformTestWrapperRaw(CoalesceRaw, cfg, testJSONInput)

	if kazaamOut != jsonOut || (kazaamOutRaw != jsonOut && kazaamOutRaw != altJsonOut) {
		t.Error("Transformed data does not match expectation.")
		t.Log("Expected:   ", jsonOut)
		t.Log("Actual:     ", kazaamOut)
		t.Log("Actual Raw: ", kazaamOutRaw)
		t.FailNow()
	}
}

func TestCoalesceWithRequire(t *testing.T) {
	spec := `{"foo": ["rating.foo", "rating.primary"]}`

	cfg := getConfig(spec, true)
	_, err := getTransformTestWrapper(Coalesce, cfg, testJSONInput)
	_, errRaw := getTransformTestWrapperRaw(CoalesceRaw, cfg, testJSONInput)

	if err == nil || errRaw == nil {
		t.Error("Coalesce does not support \"require\" and should throw an error.")
		t.FailNow()
	}
}

func TestCoalesceWithMulti(t *testing.T) {
	spec := `{"foo": ["rating.foo", "rating.primary"], "bar": ["rating.bar", "rating.example.value"]}`
	jsonOut := `{"bar":3,"foo":{"value":3},"rating":{"example":{"value":3},"primary":{"value":3}}}`
	// unfortunately the libs encode in different order in this case
	altJsonOut := `{"rating":{"example":{"value":3},"primary":{"value":3}},"foo":{"value":3},"bar":3}`

	cfg := getConfig(spec, false)
	kazaamOut, _ := getTransformTestWrapper(Coalesce, cfg, testJSONInput)
	kazaamOutRaw, _ := getTransformTestWrapperRaw(CoalesceRaw, cfg, testJSONInput)

	if kazaamOut != jsonOut || (kazaamOutRaw != jsonOut && kazaamOutRaw != altJsonOut) {
		t.Error("Transformed data does not match expectation.")
		t.Log("Expected:   ", jsonOut)
		t.Log("Actual:     ", kazaamOut)
		t.Log("Actual Raw: ", kazaamOutRaw)
		t.FailNow()
	}
}

func TestCoalesceWithNotFound(t *testing.T) {
	spec := `{"foo": ["rating.foo", "rating.bar", "ratings"]}`
	jsonOut := `{"rating":{"example":{"value":3},"primary":{"value":3}}}`

	cfg := getConfig(spec, false)
	kazaamOut, _ := getTransformTestWrapper(Coalesce, cfg, testJSONInput)
	kazaamOutRaw, _ := getTransformTestWrapperRaw(CoalesceRaw, cfg, testJSONInput)

	if kazaamOut != jsonOut || kazaamOutRaw != jsonOut {
		t.Error("Transformed data does not match expectation.")
		t.Log("Expected:   ", jsonOut)
		t.Log("Actual:     ", kazaamOut)
		t.Log("Actual Raw: ", kazaamOutRaw)
		t.FailNow()
	}
}
