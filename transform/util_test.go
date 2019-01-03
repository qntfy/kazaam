package transform

import (
	"bytes"
	"encoding/json"
	"github.com/mbordner/kazaam/registry"
	"reflect"
	"strconv"
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
		{[]byte(`{"data":"value"}`), []byte(`"newValue"`), "data[1]", []byte(`{"data":[null,"newValue"]}`)},
		{[]byte(`{"data":["value"]}`), []byte(`"newValue"`), "data[-].key", []byte(`{"data":[{"key":"newValue"},"value"]}`)},
		{[]byte(`{"data":["value"]}`), []byte(`"newValue"`), "data[+]", []byte(`{"data":["value","newValue"]}`)},
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

func TestGetJsonPathValue(t *testing.T) {
	data := []byte(`{"data":{"subData":[{"key": "value"}, {"key": "value"}]}}`)
	jv, err := GetJsonPathValue(data, `data.subData[0].key`)

	if err != nil {
		t.Error("failed to get JSONValue")
		t.FailNow()
	} else {
		if jv.IsString() == false {
			t.Error("unexpected json value")
		}
	}
}

func TestJSONValueParsing(t *testing.T) {

	table := []struct {
		data     string
		dataType int
	}{
		{`null`, JSONNull},
		{`"this is a string"`, JSONString},
		{`42`, JSONInt},
		{`3.14`, JSONFloat},
		{`true`, JSONBool},
		{`false`, JSONBool},
	}

	pi := `3.141592653`
	jv, e := NewJSONValue([]byte(pi))
	if e != nil {
		t.Error("failed parsing [{}]", pi)
		t.FailNow()
	}

	jv.SetFloatStringPrecision(4)
	if jv.String() != "3.1416" {
		t.Error("float print precision issue with rounding")
	}

	for _, test := range table {
		jv, e := NewJSONValue([]byte(test.data))
		if e != nil {
			t.Error("failed parsing [{}]", test.data)
			t.FailNow()
		} else {
			if jv.GetType() != test.dataType {
				t.Error("json value data type mismatch")
				t.Log("Exepcted: {}", test.dataType)
				t.Log("Actual: {}", jv.GetType())
			}

			jv.SetFloatStringPrecision(2)

			valBytes, _ := json.Marshal(jv.GetValue())
			if bytes.Compare(valBytes, jv.GetData()) != 0 {
				t.Error("original bytes data not matching json marshal")
			}

			if string(valBytes) != test.data {
				t.Error("GetValue didnt return the original value")
			}

			if test.data != jv.String() {
				t.Error("expected String() call to produce original string")
				t.Log("Expected: {}", test.data)
				t.Log("Actual: {}", jv.String())
			}

			switch jv.GetType() {
			case JSONNull:
				if !jv.IsNull() {
					t.Error("null test function failed")
				}
			case JSONString:
				if !jv.IsString() {
					t.Error("string test function failed")
				}
				if jv.GetQuotedStringValue() != test.data {
					t.Error("not returning expected quoted string value")
					t.Log("Expected: {}", test.data)
					t.Log("Actual: {}", jv.GetQuotedStringValue())
				}
				tmp, _ := strconv.Unquote(test.data)
				if jv.GetStringValue() != tmp {
					t.Error("not returning expected string value")
					t.Log("Expected: {}", test.data)
					t.Log("Actual: {}", jv.GetStringValue())
				}
			case JSONInt:
				if !jv.IsNumber() {
					t.Error("number test function failed for int")
				}
				tmp, _ := strconv.ParseInt(test.data, 10, 64)
				if jv.GetIntValue() != tmp {
					t.Error("not returning expected int value")
				}
				if jv.GetNumber().String() != jv.String() {
					t.Error("expected number string to be same as jv string")
				}
			case JSONFloat:
				if !jv.IsNumber() {
					t.Error("number test function failed for int")
				}
				tmp, _ := strconv.ParseFloat(test.data, 64)
				if jv.GetFloatValue() != tmp {
					t.Error("not returning expected float value")
				}
			case JSONBool:
				if !jv.IsBool() {
					t.Error("bool test function failed")
				}
				if jv.GetBoolValue() && test.data == "false" || !jv.GetBoolValue() && test.data == "true" {
					t.Error("failed getting bool value")
				}
			}

		}
	}

}

func TestUnescapeString(t *testing.T) {

	s := `\ blah blah\ \n\t\\ \`

	if unescapeString(s) != ` blah blah nt\ ` {
		t.Error("unexpected behavior from unescapeString")
	}
}

type ConverterTest struct{}

func (c *ConverterTest) Init(config []byte) (err error) {
	return
}
func (c *ConverterTest) Convert(jsonData []byte, value []byte, args []byte) (newValue []byte, err error) {
	newValue = args
	return
}

func TestJsonPathParameters(t *testing.T) {

	registry.RegisterConverter("convA", &ConverterTest{})
	registry.RegisterConverter("convB", &ConverterTest{})

	data := []byte(`
{
  "tests": {
    "test_int": 500,
    "test_float": 500.01,
    "test_fraction": 0.5,
    "test_trim": "    blah   ",
    "test_money": "$6,000,000",
    "test_chars": "abcdefghijklmnopqrstuvwxyz",
	"test_mapped": "Texas"
  },
  "test_bool": true
}
`)

	table := []struct {
		path       string
		expectSkip bool
		expected   interface{}
	}{
		{
			`path.not.found?`,
			true,
			nil,
		},
		{
			`path.not.found? 1 `,
			false,
			1,
		},
		{
			`path.not.found? -2.2`,
			false,
			-2.2,
		},
		{
			`path.not.found? "blah"`,
			false,
			"blah",
		},
		{
			`path.not.found? true`,
			false,
			true,
		},
		{
			`tests.test_float ? tests.test_int == 500 && convA("tests.test_trim","blah") == "bleh" :   `,
			true,
			nil,
		},
		{
			`tests.test_float ? tests.test_int == 500 && convA("path.not.found","blah") == "blah" :  "expression error, so returns default value, even though exists"  `,
			false,
			"expression error, so returns default value, even though exists",
		},
		{
			`path.not.found ? invalid_expr( can't even parse : "default value because expression had syntax errors and is treated as false evaluation"  `,
			false,
			"default value because expression had syntax errors and is treated as false evaluation",
		},
		{
			`tests.test_float ? invalid_expr( :  `,
			true,
			nil,
		},
		{
			`path.not.found ? invalid_expr( but forgot colon so it's treated like default value, and skipped cause it's invalid json'  `,
			true,
			nil,
		},
		{
			`tests.test_float ? (tests.test_int == 500 && convA("tests.test_trim","blah") == "bleh") && true : "default value"  `,
			false,
			"default value",
		},
		{
			`tests.test_money ? (tests.test_int == 500 && test_bool ) : "$7,000,000" | convA test1 | convB test2`,
			false,
			"test2",
		},
		{ // white space is ignored around the arguments, unless escaped with a slash.. NOTE in json, this would require extra \ for escaping
			`tests.test_money ? (tests.test_int == 500 && test_bool ) && (convA("tests.test_trim","blah") == "blah") : "$7,000,000" | convA test1 | convB    \ test2 \    `,
			false,
			" test2  ",
		},
	}

	for _, test := range table {

		val, err := GetJSONRaw(data, test.path, true)

		if err != nil {
			if _, ok := err.(CPathSkipError); ok { // was a conditional path, and no default
				if test.expectSkip == false {
					t.Error("unexpected conditional skip error")
				}

				if string(err.Error()) != "Conditional Path missing and without a default value" {
					t.Error("unexpected cpath error message")
				}

			} else {
				t.Error("unexpected error parsing json path [{}]", test.path)
			}
		} else {
			expextedBytes, _ := json.Marshal(test.expected)

			if bytes.Compare(val, expextedBytes) != 0 {
				t.Error("value {} doesn't match expected {}", string(val), test.expected)
			}

		}

	}

}
