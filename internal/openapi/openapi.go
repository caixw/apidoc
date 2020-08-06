// SPDX-License-Identifier: MIT

// Package openapi 实现 openapi 的相关数据类型
//
// https://github.com/OAI/OpenAPI-Specification
package openapi

import (
	"strconv"

	"github.com/issue9/is"
	"github.com/issue9/version"

	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/internal/ast"
	"github.com/caixw/apidoc/v7/internal/locale"
)

// LatestVersion openapi 最新的版本号
const LatestVersion = "3.0.3"

// OpenAPI openAPI 的根对象
type OpenAPI struct {
	OpenAPI      string                 `json:"openapi" yaml:"openapi"`
	Info         *Info                  `json:"info" yaml:"info"`
	Servers      []*Server              `json:"servers,omitempty" yaml:"servers,omitempty"`
	Paths        map[string]*PathItem   `json:"paths" yaml:"paths"`
	Components   *Components            `json:"components,omitempty" yaml:"components,omitempty"`
	Security     []*SecurityRequirement `json:"security,omitempty" yaml:"security,omitempty"`
	Tags         []*Tag                 `json:"tags,omitempty" yaml:"tags,omitempty"`
	ExternalDocs *ExternalDocumentation `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`
}

// Components 可复用的对象
type Components struct {
	Schemas         map[string]*Schema         `json:"schemas,omitempty" yaml:"schemas,omitempty"`
	Responses       map[string]*Response       `json:"responses,omitempty" yaml:"responses,omitempty"`
	Parameters      map[string]*Parameter      `json:"parameters,omitempty" yaml:"parameters,omitempty"`
	Examples        map[string]*Example        `json:"examples,omitempty" yaml:"examples,omitempty"`
	RequestBodies   map[string]*RequestBody    `json:"requestBodies,omitempty" yaml:"requestBodies,omitempty"`
	Headers         map[string]*Header         `json:"headers,omitempty" yaml:"headers,omitempty"`
	SecuritySchemes map[string]*SecurityScheme `json:"securitySchemes,omitempty" yaml:"securitySchemes,omitempty"`
	Links           map[string]*Link           `json:"links,omitempty" yaml:"links,omitempty"`
	Callbacks       map[string]*Callback       `json:"callbacks,omitempty" yaml:"callbacks,omitempty"`
}

// ExternalDocumentation 引用外部资源的扩展文档
type ExternalDocumentation struct {
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	URL         string `json:"url" yaml:"url"`
}

// Link 链接信息
type Link struct {
	OperationRef string            `json:"operationRef,omitempty" yaml:"operationRef,omitempty"`
	OperationID  string            `json:"operationId,omitempty" yaml:"operationId,omitempty"`
	Parameters   map[string]string `json:"parameters,omitempty" yaml:"parameters,omitempty"`
	RequestBody  map[string]string `json:"requestBody,omitempty" yaml:"requestBody,omitempty"`
	Description  string            `json:"description,omitempty" yaml:"description,omitempty"`
	Server       *Server           `json:"server,omitempty" yaml:"server,omitempty"`

	Ref string `json:"$ref,omitempty" yaml:"$ref,omitempty"`
}

// Tag 标签内容
type Tag struct {
	Name         string                 `json:"name" yaml:"name"`
	Description  string                 `json:"description,omitempty" yaml:"description,omitempty"`
	ExternalDocs *ExternalDocumentation `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`
}

// Example 示例代码
type Example struct {
	Summary       string       `json:"summary,omitempty" yaml:"summary,omitempty"`
	Description   string       `json:"description,omitempty" yaml:"description,omitempty"`
	Value         ExampleValue `json:"value,omitempty" yaml:"value,omitempty"`
	ExternalValue string       `json:"external,omitempty" yaml:"external,omitempty"`

	Ref string `json:"$ref,omitempty" yaml:"$ref,omitempty"`
}

// ExampleValue 表示示例的内容类型。
type ExampleValue string

func newTag(tag *ast.Tag) *Tag {
	return &Tag{
		Name:        tag.Name.V(),
		Description: tag.Title.V(),
	}
}

func (oa *OpenAPI) sanitize() *core.Error {
	if oa.OpenAPI == "" {
		oa.OpenAPI = LatestVersion
	}

	if !version.SemVerValid(oa.OpenAPI) {
		return core.NewError(locale.ErrInvalidFormat).WithField("openapi")
	}

	if oa.Info == nil {
		return core.NewError(locale.ErrIsEmpty, "info").WithField("info")
	}
	if err := oa.Info.sanitize(); err != nil {
		err.Field = "info." + err.Field
		return err
	}

	// 没有，则采用默认值
	if len(oa.Servers) == 0 {
		oa.Servers = []*Server{{
			URL: "/",
		}}
	}

	for index, srv := range oa.Servers {
		if err := srv.sanitize(); err != nil {
			err.Field = "servers[" + strconv.Itoa(index) + "]."
			return err
		}
	}

	if len(oa.Paths) == 0 {
		return core.NewError(locale.ErrIsEmpty, "paths").WithField("paths")
	}
	for k, path := range oa.Paths {
		if err := path.sanitize(); err != nil {
			err.Field = "paths[" + k + "]." + err.Field
			return err
		}
	}

	if oa.Components != nil {
		if err := oa.Components.sanitize(); err != nil {
			err.Field = "components." + err.Field
			return err
		}
	}

	for index, item := range oa.Tags {
		if err := item.sanitize(); err != nil {
			err.Field = "tags[" + strconv.Itoa(index) + "]." + err.Field
			return err
		}
	}

	if oa.ExternalDocs != nil {
		if err := oa.ExternalDocs.sanitize(); err != nil {
			err.Field = "externalDocs." + err.Field
			return err
		}
	}

	return nil
}

func (c *Components) sanitize() *core.Error {
	for key, item := range c.Schemas {
		if err := item.sanitize(); err != nil {
			err.Field = "schemas[" + key + "]." + err.Field
			return err
		}
	}

	for key, item := range c.Responses {
		if err := item.sanitize(); err != nil {
			err.Field = "response[" + key + "]." + err.Field
			return err
		}
	}

	for key, item := range c.Parameters {
		if err := item.sanitize(); err != nil {
			err.Field = "parameters[" + key + "]." + err.Field
			return err
		}
	}

	for key, item := range c.RequestBodies {
		if err := item.sanitize(); err != nil {
			err.Field = "requestBodies[" + key + "]." + err.Field
			return err
		}
	}

	for key, item := range c.Headers {
		if err := item.sanitize(); err != nil {
			err.Field = "headers[" + key + "]." + err.Field
			return err
		}
	}

	for key, item := range c.Links {
		if err := item.sanitize(); err != nil {
			err.Field = "links[" + key + "]." + err.Field
			return err
		}
	}

	return nil
}

func (ext *ExternalDocumentation) sanitize() *core.Error {
	if !is.URL(ext.URL) {
		return core.NewError(locale.ErrInvalidFormat).WithField("url")
	}

	return nil
}

func (l *Link) sanitize() *core.Error {
	if err := l.Server.sanitize(); err != nil {
		err.Field = "server." + err.Field
		return err
	}

	return nil
}

func (tag *Tag) sanitize() *core.Error {
	if tag.Name == "" {
		return core.NewError(locale.ErrInvalidFormat).WithField("name")
	}

	if tag.ExternalDocs != nil {
		if err := tag.ExternalDocs.sanitize(); err != nil {
			err.Field = "externalDocs." + err.Field
			return err
		}
	}

	return nil
}
