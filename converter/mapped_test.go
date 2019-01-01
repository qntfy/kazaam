package converter

import (
	"github.com/qntfy/kazaam/registry"
	"strconv"
	"testing"
)

func TestMapped_Convert(t *testing.T) {
	registry.RegisterConverter("mapped", &Mapped{})
	c := registry.GetConverter("mapped")

	c.Init([]byte(`
{
  "states": {
    "Ohio": "OH",
    "Texas": "TX"
  }
}
`))

	table := []struct {
		value     string
		arguments string
		expected  string
	}{
		{`"Ohio"`,`states`,`"OH"`},
		{`"Kentucky"`,`states`,`"Kentucky"`},
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