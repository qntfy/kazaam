package kazaam

import (
	"encoding/json"
	"testing"
)

func TestSpecUnmarshalFailure(t *testing.T) {
	testSpec := `{"operation": "unimplemented", "spec": {}}`
	var s spec

	err := json.Unmarshal([]byte(testSpec), &s)

	if err == nil {
		t.Error("Should have returned an error unmarshaling spec")
	}
}
