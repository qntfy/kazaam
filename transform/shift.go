package transform

import (
	"fmt"
	"strings"

	"github.com/buger/jsonparser"
)

// Shift moves values from one provided json path to another in raw []byte.
func Shift(spec *Config, data []byte) ([]byte, error) {
	outData := []byte(`{}`)
	for k, v := range *spec.Spec {
		array := true
		outPath := strings.Split(k, ".")

		var keyList []string

		// check if `v` is a string or list and build a list of keys to evaluate
		switch v.(type) {
		case string:
			keyList = append(keyList, v.(string))
			array = false
		case []interface{}:
			for _, vItem := range v.([]interface{}) {
				vItemStr, found := vItem.(string)
				if !found {
					return nil, ParseError(fmt.Sprintf("Warn: Unable to coerce element to json string: %v", vItem))
				}
				keyList = append(keyList, vItemStr)
			}
		default:
			return nil, ParseError(fmt.Sprintf("Warn: Unknown type in message for key: %s", k))
		}

		// iterate over keys to evaluate
		// Note: this could be sped up significantly (especially for large shift transforms)
		// by using `jsonparser.EachKey()` to iterate through data once and pick up all the
		// needed underlying data. It would be a non-trivial update since you'd have to make
		// recursive calls and keep track of all the key paths at each level.
		// Currently we iterate at worst once per key in spec, with a better design it would be once
		// per spec.
		for _, v := range keyList {
			var dataForV []byte
			var err error

			// grab the data
			if v == "$" {
				dataForV = data
			} else {
				dataForV, err = getJSONRaw(data, v, spec.Require)
				if err != nil {
					return nil, err
				}
			}

			// if array flag set, encapsulate data
			if array {
				// bookend() is destructive to underlying slice, need to copy.
				// extra capacity saves an allocation and copy during bookend.
				tmp := make([]byte, len(dataForV), len(dataForV)+2)
				copy(tmp, dataForV)
				dataForV = bookend(tmp, '[', ']')
			}
			// Note: following pattern from current Shift() - if multiple elements are included in an array,
			// they will each successively overwrite each other and only the last element will be included
			// in the transformed data.
			outData, err = jsonparser.Set(outData, dataForV, outPath...)
			if err != nil {
				return nil, err
			}
		}
	}
	return outData, nil
}
