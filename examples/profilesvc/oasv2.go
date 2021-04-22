// Code generated by kok; DO NOT EDIT.
// github.com/RussellLuo/kok

package profilesvc

import (
	"reflect"

	"github.com/RussellLuo/kok/pkg/oasv2"
)

var (
	base = `swagger: "2.0"
info:
  version: "1.0.0"
  title: "Swagger Example"
  description: "Service is a simple CRUD interface for user profiles."
  license:
    name: "MIT"
host: "example.com"
basePath: "/api"
schemes:
  - "https"
consumes:
  - "application/json"
produces:
  - "application/json"
`

	paths = `
paths:
  /profiles/{id}/addresses/{addressID}:
    delete:
      description: ""
      operationId: "DeleteAddress"
      parameters:
        - name: id
          in: path
          required: true
          type: string
          description: ""
        - name: addressID
          in: path
          required: true
          type: string
          description: ""  
      %s
    get:
      description: ""
      operationId: "GetAddress"
      parameters:
        - name: id
          in: path
          required: true
          type: string
          description: ""
        - name: addressID
          in: path
          required: true
          type: string
          description: ""  
      %s 
  /profiles/{id}:
    delete:
      description: ""
      operationId: "DeleteProfile"
      parameters:
        - name: id
          in: path
          required: true
          type: string
          description: ""  
      %s
    get:
      description: ""
      operationId: "GetProfile"
      parameters:
        - name: id
          in: path
          required: true
          type: string
          description: ""  
      %s
    patch:
      description: ""
      operationId: "PatchProfile"
      parameters:
        - name: id
          in: path
          required: true
          type: string
          description: "" 
        - name: body
          in: body
          schema:
            $ref: "#/definitions/PatchProfileRequestBody" 
      %s
    put:
      description: ""
      operationId: "PutProfile"
      parameters:
        - name: id
          in: path
          required: true
          type: string
          description: "" 
        - name: body
          in: body
          schema:
            $ref: "#/definitions/PutProfileRequestBody" 
      %s 
  /profiles/{id}/addresses:
    get:
      description: ""
      operationId: "GetAddresses"
      parameters:
        - name: id
          in: path
          required: true
          type: string
          description: ""  
      %s
    post:
      description: ""
      operationId: "PostAddress"
      parameters:
        - name: id
          in: path
          required: true
          type: string
          description: "" 
        - name: body
          in: body
          schema:
            $ref: "#/definitions/PostAddressRequestBody" 
      %s 
  /profiles:
    post:
      description: ""
      operationId: "PostProfile"
      parameters: 
        - name: body
          in: body
          schema:
            $ref: "#/definitions/PostProfileRequestBody" 
      %s  
`
)

func getResponses(schema oasv2.Schema) []oasv2.OASResponses {
	return []oasv2.OASResponses{
		oasv2.GetOASResponses(schema, "DeleteAddress", 200, &DeleteAddressResponse{}),
		oasv2.GetOASResponses(schema, "GetAddress", 200, &GetAddressResponse{}),
		oasv2.GetOASResponses(schema, "DeleteProfile", 200, &DeleteProfileResponse{}),
		oasv2.GetOASResponses(schema, "GetProfile", 200, &GetProfileResponse{}),
		oasv2.GetOASResponses(schema, "PatchProfile", 200, &PatchProfileResponse{}),
		oasv2.GetOASResponses(schema, "PutProfile", 200, &PutProfileResponse{}),
		oasv2.GetOASResponses(schema, "GetAddresses", 200, &GetAddressesResponse{}),
		oasv2.GetOASResponses(schema, "PostAddress", 200, &PostAddressResponse{}),
		oasv2.GetOASResponses(schema, "PostProfile", 200, &PostProfileResponse{}),
	}
}

func getDefinitions(schema oasv2.Schema) map[string]oasv2.Definition {
	defs := make(map[string]oasv2.Definition)

	oasv2.AddResponseDefinitions(defs, schema, "DeleteAddress", 200, (&DeleteAddressResponse{}).Body())

	oasv2.AddResponseDefinitions(defs, schema, "DeleteProfile", 200, (&DeleteProfileResponse{}).Body())

	oasv2.AddResponseDefinitions(defs, schema, "GetAddress", 200, (&GetAddressResponse{}).Body())

	oasv2.AddResponseDefinitions(defs, schema, "GetAddresses", 200, (&GetAddressesResponse{}).Body())

	oasv2.AddResponseDefinitions(defs, schema, "GetProfile", 200, (&GetProfileResponse{}).Body())

	oasv2.AddDefinition(defs, "PatchProfileRequestBody", reflect.ValueOf(&struct {
		Profile Profile `json:"profile"`
	}{}))
	oasv2.AddResponseDefinitions(defs, schema, "PatchProfile", 200, (&PatchProfileResponse{}).Body())

	oasv2.AddDefinition(defs, "PostAddressRequestBody", reflect.ValueOf(&struct {
		Address Address `json:"address"`
	}{}))
	oasv2.AddResponseDefinitions(defs, schema, "PostAddress", 200, (&PostAddressResponse{}).Body())

	oasv2.AddDefinition(defs, "PostProfileRequestBody", reflect.ValueOf(&struct {
		Profile Profile `json:"profile"`
	}{}))
	oasv2.AddResponseDefinitions(defs, schema, "PostProfile", 200, (&PostProfileResponse{}).Body())

	oasv2.AddDefinition(defs, "PutProfileRequestBody", reflect.ValueOf(&struct {
		Profile Profile `json:"profile"`
	}{}))
	oasv2.AddResponseDefinitions(defs, schema, "PutProfile", 200, (&PutProfileResponse{}).Body())

	return defs
}

func OASv2APIDoc(schema oasv2.Schema) string {
	resps := getResponses(schema)
	paths := oasv2.GenPaths(resps, paths)

	defs := getDefinitions(schema)
	definitions := oasv2.GenDefinitions(defs)

	return base + paths + definitions
}
