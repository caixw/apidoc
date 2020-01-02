// SPDX-License-Identifier: MIT

package openapi

import (
	"encoding/json"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/caixw/apidoc/v5/doc"
	"github.com/caixw/apidoc/v5/internal/locale"
	"github.com/caixw/apidoc/v5/internal/vars"
	"github.com/caixw/apidoc/v5/message"
)

// 将 doc.Doc 转换成 openapi
func convert(doc *doc.Doc) (*OpenAPI, error) {
	langID := doc.Lang
	if langID == "" {
		langID = "und"
	}

	openapi := &OpenAPI{
		OpenAPI: LatestVersion,
		Info: &Info{
			Title:       doc.Title,
			Description: doc.Description.Text,
			Contact:     newContact(doc.Contact),
			License:     newLicense(doc.License),
			Version:     string(doc.Version),
		},
		Servers: make([]*Server, 0, len(doc.Servers)),
		Tags:    make([]*Tag, 0, len(doc.Tags)),
		Paths:   make(map[string]*PathItem, len(doc.Apis)),
		ExternalDocs: &ExternalDocumentation{
			Description: locale.Translate(langID, locale.GeneratorBy, vars.Name),
			URL:         vars.OfficialURL,
		},
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

func parsePaths(openapi *OpenAPI, d *doc.Doc) *message.SyntaxError {
	for _, api := range d.Apis {
		p := openapi.Paths[api.Path.Path]
		if p == nil {
			p = &PathItem{}
			openapi.Paths[api.Path.Path] = p
		}

		operation, err := setOperation(p, string(api.Method))
		if err != nil {
			err.Field = "paths." + err.Field
			return err
		}

		operation.Tags = api.Tags
		operation.Deprecated = api.Deprecated != ""
		operation.OperationID = api.ID
		operation.Summary = api.Summary
		operation.Description = api.Description.Text
		setOperationParams(operation, api)

		// servers
		// 不为 PathItem 设置 servers，直接写在 operation
		operation.Servers = make([]*Server, 0, len(api.Servers))
		for _, srv := range api.Servers {
			// 找到对应的 doc.Server.URL 值，之后根据此值从 openapi 中取 Server 对象
			var srvURL string
			for _, docSrv := range d.Servers {
				if docSrv.Name == srv {
					srvURL = docSrv.URL
					break
				}
			}

			if srvURL == "" {
				continue
			}

			for _, ss := range openapi.Servers {
				if ss.URL == srvURL {
					operation.Servers = append(operation.Servers, ss)
				}
			}
		}

		// requests
		if len(api.Requests) > 0 {
			content := make(map[string]*MediaType, len(api.Requests))
			for _, r := range api.Requests {
				examples := make(map[string]*Example, len(r.Examples))
				for _, exp := range r.Examples {
					examples[exp.Mimetype] = &Example{
						Value: ExampleValue(exp.Content),
					}
				}

				content[r.Mimetype] = &MediaType{
					Schema:   newSchemaFromRequest(r, true),
					Examples: examples,
				}
			}

			operation.RequestBody = &RequestBody{
				Content: content,
			}
		}

		// responses
		operation.Responses = make(map[string]*Response, len(api.Responses))
		for _, resp := range api.Responses {
			status := resp.Status.String()
			r, found := operation.Responses[status]
			if !found {
				r = &Response{}
				operation.Responses[status] = r
			}

			if r.Headers == nil {
				r.Headers = make(map[string]*Header, 10)
			}
			for _, h := range resp.Headers {
				desc := h.Description.Text
				if desc == "" {
					desc = h.Summary
				}
				r.Headers[h.Name] = &Header{
					Description: desc,
				}
			}

			if r.Content == nil {
				r.Content = make(map[string]*MediaType, 10)
			}
			examples := make(map[string]*Example, len(resp.Examples))
			for _, exp := range resp.Examples {
				desc := exp.Description.Text
				if desc == "" {
					desc = exp.Summary
				}
				examples[exp.Mimetype] = &Example{
					Summary: desc,
					Value:   ExampleValue(exp.Content),
				}
			}
			r.Content[resp.Mimetype] = &MediaType{
				Schema:   newSchemaFromRequest(resp, true),
				Examples: examples,
			}
		}
	} // end for doc.Apis

	return nil
}

func setOperationParams(operation *Operation, api *doc.API) {
	l := len(api.Path.Params) + len(api.Path.Queries)
	operation.Parameters = make([]*Parameter, 0, l)

	for _, param := range api.Path.Params {
		operation.Parameters = append(operation.Parameters, &Parameter{
			Name:        param.Name,
			IN:          ParameterINPath,
			Description: param.Summary,
			Required:    !param.Optional,
			Schema:      newSchema(param, true),
		})
	}

	for _, param := range api.Path.Queries {
		operation.Parameters = append(operation.Parameters, &Parameter{
			Name:        param.Name,
			IN:          ParameterINQuery,
			Description: param.Summary,
			Required:    !param.Optional,
			Schema:      newSchema(param, true),
		})
	}

	// 将各个类型的 Request 中的报头都集中到 operation.Parameters
	for _, r := range api.Requests {
		for _, param := range r.Headers {
			desc := param.Description.Text
			if desc == "" {
				desc = param.Summary
			}
			operation.Parameters = append(operation.Parameters, &Parameter{
				Name:        param.Name,
				IN:          ParameterINHeader,
				Description: desc,
			})
		}
	}
}

func setOperation(path *PathItem, method string) (*Operation, *message.SyntaxError) {
	operation := &Operation{}

	switch strings.ToUpper(method) {
	case "GET":
		if path.Get != nil {
			return nil, message.NewLocaleError("", "get", 0, locale.ErrDuplicateValue)
		}
		path.Get = operation
	case "DELETE":
		if path.Delete != nil {
			return nil, message.NewLocaleError("", "delete", 0, locale.ErrDuplicateValue)
		}
		path.Delete = operation
	case "POST":
		if path.Post != nil {
			return nil, message.NewLocaleError("", "post", 0, locale.ErrDuplicateValue)
		}
		path.Post = operation
	case "PUT":
		if path.Put != nil {
			return nil, message.NewLocaleError("", "put", 0, locale.ErrDuplicateValue)
		}
		path.Put = operation
	case "PATCH":
		if path.Patch != nil {
			return nil, message.NewLocaleError("", "patch", 0, locale.ErrDuplicateValue)
		}
		path.Patch = operation
	case "OPTIONS":
		if path.Options != nil {
			return nil, message.NewLocaleError("", "options", 0, locale.ErrDuplicateValue)
		}
		path.Options = operation
	case "HEAD":
		if path.Head != nil {
			return nil, message.NewLocaleError("", "head", 0, locale.ErrDuplicateValue)
		}
		path.Head = operation
	case "TRACE":
		if path.Trace != nil {
			return nil, message.NewLocaleError("", "trace", 0, locale.ErrDuplicateValue)
		}
		path.Trace = operation
	}

	return operation, nil
}

// JSON 输出 JSON 格式数据
func JSON(doc *doc.Doc) ([]byte, error) {
	openapi, err := convert(doc)
	if err != nil {
		return nil, err
	}

	return json.MarshalIndent(openapi, "", "\t")
}

// YAML 输出 YAML 格式数据
func YAML(doc *doc.Doc) ([]byte, error) {
	openapi, err := convert(doc)
	if err != nil {
		return nil, err
	}

	return yaml.Marshal(openapi)
}
