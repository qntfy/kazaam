package converter

import (
	"github.com/mbordner/kazaam/registry"
	"strconv"
	"testing"
)

func TestLower_Convert(t *testing.T) {
	registry.RegisterConverter("lower", &Lower{})
	c := registry.GetConverter("lower")

	table := []struct {
		value     string
		arguments string
		expected  string
	}{
		{`"THIS IS A TEST"`, ``, `"this is a test"`,},
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
