package transform

import (
	"encoding/json"
	"reflect"
	"testing"
)

const testJSONInput = `{"rating":{"example":{"value":3},"primary":{"value":3}}}`

func getConfig(spec string, require bool) Config {
	var f map[string]interface{}
	json.Unmarshal([]byte(spec), &f)
	return Config{Spec: &f, Require: require}
}

func getTransformTestWrapper(tform func(spec *Config, data []byte) ([]byte, error), cfg Config, input string) ([]byte, error) {
	output, e := tform(&cfg, []byte(input))
	if e != nil {
		return nil, e
	}
	return output, nil
}

func checkJSONBytesEqual(item1, item2 []byte) (bool, error) {
	var out1, out2 interface{}

	err := json.Unmarshal(item1, &out1)
	if err != nil {
		return false, nil
	}

	err = json.Unmarshal(item2, &out2)
	if err != nil {
		return false, nil
	}

	return reflect.DeepEqual(out1, out2), nil
}

func TestCheckJSONBytesAreEqual(t *testing.T) {
	item1 := []byte(`{"test":["data1", "data2"],"key":"value"}`)
	item2 := []byte(`{"key":"value","test":["data1", "data2"]}`)
	areEqual, _ := checkJSONBytesEqual(item1, item2)
	if !areEqual {
		t.Error("JSON equality check failed")
		t.Log("Item 1: ", string(item1))
		t.Log("Item 2: ", string(item2))
		t.FailNow()
	}
}

func TestCheckJSONBytesAreNotEqual(t *testing.T) {
	item1 := []byte(`{"test":["data1", "data2"]}`)
	item2 := []byte(`{"test":["data1", "data1"]}`)
	areEqual, _ := checkJSONBytesEqual(item1, item2)
	if areEqual {
		t.Error("JSON inequality check failed")
		t.Log("Item 1: ", string(item1))
		t.Log("Item 2: ", string(item2))
		t.FailNow()
	}
}

func TestBookend(t *testing.T) {
	input := []byte(`"foo", "bar"`)
	expected := []byte(`["foo", "bar"]`)

	result := bookend(input, '[', ']')
	if string(result) != string(expected) {
		t.Error("Bookend result does not match expectation.")
		t.Log("Expected: ", expected)
		t.Log("Actual:   ", result)
		t.FailNow()
	}

	input = []byte("fooString")
	expected = []byte(`"fooString"`)
	result = bookend(input, '"', '"')
	areEqual, _ := checkJSONBytesEqual(result, expected)
	if !areEqual {
		t.Error("Bookend result does not match expectation.")
		t.Log("Expected: ", expected)
		t.Log("Actual:   ", result)
		t.FailNow()
	}
}

func TestSetJSONRaw(t *testing.T) {
	setPathTests := []struct {
		inputData      []byte
		inputValue     []byte
		path           string
		expectedOutput []byte
	}{
		{[]byte(`{"data":"value"}`), []byte(`"newValue"`), "data", []byte(`{"data":"newValue"}`)},
		{[]byte(`{"data":["value", "notValue"]}`), []byte(`"newValue"`), "data[0]", []byte(`{"data":["newValue", "notValue"]}`)},
		{[]byte(`{"data":["value", "notValue"]}`), []byte(`"newValue"`), "data[*]", []byte(`{"data":["newValue", "newValue"]}`)},
		{[]byte(`{"data":[{"key": "value"}, {"key": "value"}]}`), []byte(`"newValue"`), "data[*].key", []byte(`{"data":[{"key": "newValue"}, {"key": "newValue"}]}`)},
		{[]byte(`{"data":[{"key": "value"}, {"key": "value"}]}`), []byte(`"newValue"`), "data[1].key", []byte(`{"data":[{"key": "value"}, {"key": "newValue"}]}`)},
		{[]byte(`{"data":{"subData":[{"key": "value"}, {"key": "value"}]}}`), []byte(`"newValue"`), "data.subData[*].key", []byte(`{"data":{"subData":[{"key": "newValue"}, {"key": "newValue"}]}}`)},
	}
	for _, testItem := range setPathTests {
		actual, _ := setJSONRaw(testItem.inputData, testItem.inputValue, testItem.path)
		areEqual, _ := checkJSONBytesEqual(actual, testItem.expectedOutput)
		if !areEqual {
			t.Error("Error data does not match expectation.")
			t.Log("Expected:   ", testItem.expectedOutput)
			t.Log("Actual:     ", string(actual))
		}
	}
}

func TestSetJSONRawBadIndex(t *testing.T) {
	_, err := setJSONRaw([]byte(`{"data":["value"]}`), []byte(`"newValue"`), "data[g].key")

	errMsg := `Warn: Unable to coerce index to integer: g`
	if err.Error() != errMsg {
		t.Error("Error data does not match expectation.")
		t.Log("Expected:   ", errMsg)
		t.Log("Actual:     ", err.Error())
		t.FailNow()
	}
}

func TestGetJSONRaw(t *testing.T) {
	getPathTests := []struct {
		inputData      []byte
		path           string
		required       bool
		expectedOutput []byte
	}{
		{[]byte(`{"data":"value"}`), "data", true, []byte(`"value"`)},
		{[]byte(`{"data":"value"}`), "data", false, []byte(`"value"`)},
		{[]byte(`{"notData":"value"}`), "data", false, []byte(`null`)},
		{[]byte(`{"data":["value", "notValue"]}`), "data[0]", true, []byte(`"value"`)},
		{[]byte(`{"data":["value", "notValue"]}`), "data[*]", true, []byte(`["value","notValue"]`)},
		{[]byte(`{"data":[{"key": "value"}, {"key": "value"}]}`), "data[*].key", true, []byte(`["value","value"]`)},
		{[]byte(`{"data":[{"key": "value"}, {"key": "otherValue"}]}`), "data[1].key", true, []byte(`"otherValue"`)},
		{[]byte(`{"data":{"subData":[{"key": "value"}, {"key": "value"}]}}`), "data.subData[*].key", true, []byte(`["value","value"]`)},
	}
	for _, testItem := range getPathTests {
		actual, _ := getJSONRaw(testItem.inputData, testItem.path, testItem.required)
		areEqual, _ := checkJSONBytesEqual(actual, testItem.expectedOutput)
		if !areEqual {
			t.Error("Error data does not match expectation.")
			t.Log("Expected:   ", string(testItem.expectedOutput))
			t.Log("Actual:     ", string(actual))
		}
	}
}

func TestGetJSONRawBadIndex(t *testing.T) {
	_, err := getJSONRaw([]byte(`{"data":["value"]}`), "data[-1].key", true)

	errMsg := `Warn: Unable to coerce index to integer: -1`
	if err.Error() != errMsg {
		t.Error("Error data does not match expectation.")
		t.Log("Expected:   ", errMsg)
		t.Log("Actual:     ", err.Error())
		t.FailNow()
	}
}
