package transform

import (
	"github.com/mbordner/kazaam/registry"
	"testing"
)

type ExprConverterTest struct{}

func (c *ExprConverterTest) Init(config []byte) (err error) {
	return
}
func (c *ExprConverterTest) Convert(jsonData []byte, value []byte, args []byte) (newValue []byte, err error) {
	v,e := NewJSONValue(args)
	if e != nil {
		err = e
	} else {
		v,e := NewJSONValue([]byte(v.GetStringValue()))
		if e != nil {
			err = e
		} else {
			if v.GetType() == JSONString {
				newValue = args
			} else {
				newValue =[]byte(v.String())
			}
		}
	}
	return
}

func TestExpressions(t *testing.T) {

	registry.RegisterConverter("exprConvA", &ExprConverterTest{})
	registry.RegisterConverter("exprConvB", &ExprConverterTest{})

	data := []byte(`
{
  "tests": {
    "test_int": 500,
    "test_float": 500.01,
    "test_float2": 500.0,
    "test_fraction": 0.5,
    "test_trim": "    blah   ",
    "test_money": "$6,000,000",
    "test_chars": "abcdefghijklmnopqrstuvwxyz",
	"test_mapped": "Texas",
	"test_null": null
  },
  "test_bool": true
}
`)

	table := []struct {
		expr          string
		expected      bool
		expectedError bool
	}{
		{
			`  ( tests.test_int == 500 ) `,
			true,
			false,
		},
		{
			`tests.test_int > 400 `,
			true,
			false,
		},
		{
			`tests.test_int < 400`,
			false,
			false,
		},
		{
			`tests.test_int >= 500`,
			true,
			false,
		},
		{
			`tests.test_int <= 500`,
			true,
			false,
		},
		{
			`!(tests.test_int != 500)`,
			true,
			false,
		},
		{
			`tests.test_float == 500 || tests.test_int == 500`,
			true,
			false,
		},
		{
			`tests.test_float2 == tests.test_int`,
			true,
			false,
		},
		{
			`true && false`,
			false,
			false,
		},
		{
			`false || true`,
			true,
			false,
		},
		{
			`"string1" != "string2"`,
			true,
			false,
		},
		{
			`"string1" == "string2" || (exprConvA("tests.test_mapped","\"Texas\"") == "Texas") && true || false`,
			true,
			false,
		},
		{
			`1 && 1`,
			false, // && and || needs boolean expressions
			true,
		},
		{
			`"string1" == 1`,
			false,
			true,
		},
		{
			`1 ^ 1`,
			false,
			true,
		},
		{
			`50.0 == 50`,
			true,
			false,
		},
		{
			`exprConvA("tests.test_int","null") == null`,
			true,
			false,
		},
		{
			`tests.test_null == null && tests.test_null == nil`,
			true,
			false,
		},
		{
			`tests.test_chars == "abcdefghijklmnopqrstuvwxyz"`,
			true,
			false,
		},
		{
			`exprConvA("tests.test_int","true") == true`,
			true,
			false,
		},
		{
			`exprConvA("tests.test_int","5") == 5`,
			true,
			false,
		},
		{
			`exprConvA("tests.test_int","5.01") == 5.01`,
			true,
			false,
		},
		{
			`blah("tests.test_int","test") == true`,
			false,
			true,
		},
		{
			`exprConvA("tests.test_int",1) == 1`, // will fail because the arguments are not a string, and this is required
			false,
			true,
		},
		{
			`exprConvA(1,1) == 1`, // will fail because the path is not a string, and this is required
			false,
			true,
		},
		{
			`exprConvA() == 1`, // will fail because expects path string
			false,
			true,
		},
	}

	for _, test := range table {
		be, err := NewBasicExpr(data, test.expr)
		if err != nil {
			t.Error("error parsing expression: {}", test.expr)
		}

		if evaluation, err := be.Eval(); err != nil || evaluation != test.expected {
			if err != nil && test.expectedError == false {
				t.Error("unexpected expression evaluation value")
			}
		}
	}

}
