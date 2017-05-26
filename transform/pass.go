package transform

// Pass performs no manipulation of the passed-in data. It is useful
// for testing/default behavior.
func Pass(spec *Config, data []byte) ([]byte, error) {
	return data, nil
}
