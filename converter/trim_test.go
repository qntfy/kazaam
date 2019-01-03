package converter

import (
	"github.com/mbordner/kazaam/registry"
	"strconv"
	"testing"
)

func TestTrim_Convert(t *testing.T) {
	registry.RegisterConverter("trim", &Trim{})
	c := registry.GetConverter("trim")

	table := []struct {
		value     string
		arguments string
		expected  string
	}{
		{`"    THIS IS A TEST    "`, ``, `"THIS IS A TEST"`,},
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