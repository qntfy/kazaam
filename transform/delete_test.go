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

func TestDeleteSpecErrorNoPathsKey(t *testing.T) {
	spec := `{"pathz": ["a.path"]}`
	expectedErr := "Unable to get paths to delete"

	cfg := getConfig(spec, false)
	_, err := getTransformTestWrapper(Delete, cfg, testJSONInput)

	if err == nil {
		t.Error("Should have generated error for invalid paths")
		t.Log("Spec:   ", spec)
		t.FailNow()
	}
	e, ok := err.(SpecError)
	if !ok {
		t.Error("Unexpected error type")
		t.FailNow()
	}

	if e.Error() != expectedErr {
		t.Error("Unexpected error details")
		t.Log("Expected:   ", expectedErr)
		t.Log("Actual:     ", e.Error())
		t.FailNow()
	}
}

func TestDeleteSpecErrorInvalidPaths(t *testing.T) {
	spec := `{"paths": false}`
	expectedErr := "paths should be a slice of strings: false"

	cfg := getConfig(spec, false)
	_, err := getTransformTestWrapper(Delete, cfg, testJSONInput)

	if err == nil {
		t.Error("Should have generated error for invalid paths")
		t.Log("Spec:   ", spec)
		t.FailNow()
	}
	e, ok := err.(SpecError)
	if !ok {
		t.Error("Unexpected error type")
		t.FailNow()
	}

	if e.Error() != expectedErr {
		t.Error("Unexpected error details")
		t.Log("Expected:   ", expectedErr)
		t.Log("Actual:     ", e.Error())
		t.FailNow()
	}
}

func TestDeleteSpecErrorInvalidPathItem(t *testing.T) {
	spec := `{"paths": ["foo", 42]}`
	expectedErr := "Error processing 42: path should be a string"

	cfg := getConfig(spec, false)
	_, err := getTransformTestWrapper(Delete, cfg, testJSONInput)

	if err == nil {
		t.Error("Should have generated error for invalid paths")
		t.Log("Spec:   ", spec)
		t.FailNow()
	}
	e, ok := err.(SpecError)
	if !ok {
		t.Error("Unexpected error type")
		t.FailNow()
	}

	if e.Error() != expectedErr {
		t.Error("Unexpected error details")
		t.Log("Expected:   ", expectedErr)
		t.Log("Actual:     ", e.Error())
		t.FailNow()
	}
}

func TestDeleteSpecErrorWildcardNotSupported(t *testing.T) {
	spec := `{"paths": ["ratings[*].value"]}`
	jsonIn := `{"ratings: [{"value": 3, "user": "rick"}, {"value": 7, "user": "jerry"}]}`
	expectedErr := "Array wildcard not supported for this operation."

	cfg := getConfig(spec, false)
	_, err := getTransformTestWrapper(Delete, cfg, jsonIn)

	if err == nil {
		t.Error("Should have generated error for invalid paths")
		t.Log("Spec:   ", spec)
		t.FailNow()
	}
	e, ok := err.(SpecError)
	if !ok {
		t.Error("Unexpected error type")
		t.FailNow()
	}

	if e.Error() != expectedErr {
		t.Error("Unexpected error details")
		t.Log("Expected:   ", expectedErr)
		t.Log("Actual:     ", e.Error())
		t.FailNow()
	}
}

func TestDeleteWithRequire(t *testing.T) {
	spec := `{"paths": ["rating.examplez"]}`

	cfg := getConfig(spec, true)
	_, err := getTransformTestWrapper(Delete, cfg, testJSONInput)

	if err == nil {
		t.Error("Should have generated error for invalid paths")
		t.Log("Spec:   ", spec)
		t.FailNow()
	}
	_, ok := err.(RequireError)
	if !ok {
		t.Error("Unexpected error type")
		t.Error(err.Error())
		t.FailNow()
	}

}
