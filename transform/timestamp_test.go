package transform

import (
	"testing"
)

func TestTimestamp(t *testing.T) {
	jsonIn := `{"timestampA":"Sun Jul 23 08:15:27 +0000 2017","topLevel":{"timestampB":"Fri Jul 21 08:15:27 +0000 2017"},"timestampC":["Sat Jul 22 08:15:27 +0000 2017","Sun Jul 23 08:15:27 +0000 2017","Mon Jul 24 08:15:27 +0000 2017"]}`
	spec := `{"timestampA":{"inputFormat":"Mon Jan _2 15:04:05 -0700 2006","outputFormat":"2006-01-02T15:04:05-0700"},"topLevel.timestampB":{"inputFormat":"Mon Jan _2 15:04:05 -0700 2006","outputFormat":"2006-01-02T15:04:05-0700"},"timestampC":{"inputFormat":"Mon Jan _2 15:04:05 -0700 2006","outputFormat":"2006-01-02T15:04:05-0700"}}`
	jsonOut := `{"timestampA":"2017-07-23T08:15:27+0000","topLevel":{"timestampB":"2017-07-21T08:15:27+0000"},"timestampC":["2017-07-22T08:15:27+0000","2017-07-23T08:15:27+0000","2017-07-24T08:15:27+0000"]}`

	cfg := getConfig(spec, false)
	kazaamOut, _ := getTransformTestWrapper(Timestamp, cfg, jsonIn)

	if kazaamOut != jsonOut {
		t.Error("Transformed data does not match expectation.")
		t.Log("Expected:   ", jsonOut)
		t.Log("Actual:     ", kazaamOut)
		t.FailNow()
	}
}

func TestTimestampWithIndex(t *testing.T) {
	jsonIn := `{"timestampC":["Sat Jul 22 08:15:27 +0000 2017","Sun Jul 23 08:15:27 +0000 2017","Mon Jul 24 08:15:27 +0000 2017"]}`
	spec := `{"timestampC[0]":{"inputFormat":"Mon Jan _2 15:04:05 -0700 2006","outputFormat":"2006-01-02T15:04:05-0700"}}`
	jsonOut := `{"timestampC":["2017-07-22T08:15:27+0000","Sun Jul 23 08:15:27 +0000 2017","Mon Jul 24 08:15:27 +0000 2017"]}`

	cfg := getConfig(spec, false)
	kazaamOut, _ := getTransformTestWrapper(Timestamp, cfg, jsonIn)

	if kazaamOut != jsonOut {
		t.Error("Transformed data does not match expectation.")
		t.Log("Expected:   ", jsonOut)
		t.Log("Actual:     ", kazaamOut)
		t.FailNow()
	}
}

func TestTimestampWithWildcard(t *testing.T) {
	jsonIn := `{"timestampC":["Sat Jul 22 08:15:27 +0000 2017","Sun Jul 23 08:15:27 +0000 2017","Mon Jul 24 08:15:27 +0000 2017"]}`
	spec := `{"timestampC[*]":{"inputFormat":"Mon Jan _2 15:04:05 -0700 2006","outputFormat":"2006-01-02T15:04:05-0700"}}`
	jsonOut := `{"timestampC":["2017-07-22T08:15:27+0000","2017-07-23T08:15:27+0000","2017-07-24T08:15:27+0000"]}`

	cfg := getConfig(spec, false)
	kazaamOut, _ := getTransformTestWrapper(Timestamp, cfg, jsonIn)

	if kazaamOut != jsonOut {
		t.Error("Transformed data does not match expectation.")
		t.Log("Expected:   ", jsonOut)
		t.Log("Actual:     ", kazaamOut)
		t.FailNow()
	}
}

func TestTimestampWithMissingKey(t *testing.T) {
	jsonIn := `{"notTheRightField": 9999999,"topLevel":{"timestampB":"Fri Jul 21 08:15:27 +0000 2017"}}`
	spec := `{"timestampA":{"inputFormat":"Mon Jan _2 15:04:05 -0700 2006","outputFormat":"2006-01-02T15:04:05-0700"},"topLevel.timestampB":{"inputFormat":"Mon Jan _2 15:04:05 -0700 2006","outputFormat":"2006-01-02T15:04:05-0700"}}`
	jsonOut := `{"notTheRightField": 9999999,"topLevel":{"timestampB":"2017-07-21T08:15:27+0000"}}`

	cfg := getConfig(spec, false)
	kazaamOut, _ := getTransformTestWrapper(Timestamp, cfg, jsonIn)

	if kazaamOut != jsonOut {
		t.Error("Transformed data does not match expectation.")
		t.Log("Expected:   ", jsonOut)
		t.Log("Actual:     ", kazaamOut)
		t.FailNow()
	}
}

func TestTimestampWithMissingKeyAndRequire(t *testing.T) {
	jsonIn := `{"notTheRightField": 9999999}`
	spec := `{"timestampA":{"inputFormat":"Mon Jan _2 15:04:05 -0700 2006","outputFormat":"2006-01-02T15:04:05-0700"}}`

	cfg := getConfig(spec, true)
	_, err := getTransformTestWrapper(Timestamp, cfg, jsonIn)

	if err == nil {
		t.Error("Transform path does not exist in message and should throw an error")
		t.FailNow()
	}
}

func TestTimestampWithMissingInputFormatKey(t *testing.T) {
	spec := `{"timestampA":{"outputFormat":"2006-01-02T15:04:05-0700"}}`
	jsonIn := `{"data": {"id": true}}`

	cfg := getConfig(spec, false)
	_, err := getTransformTestWrapper(Timestamp, cfg, jsonIn)

	errMsg := "Warn: Invalid spec. Unable to get \"inputFormat\" for key: timestampA"
	if err.Error() != errMsg {
		t.Error("Error data does not match expectation.")
		t.Log("Expected:   ", errMsg)
		t.Log("Actual:     ", err.Error())
		t.FailNow()
	}
}

func TestTimestampWithMissingOutputFormatKey(t *testing.T) {
	spec := `{"timestampA":{"inputFormat":"Mon Jan _2 15:04:05 -0700 2006"}}`
	jsonIn := `{"data": {"id": true}}`

	cfg := getConfig(spec, false)
	_, err := getTransformTestWrapper(Timestamp, cfg, jsonIn)

	errMsg := "Warn: Invalid spec. Unable to get \"outputFormat\" for key: timestampA"
	if err.Error() != errMsg {
		t.Error("Error data does not match expectation.")
		t.Log("Expected:   ", errMsg)
		t.Log("Actual:     ", err.Error())
		t.FailNow()
	}
}

func TestTimestampWithMissingOpsKey(t *testing.T) {
	spec := `{"operations": null}`
	jsonIn := `{"data": {"id": true}}`

	cfg := getConfig(spec, false)
	_, err := getTransformTestWrapper(Timestamp, cfg, jsonIn)

	errMsg := "Warn: Invalid spec. Unable to get value for key: operations"
	if err.Error() != errMsg {
		t.Error("Error data does not match expectation.")
		t.Log("Expected:   ", errMsg)
		t.Log("Actual:     ", err.Error())
		t.FailNow()
	}
}
