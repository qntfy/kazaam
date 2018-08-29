package transform

import (
	"strings"

	uuid "github.com/gofrs/uuid"
)

var (
	versionError = SpecError("Please set version 3 || 4 || 5")
)

// UUID tries to generate a UUID based on spec components
func UUID(spec *Config, data []byte) ([]byte, error) {

	// iterate through the spec
	for k, v := range *spec.Spec {
		// convert spec to correct type
		uuidSpec, ok := v.(map[string]interface{})
		if !ok {
			return nil, SpecError("Invalid Spec for UUID")
		}
		version := getUUIDVersion(uuidSpec)
		if version < 3 || version > 5 {
			return nil, versionError
		}

		var u uuid.UUID
		var err error

		switch version {
		case 4:
			u, err = uuid.NewV4()
			if err != nil {
				return nil, err
			}

		case 3, 5:
			// choose the correct UUID function
			var NewUUID func(uuid.UUID, string) uuid.UUID
			NewUUID = uuid.NewV3
			if version == 5 {
				NewUUID = uuid.NewV5
			}

			// pull required configuration from spec and do validation
			names, ok := uuidSpec["names"]
			if !ok {
				return nil, SpecError("Must provide names field")
			}
			namespaceString, ok := uuidSpec["namespace"].(string)
			if !ok {
				return nil, SpecError("Must provide `namespace` as a string")
			}
			nameFields, ok := names.([]interface{})
			if !ok {
				return nil, SpecError("Spec is invalid. `Names` field must be an array.")
			}

			// generate the required namespace
			u, err = namespaceFromString(namespaceString)
			if err != nil {
				return nil, SpecError("Namespace is not a valid UUID or is not DNS, URL, OID, X500")
			}

			// loop over the names field
			for _, field := range nameFields {
				p, _ := field.(map[string]interface{})["path"].(string)

				name, pathErr := getJSONRaw(data, p, true)
				// if a string, remove the heading and trailing quote
				nameString := strings.TrimPrefix(strings.TrimSuffix(string(name), "\""), "\"")
				if pathErr == NonExistentPath {
					nameString, ok = field.(map[string]interface{})["default"].(string)
					if !ok {
						return nil, SpecError("Spec is invalid. Unable to get path or default")
					}
				}
				u = NewUUID(u, nameString)
			}

		default:
			return nil, versionError

		}
		// set the uuid in the appropriate place
		data, err = setJSONRaw(data, bookend([]byte(u.String()), '"', '"'), k)
		if err != nil {
			return nil, err
		}
	}
	return data, nil
}

func namespaceFromString(namespace string) (uuid.UUID, error) {
	var u uuid.UUID
	var err error
	switch namespace {
	case "DNS":
		u = uuid.NamespaceDNS
	case "URL":
		u = uuid.NamespaceURL
	case "OID":
		u = uuid.NamespaceOID
	case "X500":
		u = uuid.NamespaceX500
	default:
		u, err = uuid.FromString(namespace)
	}
	return u, err
}

func getUUIDVersion(uuidSpec map[string]interface{}) int {
	var version int
	versionInterface, ok := uuidSpec["version"]
	if !ok {
		return -1
	}
	versionFloat, ok := versionInterface.(float64)
	version = int(versionFloat)
	if !ok || version < 3 || version > 5 {
		return -2
	}
	return version
}
