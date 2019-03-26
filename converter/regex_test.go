package converter

import (
	"github.com/mbordner/kazaam/registry"
	"strconv"
	"testing"
)

func TestRegex_Convert(t *testing.T) {
	registry.RegisterConverter("regex", &Regex{})
	c := registry.GetConverter("regex")

	c.Init([]byte(`
{
  "remove_dollar_sign": {
    "match": "\\$\\s*(.*)",
    "replace": "$1"
  },
  "remove_comma": {
    "match": ",",
    "replace": ""
  },
  "convert_naics": [
	{
		"match": "^8111.*",
 		"replace": "automotive_services"
	},
	{
		"match": "^4413.*",
 		"replace": "automotive_services"
	},
	{
		"match": "^531.*",
		"replace": "real_estate"
	},
	{
		"match": "real_estate",
		"replace": "did not stop when matched"
	}
  ]
}
`))

	table := []struct {
		value     string
		arguments string
		expected  string
	}{
		{`"$5,000,000"`, `remove_dollar_sign`, `"5,000,000"`},
		{`"5,000,000"`, `remove_comma`, `"5000000"`},
		{`"531312"`, `convert_naics`, `"real_estate"`},
	}

	for _, test := range table {
		v, e := c.Convert(getTestData(), []byte(test.value), []byte(strconv.Quote(test.arguments)))

		if e != nil {
			t.Error("error running convert function")
		}

		if string(v) != test.expected {
			t.Error("unexpected result from convert")
			t.Log("Expected: {}", test.expected)
			t.Log("Actual: {}", string(v))
		}
	}
}
