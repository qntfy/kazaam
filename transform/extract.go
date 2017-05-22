package transform

// ExtractRaw returns the specified path as the top-level object in raw []byte.
func ExtractRaw(spec *Config, data []byte) ([]byte, error) {
	outPath, ok := (*spec.Spec)["path"]
	if !ok {
		return nil, SpecError("Unable to get path")
	}
	result, err := getJSONRaw(data, outPath.(string), spec.Require)
	if err != nil {
		return nil, err
	}
	return result, nil
}
