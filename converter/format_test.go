package converter

import (
	"github.com/qntfy/kazaam/registry"
	"strconv"
	"testing"
)

func TestFormat_Convert(t *testing.T) {
	registry.RegisterConverter("format", &Format{})
	c := registry.GetConverter("format")

	table := []struct {
		value     string
		arguments string
		expected  string
	}{
		{`5.1`, `%.0f`, `"5"`,},
		{`5.23546`, `%.2f`, `"5.24"`,},
		{`0.01`, `%.4f`, `"0.0100"`},
		{`true`, `%t`, `"true"`},
		{`"the something fox"`, `%s jumped over something.`, `"the something fox jumped over something."`},
		{`42`, `%d`,`"42"`},
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
