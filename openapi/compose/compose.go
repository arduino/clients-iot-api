package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"slices"

	"github.com/getkin/kin-openapi/openapi3"
	"gopkg.in/yaml.v3"

	_ "embed"
)

const (
	BasePath                = "../"
	OpenapiComponentsFolder = BasePath + "components/"
	IotApiOpenApi           = OpenapiComponentsFolder + "iotapi-openapi.yaml"
	HeaderOpenApi           = OpenapiComponentsFolder + "header.yaml"
	FooterOpenApi           = OpenapiComponentsFolder + "footer.yaml"

	MergedOpenYAMLApi = BasePath + "openapi.yaml"
	MergedOpenJSONApi = BasePath + "openapi.json"
)

//go:embed path-blacklist.json
var pathBlacklist string

type blacklist struct {
	Blacklist []string `json:"blacklist"`
}

func main() {
	// Load blacklist
	var bl blacklist
	_ = json.Unmarshal([]byte(pathBlacklist), &bl)

	log.Println("Compose OpenAPI...")

	// Load the main (IOT-API) OpenAPI specification
	loader := openapi3.NewLoader()
	masterOpenapi, err := loader.LoadFromFile(IotApiOpenApi)
	if err != nil {
		log.Fatalf("Failed to load %s: %v", IotApiOpenApi, err)
	}

	mergedPaths := []openapi3.NewPathsOption{}
	for _, path := range masterOpenapi.Paths.InMatchingOrder() {
		if slices.Contains(bl.Blacklist, path) {
			continue
		}
		pathItem := masterOpenapi.Paths.Find(path)
		mergedPaths = append(mergedPaths, openapi3.WithPath(fmt.Sprintf("/iot%s", path), pathItem))
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

	// Add merged paths to the master document
	masterOpenapi.Paths = openapi3.NewPaths(mergedPaths...)

	ym, err := masterOpenapi.MarshalYAML()
	if err != nil {
		log.Fatalf("Failed to marshal merged document to YAML: %v", err)
	}
	yamlData, err := yaml.Marshal(ym)
	if err != nil {
		log.Fatalf("Failed to marshal merged document to YAML: %v", err)
	}

	err = os.WriteFile(MergedOpenYAMLApi, yamlData, 0644)
	if err != nil {
		log.Fatalf("Failed to write %s: %v", MergedOpenYAMLApi, err)
	}

	log.Println("Merged OpenAPI saved to " + MergedOpenYAMLApi)

	js, err := masterOpenapi.MarshalJSON()
	if err != nil {
		log.Fatalf("Failed to marshal merged document to JSON: %v", err)
	}

	jsonData, err := json.Marshal(js)
	if err != nil {
		log.Fatalf("Failed to marshal merged document to JSON: %v", err)
	}

	err = os.WriteFile(MergedOpenJSONApi, jsonData, 0644)
	if err != nil {
		log.Fatalf("Failed to write %s: %v", MergedOpenJSONApi, err)
	}

	log.Println("Merged OpenAPI saved to " + MergedOpenJSONApi)
}
