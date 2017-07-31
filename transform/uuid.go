package transform

import (
	"strings"

	"github.com/buger/jsonparser"
	uuid "github.com/satori/go.uuid"
)

var (
	versionError = SpecError("Please set version 3 || 4 || 5")
)

// UUID tries to generate a UUID based on spec components
func UUID(spec *Config, data []byte) ([]byte, error) {

	// iterate through the spec
	for k, v := range *spec.Spec {
		outPath := strings.Split(k, ".")

		// convert to corrct type
		uuidSpec, ok := v.(map[string]interface{})
		if !ok {
			return nil, SpecError("Invalid Spec for UUID")
		}

		//grab version
		version, ok := uuidSpec["version"]

		if !ok {
			return nil, versionError
		}

		var u uuid.UUID
		var err error

		switch version {
		case 4.0:
			u = uuid.NewV4()

		case 3.0, 5.0:

			names, ok := uuidSpec["names"]
			if !ok {
				return nil, SpecError("Must provide names field")
			}

			nameSpace, ok := uuidSpec["nameSpace"].(string)
			if !ok {
				return nil, SpecError("Must provide namesapce, Must be a string")
			}

			var nameSpaceUUID uuid.UUID

			// swtich on the namespace
			switch nameSpace {
			case "DNS":
				nameSpaceUUID = uuid.NamespaceDNS
			case "URL":
				nameSpaceUUID = uuid.NamespaceURL
			case "OID":
				nameSpaceUUID = uuid.NamespaceOID
			case "X500":
				nameSpaceUUID = uuid.NamespaceX500
			default:
				nameSpaceUUID, err = uuid.FromString(nameSpace)
				if err != nil {
					return nil, SpecError("nameSpace is not a valid UUID or is not DNS, URL, OID, X500")
				}
			}

			nameFields, ok := names.([]interface{})
			if !ok {
				return nil, SpecError("Spec is invalid")
			}

			// loop over the names field
			for _, field := range nameFields {
				p, _ := field.(map[string]interface{})["path"].(string)

				name, err := getJSONRaw(data, p, false)
				if err == jsonparser.KeyPathNotFoundError {

					d, ok := field.(map[string]interface{})["default"].(string)
					if !ok {
						return nil, SpecError("Spec is invalid")

					}
					name = []byte(d)

				}

				// check if there is an empty uuid & version 3
				if u.String() == "00000000-0000-0000-0000-000000000000" && version == 3.0 {

					u = uuid.NewV3(nameSpaceUUID, string(name))

					// same as above except version 5
				} else if u.String() == "00000000-0000-0000-0000-000000000000" && version == 5.0 {

					u = uuid.NewV5(nameSpaceUUID, string(name))

				} else if version == 3.0 {

					u = uuid.NewV3(u, string(name))

				} else if version == 5.0 {

					u = uuid.NewV3(u, string(name))
				}
			}

		default:
			return nil, versionError

		}
		d, err := jsonparser.Set(data, bookend([]byte(u.String()), '"', '"'), outPath...)
		if err != nil {
			return nil, err
		}

		return d, nil
	}
	return nil, SpecError("Spec invalid for UUID")
}
