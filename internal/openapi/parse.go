// SPDX-License-Identifier: MIT

package openapi

import (
	"encoding/json"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/internal/ast"
	"github.com/caixw/apidoc/v7/internal/locale"
	"github.com/caixw/apidoc/v7/internal/vars"
)

// 将 doc.APIDoc 转换成 openapi
func convert(doc *ast.APIDoc) (*OpenAPI, error) {
	langID := doc.Lang.V()
	if langID == "" {
		langID = "und"
	}

	openapi := &OpenAPI{
		OpenAPI: LatestVersion,
		Info: &Info{
			Title:       doc.Title.Content.Value,
			Description: doc.Description.V(),
			Contact:     newContact(doc.Contact),
			License:     newLicense(doc.License),
			Version:     doc.Version.V(),
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

	if err := openapi.sanitize(); err != nil {
		return nil, err
	}
	return openapi, nil
}

func parsePaths(openapi *OpenAPI, d *ast.APIDoc) *core.SyntaxError {
	for _, api := range d.Apis {
		p := openapi.Paths[api.Path.Path.V()]
		if p == nil {
			p = &PathItem{}
			openapi.Paths[api.Path.Path.V()] = p
		}

		operation, err := setOperation(p, api.Method.V())
		if err != nil {
			err.Field = "paths." + err.Field
			return err
		}

		if len(api.Tags) > 0 {
			operation.Tags = make([]string, 0, len(api.Tags))
			for _, tag := range api.Tags {
				operation.Tags = append(operation.Tags, tag.Content.Value)
			}
		}
		operation.Deprecated = api.Deprecated != nil
		if api.ID != nil {
			operation.OperationID = api.ID.V()
		}
		if api.Summary != nil {
			operation.Summary = api.Summary.V()
		}
		if api.Description != nil {
			operation.Description = api.Description.V()
		}
		setOperationParams(operation, api)

		// servers
		// 不为 PathItem 设置 servers，直接写在 operation
		operation.Servers = make([]*Server, 0, len(api.Servers))
		for _, srv := range api.Servers {
			// 找到对应的 doc.Server.URL 值，之后根据此值从 openapi 中取 Server 对象
			var srvURL string
			for _, docSrv := range d.Servers {
				if docSrv.Name != nil && docSrv.Name.V() == srv.Content.Value {
					srvURL = docSrv.URL.V()
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
					examples[exp.Mimetype.V()] = &Example{
						Value: ExampleValue(exp.Content.Value.Value),
					}
				}

				content[r.Mimetype.V()] = &MediaType{
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
			status := strconv.Itoa(resp.Status.V())
			r, found := operation.Responses[status]
			if !found {
				r = &Response{
					Description: getDescription(resp.Description, resp.Summary),
					Headers:     make(map[string]*Header, 10),
					Content:     make(map[string]*MediaType, 10),
				}
				operation.Responses[status] = r
			}

			for _, h := range resp.Headers {
				r.Headers[h.Name.V()] = &Header{
					Style:       Style{Style: StyleSimple},
					Description: getDescription(h.Description, h.Summary),
				}
			}

			examples := make(map[string]*Example, len(resp.Examples))
			for _, exp := range resp.Examples {
				examples[exp.Mimetype.V()] = &Example{
					Summary: exp.Summary.V(),
					Value:   ExampleValue(exp.Content.Value.Value),
				}
			}
			r.Content[resp.Mimetype.V()] = &MediaType{
				Schema:   newSchemaFromRequest(resp, true),
				Examples: examples,
			}
		}
	} // end for doc.Apis

	return nil
}

func setOperationParams(operation *Operation, api *ast.API) {
	l := len(api.Path.Params) + len(api.Path.Queries)
	operation.Parameters = make([]*Parameter, 0, l)

	for _, param := range api.Path.Params {
		operation.Parameters = append(operation.Parameters, &Parameter{
			Name:        param.Name.V(),
			IN:          ParameterINPath,
			Description: getDescription(param.Description, param.Summary),
			Required:    !param.Optional.V(),
			Schema:      newSchema(param, true),
		})
	}

	for _, param := range api.Path.Queries {
		operation.Parameters = append(operation.Parameters, &Parameter{
			Name:        param.Name.V(),
			IN:          ParameterINQuery,
			Description: getDescription(param.Description, param.Summary),
			Required:    !param.Optional.V(),
			Schema:      newSchema(param, true),
		})
	}

	// 将各个类型的 Request 中的报头都集中到 operation.Parameters
	for _, r := range api.Requests {
		for _, param := range r.Headers {
			operation.Parameters = append(operation.Parameters, &Parameter{
				Style:       Style{Style: StyleSimple},
				Name:        param.Name.V(),
				IN:          ParameterINHeader,
				Description: getDescription(param.Description, param.Summary),
			})
		}
	}
}

func getDescription(desc *ast.Richtext, summary *ast.Attribute) string {
	if desc.V() != "" {
		return desc.V()
	}

	if summary != nil {
		return summary.V()
	}

	return ""
}

func setOperation(path *PathItem, method string) (*Operation, *core.SyntaxError) {
	operation := &Operation{}

	switch strings.ToUpper(method) {
	case "GET":
		if path.Get != nil {
			return nil, core.NewSyntaxError(core.Location{}, "get", locale.ErrDuplicateValue)
		}
		path.Get = operation
	case "DELETE":
		if path.Delete != nil {
			return nil, core.NewSyntaxError(core.Location{}, "delete", locale.ErrDuplicateValue)
		}
		path.Delete = operation
	case "POST":
		if path.Post != nil {
			return nil, core.NewSyntaxError(core.Location{}, "post", locale.ErrDuplicateValue)
		}
		path.Post = operation
	case "PUT":
		if path.Put != nil {
			return nil, core.NewSyntaxError(core.Location{}, "put", locale.ErrDuplicateValue)
		}
		path.Put = operation
	case "PATCH":
		if path.Patch != nil {
			return nil, core.NewSyntaxError(core.Location{}, "patch", locale.ErrDuplicateValue)
		}
		path.Patch = operation
	case "OPTIONS":
		if path.Options != nil {
			return nil, core.NewSyntaxError(core.Location{}, "options", locale.ErrDuplicateValue)
		}
		path.Options = operation
	case "HEAD":
		if path.Head != nil {
			return nil, core.NewSyntaxError(core.Location{}, "head", locale.ErrDuplicateValue)
		}
		path.Head = operation
	case "TRACE":
		if path.Trace != nil {
			return nil, core.NewSyntaxError(core.Location{}, "trace", locale.ErrDuplicateValue)
		}
		path.Trace = operation
	}

	return operation, nil
}

// JSON 输出 JSON 格式数据
func JSON(doc *ast.APIDoc) ([]byte, error) {
	openapi, err := convert(doc)
	if err != nil {
		return nil, err
	}

	return json.MarshalIndent(openapi, "", "\t")
}

// YAML 输出 YAML 格式数据
func YAML(doc *ast.APIDoc) ([]byte, error) {
	openapi, err := convert(doc)
	if err != nil {
		return nil, err
	}

	return yaml.Marshal(openapi)
}
