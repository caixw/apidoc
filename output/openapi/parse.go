// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package openapi

import (
	"encoding/json"
	"strings"

	yaml "gopkg.in/yaml.v2"

	"github.com/caixw/apidoc/doc"
)

func parse(doc *doc.Doc) (*OpenAPI, error) {
	openapi := &OpenAPI{
		OpenAPI: doc.APIDoc,
		Info: &Info{
			Title:       doc.Title,
			Description: doc.Content,
			Contact:     newContact(doc.Contact),
			License:     newLicense(doc.License),
			Version:     doc.Version,
		},
		Servers: make([]*Server, 0, len(doc.Servers)),
		Tags:    make([]*Tag, 0, len(doc.Tags)),
	}

	for _, srv := range doc.Servers {
		openapi.Servers = append(openapi.Servers, newServer(srv))
	}

	for _, tag := range doc.Tags {
		openapi.Tags = append(openapi.Tags, newTag(tag))
	}

	if err := parsePaths(openapi, doc); err != nil {
		return nil, err
	}

	return openapi, nil
}

func parsePaths(openapi *OpenAPI, doc *doc.Doc) *Error {
	for _, api := range doc.Apis {
		p := openapi.Paths[api.Path]
		if p == nil {
			p = &PathItem{}
			openapi.Paths[api.Path] = p
		}

		operation, err := setOperation(p, api.Method)
		if err != nil {
			err.Field = "paths." + err.Field
			return err
		}

		operation.Tags = api.Tags
		operation.Description = api.Description
		operation.Deprecated = api.Deprecated != ""
		setOperationParams(operation, api)

		if len(api.Requests) > 0 {
			content := make(map[string]*MediaType, len(api.Requests))
			for _, r := range api.Requests {
				examples := make(map[string]*Example, len(r.Examples))
				for _, exp := range r.Examples {
					examples[exp.Mimetype] = &Example{
						Summary: exp.Summary,
						Value:   ExampleValue(exp.Value),
					}
				}

				content[r.Mimetype] = &MediaType{
					Schema:   &Schema{Schema: r.Type},
					Examples: examples,
				}
			}

			operation.RequestBody = &RequestBody{
				Content: content,
			}
		}

		operation.Responses = make(map[string]*Response, len(api.Responses))
		for _, resp := range api.Responses {
			r, found := operation.Responses[resp.Status]
			if !found {
				r = &Response{}
				operation.Responses[resp.Status] = r
			}

			if r.Headers == nil {
				r.Headers = make(map[string]*Header, 10)
			}
			for _, h := range resp.Headers {
				r.Headers[h.Name] = &Header{
					Description: Description(h.Summary),
				}
			}

			if r.Content == nil {
				r.Content = make(map[string]*MediaType, 10)
			}
			examples := make(map[string]*Example, len(resp.Examples))
			for _, exp := range resp.Examples {
				examples[exp.Mimetype] = &Example{
					Summary: exp.Summary,
					Value:   ExampleValue(exp.Value),
				}
			}
			r.Content[resp.Mimetype] = &MediaType{
				Schema:   &Schema{Schema: resp.Type},
				Examples: examples,
			}
		}
	} // end for doc.Apis

	return nil
}

func setOperationParams(operation *Operation, api *doc.API) {
	operation.Parameters = make([]*Parameter, 0, len(api.Params)+len(api.Queries)+len(api.Requests[0].Headers))

	for _, param := range api.Params {
		operation.Parameters = append(operation.Parameters, &Parameter{
			Name:        param.Name,
			IN:          ParameterINPath,
			Description: Description(param.Summary),
			Required:    !param.Optional,
			Schema:      &Schema{Schema: param.Type},
		})
	}

	for _, param := range api.Queries {
		operation.Parameters = append(operation.Parameters, &Parameter{
			Name:        param.Name,
			IN:          ParameterINQuery,
			Description: Description(param.Summary),
			Required:    !param.Optional,
			Schema:      &Schema{Schema: param.Type},
		})
	}

	if len(api.Requests) > 0 {
		for _, param := range api.Requests[0].Headers {
			operation.Parameters = append(operation.Parameters, &Parameter{
				Name:        param.Name,
				IN:          ParameterINHeader,
				Description: Description(param.Summary),
				Required:    !param.Optional,
			})
		}
	}
}

func setOperation(path *PathItem, method string) (*Operation, *Error) {
	operation := &Operation{}

	switch strings.ToUpper(method) {
	case "GET":
		if path.Get != nil {
			return nil, &Error{Field: "Get", Message: "已经存在"}
		}
		path.Get = operation
	case "DELETE":
		if path.Delete != nil {
			return nil, &Error{Field: "Delete", Message: "已经存在"}
		}
		path.Delete = operation
	case "POST":
		if path.Post != nil {
			return nil, &Error{Field: "Post", Message: "已经存在"}
		}
		path.Post = operation
	case "PUT":
		if path.Put != nil {
			return nil, &Error{Field: "Put", Message: "已经存在"}
		}
		path.Put = operation
	case "PATCH":
		if path.Patch != nil {
			return nil, &Error{Field: "Patch", Message: "已经存在"}
		}
		path.Patch = operation
	case "OPTIONS":
		if path.Options != nil {
			return nil, &Error{Field: "Options", Message: "已经存在"}
		}
		path.Options = operation
	case "HEAD":
		if path.Head != nil {
			return nil, &Error{Field: "Head", Message: "已经存在"}
		}
		path.Head = operation
	case "TRACE":
		if path.Trace != nil {
			return nil, &Error{Field: "Trace", Message: "已经存在"}
		}
		path.Trace = operation
	}

	return operation, nil
}

// JSON 输出 JSON 格式数据
func JSON(doc *doc.Doc) ([]byte, error) {
	openapi, err := parse(doc)
	if err != nil {
		return nil, err
	}

	return json.Marshal(openapi)
}

// YAML 输出 YAML 格式数据
func YAML(doc *doc.Doc) ([]byte, error) {
	openapi, err := parse(doc)
	if err != nil {
		return nil, err
	}

	return yaml.Marshal(openapi)
}
