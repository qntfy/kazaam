package converter

import (
	"github.com/mbordner/kazaam/registry"
	"strconv"
	"testing"
)

func TestLen_Convert(t *testing.T) {

	registry.RegisterConverter("len", &Len{})
	c := registry.GetConverter("len")

	table := []struct {
		value     string
		arguments string
		expected  string
	}{
		{`"The quick brown fox jumps over the lazy dog"`, ``, `43`,},
		{`"the lazy dog"`, ``, `12`,},
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
