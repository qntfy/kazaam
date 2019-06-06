package transform

import (
	"encoding/json"
)

// Merge joins multiple array values into array of objects with property names containing their matching array. Arrays
// should be the same length.
func Merge(spec *Config, data []byte) ([]byte, error) {
	var outData []byte

	if spec.InPlace {
		outData = data
	} else {
		outData = []byte(`{}`)
	}

	// iterate through the spec
	for k, v := range *spec.Spec {

		// map[prop_name] = [ values... ]
		arrayVals := make(map[string][]interface{})
		outVals := make([]map[string]interface{}, 0)

		mergeSpec, ok := v.([]interface{})
		if !ok {
			return nil, SpecError("Invalid Spec for Merge")
		}

		l := 0

		for i, v := range mergeSpec {
			arraySpec := v.(map[string]interface{})
			var name, array string
			name, ok = arraySpec["name"].(string)
			if !ok {
				return nil, SpecError("Array spec missing name for Merge")
			}
			array, ok = arraySpec["array"].(string)
			if !ok {
				return nil, SpecError("Array spec missing array for Merge")
			}

			var dataForV []byte
			var err error

			dataForV, err = getJSONRaw(data, array, true)
			if err != nil {
				return nil, err
			}

			var arrayValues []interface{}
			err = json.Unmarshal(dataForV, &arrayValues)

			arrayVals[name] = arrayValues
			if i == 0 {
				l = len(arrayValues)
			} else if l != len(arrayValues) {
				return nil, SpecError("Arrays must be the same length for Merge")
			}
		}

		for i := 0; i < l; i++ {
			m := make(map[string]interface{})
			for k, v := range arrayVals {
				m[k] = v[0]
				arrayVals[k] = v[1:]
			}
			outVals = append(outVals, m)
		}

		dataForV, err := json.Marshal(outVals)
		if err != nil {
			return nil, err
		}

		outData, err = setJSONRaw(outData, dataForV, k)
		if err != nil {
			return nil, err
		}

	}

	return outData, nil

}
