package main

import (
	"log"
	"os"

	"github.com/getkin/kin-openapi/openapi3"
	"gopkg.in/yaml.v3"

	_ "embed"
)

const (
	BasePath                = "../"
	OpenapiComponentsFolder = BasePath + "components/"
	RawOpenApi              = OpenapiComponentsFolder + "raw-openapi.yaml"
	HeaderOpenApi           = OpenapiComponentsFolder + "header.yaml"
	FooterOpenApi           = OpenapiComponentsFolder + "footer.yaml"

	MergedOpenApi = BasePath + "openapi.yaml"
)

func main() {

	log.Println("Compose OpenAPI...")

	// Load the main OpenAPI specification
	loader := openapi3.NewLoader()
	masterOpenapi, err := loader.LoadFromFile(RawOpenApi)
	if err != nil {
		log.Fatalf("Failed to load %s: %v", RawOpenApi, err)
	}

	if masterOpenapi.Components.Schemas == nil {
		masterOpenapi.Components.Schemas = make(map[string]*openapi3.SchemaRef)
	}
	if masterOpenapi.Components.Responses == nil {
		masterOpenapi.Components.Responses = make(map[string]*openapi3.ResponseRef)
	}
	if masterOpenapi.Components.Parameters == nil {
		masterOpenapi.Components.Parameters = make(map[string]*openapi3.ParameterRef)
	}
	if masterOpenapi.Components.RequestBodies == nil {
		masterOpenapi.Components.RequestBodies = make(map[string]*openapi3.RequestBodyRef)
	}
	if masterOpenapi.Components.Headers == nil {
		masterOpenapi.Components.Headers = make(map[string]*openapi3.HeaderRef)
	}
	if masterOpenapi.Components.SecuritySchemes == nil {
		masterOpenapi.Components.SecuritySchemes = make(map[string]*openapi3.SecuritySchemeRef)
	}

	// Merge the components

	// Load header and footer
	headerOpenapi, err := loader.LoadFromFile(HeaderOpenApi)
	if err != nil {
		log.Fatalf("Failed to load %s: %v", HeaderOpenApi, err)
	}

	if headerOpenapi.Servers != nil {
		masterOpenapi.Servers = append(masterOpenapi.Servers, headerOpenapi.Servers...)
	}

	footerOpenapi, err := loader.LoadFromFile(FooterOpenApi)
	if err != nil {
		log.Fatalf("Failed to load %s: %v", FooterOpenApi, err)
	}

	if footerOpenapi.Components.SecuritySchemes != nil {
		for name, securityScheme := range footerOpenapi.Components.SecuritySchemes {
			masterOpenapi.Components.SecuritySchemes[name] = securityScheme
		}
	}

	ym, err := masterOpenapi.MarshalYAML()
	if err != nil {
		log.Fatalf("Failed to marshal merged document to YAML: %v", err)
	}
	yamlData, err := yaml.Marshal(ym)
	if err != nil {
		log.Fatalf("Failed to marshal merged document to YAML: %v", err)
	}

	err = os.WriteFile(MergedOpenApi, yamlData, 0644)
	if err != nil {
		log.Fatalf("Failed to write %s: %v", MergedOpenApi, err)
	}

	log.Println("Merged OpenAPI saved to " + MergedOpenApi)

}
