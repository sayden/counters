package counters

import (
	"bytes"
	"io"
	"os"

	"github.com/invopop/jsonschema"
	"github.com/pkg/errors"
	"github.com/xeipuuv/gojsonschema"
)

func ValidateSchemaReader[S CounterTemplate | CardsTemplate](r io.Reader) error {
	byt, err := io.ReadAll(r)
	if err != nil {
		return errors.Wrap(err, "could not read JSON file")
	}
	return ValidateSchemaBytes[S](byt)
}

func ValidateSchemaBytes[S CounterTemplate | CardsTemplate](docByt []byte) (err error) {
	if bytes.Contains(docByt, []byte("\n")) {
		docByt = bytes.ReplaceAll(docByt, []byte("\n"), []byte(""))
	}

	reflector := new(jsonschema.Reflector)
	counterTemplateSchemaMarshaller := reflector.Reflect(new(S))
	schemaByt, err := counterTemplateSchemaMarshaller.MarshalJSON()
	if err != nil {
		return errors.Wrap(err, "could not marshal counter template schema")
	}
	schema := gojsonschema.NewBytesLoader(schemaByt)

	documentLoader := gojsonschema.NewBytesLoader(docByt)
	result, err := gojsonschema.Validate(schema, documentLoader)
	if err != nil {
		return errors.Wrap(err, "could not validate JSON file")
	}

	return validateResult(result)
}

func ValidateSchemaAtPath[S CounterTemplate | CardsTemplate](inputPath string) error {
	byt, err := os.ReadFile(inputPath)
	if err != nil {
		return errors.Wrap(err, "could not read JSON file")
	}

	return ValidateSchemaBytes[S](byt)
}

func validateResult(result *gojsonschema.Result) error {
	if !result.Valid() {
		err := errors.New("JSON file is not valid")
		for _, desc := range result.Errors() {
			err = errors.Wrap(err, "\n"+desc.String())
		}

		return err
	}

	return nil
}
