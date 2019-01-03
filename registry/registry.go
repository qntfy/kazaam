package registry

import "errors"

type Converter interface {
	Init(config []byte) (error)
	Convert(jsonData []byte, value []byte, args []byte) (newValue []byte, err error)
}

var converters map[string]Converter

func init() {
	converters = make(map[string]Converter)
}

func RegisterConverter(name string, conv Converter) error {
	_, ok := converters[name]
	if ok {
		return errors.New("converter with that name already registered")
	}
	converters[name] = conv
	return nil
}

func GetConverter(name string) Converter {
	conv, ok := converters[name]
	if ok {
		return conv
	}
	return nil
}
