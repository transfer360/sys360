package validate

import (
	"fmt"
	"github.com/xeipuuv/gojsonschema"
)

type ValidationResult struct {
	Valid  bool
	Errors []string
}

func IsValidAgainstSchema(schemaFile string, object []byte) (sv ValidationResult, err error) {

	schemaLoader := gojsonschema.NewReferenceLoader(fmt.Sprintf("file://%s", schemaFile))
	documentLoader := gojsonschema.NewStringLoader(string(object))

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return sv, fmt.Errorf("there was an error validating document: %w", err)
	} else {

		if result.Valid() {

			sv.Valid = true
			return sv, nil

		} else {

			sv.Errors = make([]string, 0)
			for _, desc := range result.Errors() {
				sv.Errors = append(sv.Errors, desc.String())
			}
			return sv, nil
		}
	}

}
