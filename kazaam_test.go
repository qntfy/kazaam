package kazaam

import (
	"fmt"
	"testing"

	"github.com/qntfy/jsonparser"
	"github.com/qntfy/kazaam/transform"
)

func TestDefaultKazaamGetUnknownTransform(t *testing.T) {
	_, err := NewKazaam(`[{"operation": "doesnt-exist"}]`)

	if err == nil {
		t.Error("Should have thrown error for unknown transform")
	}
}

func TestKazaamWithRegisteredTransform(t *testing.T) {
	kc := NewDefaultConfig()
	kc.RegisterTransform("3rd-party", func(spec *transform.Config, data []byte) ([]byte, error) {
		data, _ = jsonparser.Set(data, []byte("doesnt-exist"), "does-exist")
		return data, nil
	})
	_, err := New(`[{"operation": "3rd-party"}]`, kc)
	if err != nil {
		t.Error("Shouldn't have thrown error for registered 3rd-party transform")
	}
}

func TestReregisterKazaamTransform(t *testing.T) {
	kc := NewDefaultConfig()
	err := kc.RegisterTransform("shift", nil)

	if err == nil {
		t.Error("Should have thrown error for duplicated transform name")
	}
}

func TestDefaultTransformsSetCardinarily(t *testing.T) {
	if len(validSpecTypes) != 9 {
		t.Error("Unexpected number of default transforms. Missing tests?")
	}
}

func TestErrorTypes(t *testing.T) {
	testCases := []struct {
		typ         int
		msg         string
		expectedMsg string
	}{
		{ParseError, "test1", "ParseError - test1"},
		{RequireError, "test2", "RequiredError - test2"},
		{SpecError, "test3", "SpecError - test3"},
		{5, "test3", "SpecError - test3"},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%d error type", tc.typ), func(t *testing.T) {
			e := Error{ErrType: tc.typ, ErrMsg: tc.msg}
			if e.Error() != tc.expectedMsg {
				t.Errorf("got %s; want %s", e.Error(), tc.expectedMsg)
			}
		})
	}
}

func ExampleNewKazaam() {
	k, _ := NewKazaam(`[{"operation": "shift", "spec": {"output": "input"}}]`)
	kazaamOut, _ := k.TransformJSONStringToString(`{"input":"input value"}`)

	fmt.Println(kazaamOut)
	// Output:
	// {"output":"input value"}
}

func ExampleNew() {
	// Initialize a default Kazaam instance (i.e. same as NewKazaam(spec string))
	k, _ := New(`[{"operation": "shift", "spec": {"output": "input"}}]`, NewDefaultConfig())
	kazaamOut, _ := k.TransformJSONStringToString(`{"input":"input value"}`)

	fmt.Println(kazaamOut)
	// Output:
	// {"output":"input value"}
}

func ExampleConfig_RegisterTransform() {
	// use the default config to have access to built-in kazaam transforms
	kc := NewDefaultConfig()

	// register the new custom transform called "copy" which supports copying the
	// value of a top-level key to another top-level key
	kc.RegisterTransform("copy", func(spec *transform.Config, data []byte) ([]byte, error) {
		// the internal `Spec` will contain a mapping of source and target keys
		for targetField, sourceFieldInt := range *spec.Spec {
			sourceField := sourceFieldInt.(string)
			// Note: jsonparser.Get() strips quotes from returned strings, so a real
			// transform would need handling for that. We use a Number below for simplicity.
			result, _, _, _ := jsonparser.Get(data, sourceField)
			data, _ = jsonparser.Set(data, result, targetField)
		}
		return data, nil
	})

	k, _ := New(`[{"operation": "copy", "spec": {"output": "input"}}]`, kc)
	kazaamOut, _ := k.TransformJSONStringToString(`{"input":72}`)

	fmt.Println(kazaamOut)
	// Output:
	// {"input":72,"output":72}
}
