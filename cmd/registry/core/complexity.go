// Copyright 2020 Google LLC. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package core

import (
	discovery "github.com/google/gnostic/discovery"
	metrics "github.com/google/gnostic/metrics"
	openapi_v2 "github.com/google/gnostic/openapiv2"
	openapi_v3 "github.com/google/gnostic/openapiv3"
)

func SummarizeOpenAPIv2Document(document *openapi_v2.Document) *metrics.Complexity {
	summary := &metrics.Complexity{}
	if document.Definitions != nil && document.Definitions.AdditionalProperties != nil {
		for _, pair := range document.Definitions.AdditionalProperties {
			summarizeOpenAPIv2Schema(summary, pair.Value)
		}
	}
	for _, pair := range document.Paths.Path {
		summary.PathCount++
		v := pair.Value
		if v.Get != nil {
			summary.GetCount++
		}
		if v.Post != nil {
			summary.PostCount++
		}
		if v.Put != nil {
			summary.PutCount++
		}
		if v.Delete != nil {
			summary.DeleteCount++
		}
	}
	return summary
}

func summarizeOpenAPIv2Schema(summary *metrics.Complexity, schema *openapi_v2.Schema) {
	summary.SchemaCount++
	if schema.Properties != nil {
		for _, pair := range schema.Properties.AdditionalProperties {
			summary.SchemaPropertyCount++
			summarizeOpenAPIv2Schema(summary, pair.Value)
		}
	}
}

func SummarizeOpenAPIv3Document(document *openapi_v3.Document) *metrics.Complexity {
	summary := &metrics.Complexity{}
	if document.Components != nil && document.Components.Schemas != nil {
		for _, pair := range document.Components.Schemas.AdditionalProperties {
			summarizeOpenAPIv3Schema(summary, pair.Value)
		}
	}
	for _, pair := range document.Paths.Path {
		summary.PathCount++
		v := pair.Value
		if v.Get != nil {
			summary.GetCount++
		}
		if v.Post != nil {
			summary.PostCount++
		}
		if v.Put != nil {
			summary.PutCount++
		}
		if v.Delete != nil {
			summary.DeleteCount++
		}
	}
	return summary
}

func summarizeOpenAPIv3Schema(summary *metrics.Complexity, schemaOrReference *openapi_v3.SchemaOrReference) {
	summary.SchemaCount++
	schema := schemaOrReference.GetSchema()
	if schema != nil && schema.Properties != nil {
		for _, pair := range schema.Properties.AdditionalProperties {
			summary.SchemaPropertyCount++
			summarizeOpenAPIv3Schema(summary, pair.Value)
		}
	}
}

func SummarizeDiscoveryDocument(document *discovery.Document) *metrics.Complexity {
	summary := &metrics.Complexity{}
	if document.Schemas != nil && document.Schemas.AdditionalProperties != nil {
		for _, pair := range document.Schemas.AdditionalProperties {
			summarizeDiscoverySchema(summary, pair.Value)
		}
	}
	if document.Resources != nil {
		for _, pair := range document.Resources.AdditionalProperties {
			summarizeDiscoveryResource(summary, pair.Value)
		}
	}
	if document.Methods != nil {
		for _, pair := range document.Methods.AdditionalProperties {
			summary.PathCount++
			v := pair.Value
			switch v.HttpMethod {
			case "GET":
				summary.GetCount++
			case "POST":
				summary.PostCount++
			case "PUT":
				summary.PutCount++
			case "DELETE":
				summary.DeleteCount++
			}
		}
	}
	return summary
}

func summarizeDiscoverySchema(summary *metrics.Complexity, schema *discovery.Schema) {
	summary.SchemaCount++
	if schema != nil && schema.Properties != nil {
		for _, pair := range schema.Properties.AdditionalProperties {
			summary.SchemaPropertyCount++
			summarizeDiscoverySchema(summary, pair.Value)
		}
	}
}

func summarizeDiscoveryResource(summary *metrics.Complexity, resource *discovery.Resource) {
	if resource.Resources != nil {
		for _, pair := range resource.Resources.AdditionalProperties {
			summarizeDiscoveryResource(summary, pair.Value)
		}
	}
	if resource.Methods != nil {
		for _, pair := range resource.Methods.AdditionalProperties {
			summary.PathCount++
			v := pair.Value
			switch v.HttpMethod {
			case "GET":
				summary.GetCount++
			case "POST":
				summary.PostCount++
			case "PUT":
				summary.PutCount++
			case "DELETE":
				summary.DeleteCount++
			}
		}
	}
}
