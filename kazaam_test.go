package kazaam

import (
	"testing"

	simplejson "github.com/bitly/go-simplejson"
)

func TestGetUnknownTransform(t *testing.T) {
	testJSON := `{"test":"data"}`
	tformName := "doesnt-exist"
	tform := getTransform(spec{Operation: &tformName})
	dataIn, _ := simplejson.NewJson([]byte(testJSON))
	dataOut, err := tform(&spec{}, dataIn)
	if err != nil {
		t.Error("Unexpected error: ", err)
	}
	jsonOut, _ := dataOut.MarshalJSON()
	if string(jsonOut) != testJSON {
		t.Error("Unknown transform type handled incorrectly")
	}
}
