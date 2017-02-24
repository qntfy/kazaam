package kazaam

import simplejson "github.com/bitly/go-simplejson"

// no op transform -- useful for testing/default behavior
func transformPass(spec *spec, data *simplejson.Json) (*simplejson.Json, error) {
	return data, nil
}
