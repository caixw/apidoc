// SPDX-License-Identifier: MIT

// Package openapi 实现 openapi 的相关数据类型
//
// https://github.com/OAI/OpenAPI-Specification
package openapi

import (
	"strconv"

	"github.com/issue9/is"
	"github.com/issue9/version"

	"github.com/caixw/apidoc/v5/doc"
	"github.com/caixw/apidoc/v5/internal/locale"
	"github.com/caixw/apidoc/v5/message"
)

// TODO 扩展字段未加

// LatestVersion openapi 最新的版本号
const LatestVersion = "3.0.1"

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
	SecuritySchemes map[string]*SecurityScheme `json:"securitySechemes,omitempty" yaml:"securitySechemes,omitempty"`
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

func newTag(tag *doc.Tag) *Tag {
	return &Tag{
		Name:        tag.Name,
		Description: tag.Description,
	}
}

// Sanitize 数据检测
func (oa *OpenAPI) Sanitize() *message.SyntaxError {
	if oa.OpenAPI == "" {
		oa.OpenAPI = LatestVersion
	}

	if !version.SemVerValid(oa.OpenAPI) {
		return message.NewLocaleError("", "openapi", 0, locale.ErrInvalidFormat)
	}

	if oa.Info == nil {
		return message.NewLocaleError("", "info", 0, locale.ErrRequired)
	}
	if err := oa.Info.Sanitize(); err != nil {
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
		if err := srv.Sanitize(); err != nil {
			err.Field = "servers[" + strconv.Itoa(index) + "]."
			return err
		}
	}

	if len(oa.Paths) == 0 {
		return message.NewLocaleError("", "paths", 0, locale.ErrRequired)
	}
	// TODO 验证 paths

	if oa.Components != nil {
		if err := oa.Components.Sanitize(); err != nil {
			err.Field = "components." + err.Field
			return err
		}
	}

	for index, item := range oa.Tags {
		if err := item.Sanitize(); err != nil {
			err.Field = "tags[" + strconv.Itoa(index) + "]." + err.Field
			return err
		}
	}

	if oa.ExternalDocs != nil {
		if err := oa.ExternalDocs.Sanitize(); err != nil {
			err.Field = "externalDocs." + err.Field
			return err
		}
	}

	return nil
}

// Sanitize 数据检测
func (c *Components) Sanitize() *message.SyntaxError {
	for key, item := range c.Schemas {
		if err := item.Sanitize(); err != nil {
			err.Field = "schemas[" + key + "]." + err.Field
			return err
		}
	}

	for key, item := range c.Responses {
		if err := item.Sanitize(); err != nil {
			err.Field = "response[" + key + "]." + err.Field
			return err
		}
	}

	for key, item := range c.Parameters {
		if err := item.Sanitize(); err != nil {
			err.Field = "parameters[" + key + "]." + err.Field
			return err
		}
	}

	for key, item := range c.RequestBodies {
		if err := item.Sanitize(); err != nil {
			err.Field = "requestBodies[" + key + "]." + err.Field
			return err
		}
	}

	for key, item := range c.Headers {
		if err := item.Sanitize(); err != nil {
			err.Field = "headers[" + key + "]." + err.Field
			return err
		}
	}

	for key, item := range c.Links {
		if err := item.Sanitize(); err != nil {
			err.Field = "links[" + key + "]." + err.Field
			return err
		}
	}

	return nil
}

// Sanitize 数据检测
func (ext *ExternalDocumentation) Sanitize() *message.SyntaxError {
	if !is.URL(ext.URL) {
		return message.NewLocaleError("", "url", 0, locale.ErrInvalidFormat)
	}

	return nil
}

// Sanitize 数据检测
func (l *Link) Sanitize() *message.SyntaxError {
	if err := l.Server.Sanitize(); err != nil {
		err.Field = "server." + err.Field
		return err
	}

	return nil
}

// Sanitize 数据检测
func (tag *Tag) Sanitize() *message.SyntaxError {
	if tag.Name == "" {
		return message.NewLocaleError("", "name", 0, locale.ErrInvalidFormat)
	}

	if tag.ExternalDocs != nil {
		if err := tag.ExternalDocs.Sanitize(); err != nil {
			err.Field = "externalDocs." + err.Field
			return err
		}
	}

	return nil
}
