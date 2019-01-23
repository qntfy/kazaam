package converter

import (
	"github.com/mbordner/kazaam/registry"
	"strconv"
	"testing"
)

func TestSplitn_Convert(t *testing.T) {

	registry.RegisterConverter("splitn", &Splitn{})
	c := registry.GetConverter("splitn")

	table := []struct {
		value     string
		arguments string
		expected  string
	}{
		{`"aazbbzcczdd""`, `z 4`, `"dd"`,},
		{`"abc|def|ghi|jkl|mno"`, `| 2`, `"def"`,},
		{"\"abc\\ndef\\nghi\\njkl\\nmno\"", "\n 5", `"mno"`,},
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
