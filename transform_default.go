package kazaam

import (
	"strings"

	simplejson "github.com/bitly/go-simplejson"
)

// simple transform to set default values in output json
func transformDefault(spec *spec, data *simplejson.Json) (*simplejson.Json, error) {
	for k, v := range *spec.Spec {
		data.SetPath(strings.Split(k, "."), v)
	}
	return data, nil
}
