package transform

import "testing"

func TestMerge(t *testing.T) {
	spec := `{"merge1":[{"name":"prop_1","array":"array_a"},{"name":"prop_2","array":"array_b"},{"name":"prop_3","array":"array_c"}]}`
	jsonIn := `{"array_a":["a_1","a_2","a_3"],"array_b":["b_1","b_2","b_3"],"array_c":["c_1","c_2","c_3"]}`
	jsonOut := `{"merge1":[{"prop_1":"a_1","prop_2":"b_1","prop_3":"c_1"},{"prop_1":"a_2","prop_2":"b_2","prop_3":"c_2"},{"prop_1":"a_3","prop_2":"b_3","prop_3":"c_3"}]}`

	cfg := getConfig(spec, false)
	kazaamOut, _ := getTransformTestWrapper(Merge, cfg, jsonIn)
	areEqual, _ := checkJSONBytesEqual(kazaamOut, []byte(jsonOut))

	if !areEqual {
		t.Error("Transformed data does not match expectation.")
		t.Log("Expected:   ", jsonOut)
		t.Log("Actual:     ", kazaamOut)
		t.FailNow()
	}
}

func TestMergeSingleArray(t *testing.T) {
	spec := `{"merge1":[{"name":"prop_1","array":"array_a"}]}`
	jsonIn := `{"array_a":["a_1","a_2","a_3"]}`
	jsonOut := `{"merge1":[{"prop_1":"a_1"},{"prop_1":"a_2"},{"prop_1":"a_3"}]}`

	cfg := getConfig(spec, false)
	kazaamOut, _ := getTransformTestWrapper(Merge, cfg, jsonIn)
	areEqual, _ := checkJSONBytesEqual(kazaamOut, []byte(jsonOut))

	if !areEqual {
		t.Error("Transformed data does not match expectation.")
		t.Log("Expected:   ", jsonOut)
		t.Log("Actual:     ", kazaamOut)
		t.FailNow()
	}
}
