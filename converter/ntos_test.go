package converter

import (
	"github.com/qntfy/kazaam/registry"
	"strconv"
	"testing"
)

func TestNtos_Convert(t *testing.T) {

	registry.RegisterConverter("ntos", &Ntos{})
	c := registry.GetConverter("ntos")

	table := []struct {
		value     string
		arguments string
		expected  string
	}{
		{`5`, ``, `"5"`,},
		{`5.01`, ``, `"5.01"`,},
		{`-5.01`, ``, `"-5.01"`,},
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