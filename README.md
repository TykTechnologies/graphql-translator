# graphql-translator

graphql-translator is a library that takes OpenAPI and AsyncAPI documents and translates them into GraphQL documents.

This library uses [graphql-go-tools](https://github.com/TykTechnologies/graphql-go-tools) library for GraphQL handling.

## OpenAPI

Supported OpenAPI versions:

* 3.0.0

### How to translate OpenAPI to GraphQL

```go
doc, report := ImportOpenAPIDocumentString(openapiDocument)
if report.HasErrors() {
	return report
}

// Now, we can pretty-print the resulting abstract syntax tree.
w := &bytes.Buffer{}
err = astprinter.PrintIndent(doc, nil, []byte("  "), w)
if err != nil {
	return err
}

// This will print the GraphQL document.
fmt.Println(w.String())
```

`openapi` package also provides `ImportOpenAPIDocumentByte`, `ParseOpenAPIDocument` and `ImportParsedOpenAPIv3Document` 
methods. You can check the package documents to see how they work.

## AsyncAPI

Supported AsyncAPI versions:

* 2.0.0
* 2.1.0
* 2.2.0
* 2.3.0
* 2.4.0

### How to translate AsyncAPI to GraphQL

```go
doc, report := ImportAsyncAPIDocumentString(string(asyncapiDoc))
if report.HasErrors() {
	return err
}

// Now, we can pretty-print the resulting abstract syntax tree.
w := &bytes.Buffer{}
err = astprinter.PrintIndent(doc, nil, []byte("  "), w)
if err != nil {
	return err
}

// This will print the GraphQL document.
fmt.Println(w.String())
```

`asyncapi` package also provides `ImportAsyncAPIDocumentByte`, `ParseAsyncAPIDocument` and `ImportParsedAsyncAPIDocument`
methods. You can check the package documents to see how they work.

## License

MIT License - see LICENSE for more details.

## Contributing

Feel free to file an issue in case of bugs. We're open to your ideas to enhance the repository.

You are open to contribute via PR's. Please open an issue to discuss your idea before implementing it, 
so we can have a discussion. Make sure to comply with the linting rules. You must not add untested code.
