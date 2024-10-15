package main

import (
	"fmt"
	"log"
	"os"

	"github.com/getkin/kin-openapi/openapi3"
	"gopkg.in/yaml.v3"
)

const (
	BasePath             = "../spec/v2/"
	ServiceOpenapiFolder = BasePath + "service-openapi/"
	IotAPIOpenApi        = ServiceOpenapiFolder + "iotapi-openapi.yaml"
	GroupsViewsOpenApi   = ServiceOpenapiFolder + "groups-views-openapi.yaml"
	GroupsFoldersOpenApi = ServiceOpenapiFolder + "groups-folders-openapi.yaml"

	MergedOpenApi = BasePath + "swagger.yaml"
)

func main() {

	// Load the main OpenAPI specification
	loader := openapi3.NewLoader()
	docIOT, err := loader.LoadFromFile(IotAPIOpenApi)
	if err != nil {
		log.Fatalf("Failed to load %s: %v", IotAPIOpenApi, err)
	}

	apisToMerge := []string{GroupsViewsOpenApi, GroupsFoldersOpenApi}

	mergedPaths := []openapi3.NewPathsOption{}
	for _, api := range apisToMerge {
		// Load the second OpenAPI specification
		docVIEWS, err := loader.LoadFromFile(api)
		if err != nil {
			log.Fatalf("Failed to load %s: %v", api, err)
		}

		// Merge the paths
		for _, path := range docIOT.Paths.InMatchingOrder() {
			pathItem := docIOT.Paths.Find(path)
			mergedPaths = append(mergedPaths, openapi3.WithPath(fmt.Sprintf("/iot/%s", path), pathItem))
		}
		for _, path := range docVIEWS.Paths.InMatchingOrder() {
			pathItem := docVIEWS.Paths.Find(path)
			mergedPaths = append(mergedPaths, openapi3.WithPath(path, pathItem))
		}

		// Merge the components
		if docIOT.Components.Schemas == nil {
			docIOT.Components.Schemas = make(map[string]*openapi3.SchemaRef)
		}
		for name, schema := range docVIEWS.Components.Schemas {
			docIOT.Components.Schemas[name] = schema
		}

		if docIOT.Components.Responses == nil {
			docIOT.Components.Responses = make(map[string]*openapi3.ResponseRef)
		}
		for name, response := range docVIEWS.Components.Responses {
			docIOT.Components.Responses[name] = response
		}

		if docIOT.Components.Parameters == nil {
			docIOT.Components.Parameters = make(map[string]*openapi3.ParameterRef)
		}
		for name, parameter := range docVIEWS.Components.Parameters {
			docIOT.Components.Parameters[name] = parameter
		}

		if docIOT.Components.RequestBodies == nil {
			docIOT.Components.RequestBodies = make(map[string]*openapi3.RequestBodyRef)
		}
		for name, requestBody := range docVIEWS.Components.RequestBodies {
			docIOT.Components.RequestBodies[name] = requestBody
		}

		if docIOT.Components.Headers == nil {
			docIOT.Components.Headers = make(map[string]*openapi3.HeaderRef)
		}
		for name, header := range docVIEWS.Components.Headers {
			docIOT.Components.Headers[name] = header
		}

		if docIOT.Components.SecuritySchemes == nil {
			docIOT.Components.SecuritySchemes = make(map[string]*openapi3.SecuritySchemeRef)
		}
		for name, securityScheme := range docVIEWS.Components.SecuritySchemes {
			docIOT.Components.SecuritySchemes[name] = securityScheme
		}
	}

	// Add merged
	docIOT.Paths = openapi3.NewPaths(mergedPaths...)

	ym, err := docIOT.MarshalYAML()
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

	log.Println("Merged OpenAPI specification saved to " + MergedOpenApi)

}
