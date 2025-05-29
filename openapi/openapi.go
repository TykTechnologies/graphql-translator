package openapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/TykTechnologies/graphql-go-tools/pkg/ast"
	"github.com/TykTechnologies/graphql-go-tools/pkg/introspection"
	"github.com/TykTechnologies/graphql-go-tools/pkg/operationreport"
	"github.com/getkin/kin-openapi/openapi3"
)

var (
	errTypeNameExtractionImpossible = errors.New("type name extraction is impossible")
	errNotPrimitiveType             = errors.New("not a primitive type")
)

type converter struct {
	openapi         *openapi3.T
	knownFullTypes  map[string]*knownFullTypeDetails
	knownEnums      map[string]introspection.FullType
	knownUnions     map[string]introspection.FullType
	fullTypes       []introspection.FullType
	currentPathName string
	currentPathItem *openapi3.PathItem
}

type knownFullTypeDetails struct {
	hasDescription bool
}

func newConverter(document *openapi3.T) *converter {
	return &converter{
		openapi:        document,
		knownFullTypes: make(map[string]*knownFullTypeDetails),
		knownEnums:     make(map[string]introspection.FullType),
		knownUnions:    make(map[string]introspection.FullType),
		fullTypes:      make([]introspection.FullType, 0),
	}
}

// ImportParsedOpenAPIv3Document imports a parsed OpenAPI v3 document and converts it into a GraphQL AST Document.
func ImportParsedOpenAPIv3Document(document *openapi3.T, report *operationreport.Report) *ast.Document {
	c := newConverter(document)
	data := introspection.Data{}

	queryType, err := c.importQueryType()
	if err != nil {
		report.AddInternalError(err)
		return nil
	}
	if len(queryType.Fields) > 0 {
		data.Schema.QueryType = &introspection.TypeName{
			Name: "Query",
		}
		data.Schema.Types = append(data.Schema.Types, *queryType)
	} else {
		data.Schema.QueryType = &introspection.TypeName{
			Name: "QueryPlaceholder",
		}
		stringScalarType := "String"
		placeholderObject := introspection.FullType{
			Kind:        introspection.OBJECT,
			Name:        "QueryPlaceholder",
			Description: "Placeholder object",
			Fields: []introspection.Field{
				{
					Name:        "message",
					Description: "Placeholder field",
					Type:        introspection.TypeRef{Kind: introspection.SCALAR, Name: &stringScalarType},
				},
			},
		}
		data.Schema.Types = append(data.Schema.Types, placeholderObject)
	}

	mutationType, err := c.importMutationType()
	if err != nil {
		report.AddInternalError(err)
		return nil
	}
	if len(mutationType.Fields) > 0 {
		data.Schema.MutationType = &introspection.TypeName{
			Name: "Mutation",
		}
		data.Schema.Types = append(data.Schema.Types, *mutationType)
	}

	fullTypes, err := c.importFullTypes()
	if err != nil {
		report.AddInternalError(err)
		return nil
	}
	data.Schema.Types = append(data.Schema.Types, fullTypes...)

	outputPretty, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		report.AddInternalError(err)
		return nil
	}

	jc := introspection.JsonConverter{}
	buf := bytes.NewBuffer(outputPretty)
	doc, err := jc.GraphQLDocument(buf)
	if err != nil {
		report.AddInternalError(err)
		return nil
	}
	return doc
}

// ParseOpenAPIDocument parses an OpenAPI document from a byte slice input.
// It creates a new loader with IsExternalRefsAllowed set to true.
// The loader then loads the document from the input data, returning any errors encountered.
// If the document is successfully loaded, it is validated using the loader's context.
// If validation fails, the validation error is returned.
// Otherwise, the parsed document is returned along with a nil error.
func ParseOpenAPIDocument(input []byte) (*openapi3.T, error) {
	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	document, err := loader.LoadFromData(input)
	if err != nil {
		return nil, err
	}

	// Validate the document
	if err = document.Validate(loader.Context); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	return document, nil
}

// ImportOpenAPIDocumentByte imports an OpenAPI document from a byte slice input.
// It first parses the input using the ParseOpenAPIDocument function.
// If there is an error during parsing, it adds the error to the operationreport.Report
// and returns nil with the report.
// Otherwise, it passes the parsed document to the ImportParsedOpenAPIv3Document function
// along with the operationreport.Report, and returns the imported OpenAPI document
// as ast.Document and the report.
func ImportOpenAPIDocumentByte(input []byte) (*ast.Document, operationreport.Report) {
	report := operationreport.Report{}
	document, err := ParseOpenAPIDocument(input)
	if err != nil {
		report.AddInternalError(err)
		return nil, report
	}
	return ImportParsedOpenAPIv3Document(document, &report), report
}

// ImportOpenAPIDocumentString imports an OpenAPI document from a string input.
// It delegates to ImportOpenAPIDocumentByte, passing the input as a byte slice.
// The function returns the imported OpenAPI document as ast.Document and an operationreport.Report.
// If there are any errors during the import process, they will be included in the report.
func ImportOpenAPIDocumentString(input string) (*ast.Document, operationreport.Report) {
	return ImportOpenAPIDocumentByte([]byte(input))
}
