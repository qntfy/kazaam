package transform

import (
	"encoding/json"
	"testing"
)

const testJSONInput = `{"rating":{"example":{"value":3},"primary":{"value":3}}}`

func getConfig(spec string, require bool) Config {
	var f map[string]interface{}
	json.Unmarshal([]byte(spec), &f)
	return Config{Spec: &f, Require: require}
}

func getTransformTestWrapper(tform func(spec *Config, data []byte) ([]byte, error), cfg Config, input string) (string, error) {
	output, e := tform(&cfg, []byte(input))
	if e != nil {
		return "", e
	}
	return string(output), nil
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
	if string(result) != string(expected) {
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
		expectedOutput string
	}{
		{[]byte(`{"data":"value"}`), []byte(`"newValue"`), "data", `{"data":"newValue"}`},
		{[]byte(`{"data":["value", "notValue"]}`), []byte(`"newValue"`), "data[0]", `{"data":["newValue", "notValue"]}`},
		{[]byte(`{"data":["value", "notValue"]}`), []byte(`"newValue"`), "data[*]", `{"data":["newValue", "newValue"]}`},
		{[]byte(`{"data":[{"key": "value"}, {"key": "value"}]}`), []byte(`"newValue"`), "data[*].key", `{"data":[{"key": "newValue"}, {"key": "newValue"}]}`},
		{[]byte(`{"data":[{"key": "value"}, {"key": "value"}]}`), []byte(`"newValue"`), "data[1].key", `{"data":[{"key": "value"}, {"key": "newValue"}]}`},
		{[]byte(`{"data":{"subData":[{"key": "value"}, {"key": "value"}]}}`), []byte(`"newValue"`), "data.subData[*].key", `{"data":{"subData":[{"key": "newValue"}, {"key": "newValue"}]}}`},
	}
	for _, testItem := range setPathTests {
		actual, _ := setJSONRaw(testItem.inputData, testItem.inputValue, testItem.path)
		if string(actual) != testItem.expectedOutput {
			t.Error("Error data does not match expectation.")
			t.Log("Expected:   ", testItem.expectedOutput)
			t.Log("Actual:     ", string(actual))
		}
	}
}

func TestGetJSONRaw(t *testing.T) {
	getPathTests := []struct {
		inputData      []byte
		path           string
		required       bool
		expectedOutput string
	}{
		{[]byte(`{"data":"value"}`), "data", true, `"value"`},
		{[]byte(`{"data":"value"}`), "data", false, `"value"`},
		{[]byte(`{"notData":"value"}`), "data", false, `null`},
		{[]byte(`{"data":["value", "notValue"]}`), "data[0]", true, `"value"`},
		{[]byte(`{"data":["value", "notValue"]}`), "data[*]", true, `[value,notValue]`},
		{[]byte(`{"data":[{"key": "value"}, {"key": "value"}]}`), "data[*].key", true, `["value","value"]`},
		{[]byte(`{"data":[{"key": "value"}, {"key": "otherValue"}]}`), "data[1].key", true, `"otherValue"`},
		{[]byte(`{"data":{"subData":[{"key": "value"}, {"key": "value"}]}}`), "data.subData[*].key", true, `["value","value"]`},
	}
	for _, testItem := range getPathTests {
		actual, _ := getJSONRaw(testItem.inputData, testItem.path, testItem.required)
		if string(actual) != testItem.expectedOutput {
			t.Error("Error data does not match expectation.")
			t.Log("Expected:   ", testItem.expectedOutput)
			t.Log("Actual:     ", string(actual))
		}
	}
}
