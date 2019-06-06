package converter

import (
	"strconv"
	"testing"
	"github.com/mbordner/kazaam/registry"
)

func TestEqs_Convert(t *testing.T) {

	registry.RegisterConverter("eqs", &Eqs{})
	c := registry.GetConverter("eqs")

	table := []struct {
		value     string
		arguments string
		expected  string
	}{
		{`"The quick brown fox jumps over the lazy dog"`, `"The quick brown fox jumps over the lazy dog"`, `true`,},
		{`42`,`42`,`true`,},
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