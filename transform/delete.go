package transform

import (
        "fmt"
)

// Delete deletes keys in-place from the provided data if they exist
// keys are specified in an array under "keys" in the spec.
func Delete(spec *Config, data []byte) ([]byte, error) {
	pathList, pathsOk := (*spec.Spec)["paths"]
	if !pathsOk {
		return nil, SpecError("Unable to get paths to delete")
	}
	for _, pItem := range pathList.([]interface{}) {
		path, ok := pItem.(string)
		if !ok {
			return nil, SpecError(fmt.Sprintf("Error processing %v: path should be a string", pItem))
		}

		var err error
		data, err = delJSONRaw(data, path, spec.Require)
		if err != nil {
			return nil, err
		}
	}
	return data, nil
}
