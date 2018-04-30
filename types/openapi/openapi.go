// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package openapi 实现 openapi 的相关数据类型
//
// https://github.com/OAI/OpenAPI-Specification
package openapi

// TODO 扩展字段未加

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

// Info 接口文档的基本信息
type Info struct {
	Title          string   `json:"title" yaml:"title"`
	Description    string   `json:"description,omitempty" yaml:"description,omitempty"`
	TermsOfService string   `json:"termsOfService,omitempty" json:"termsOfService,omitempty"`
	Contact        *Contact `json:"contact,omitempty" yaml:"contact,omitempty"`
	License        *License `json:"license,omitempty" yaml:"license,omitempty"`
	Version        string   `json:"version" yaml:"version"`
}

// Contact 描述联系方式
type Contact struct {
	Name  string `json:"name,omitempty" yaml:"name,omitempty"`
	URL   string `json:"url,omitempty" yaml:"url,omitempty"`
	Email string `json:"email,omitempty" yaml:"email,omitempty"`
}

// License 授权信息
type License struct {
	Name string `json:"name" yaml:"name"`
	URL  string `json:"url,omitempty" yaml:"url,omitempty"`
}

// Server 服务器描述信息
type Server struct {
	URL         string                     `json:"url" yaml:"url"`
	Description string                     `json:"description,omitempty" yaml:"description,omitempty"`
	Variables   map[string]*ServerVariable `json:"variables,omitempty" yaml:"variables,omitempty"`
}

// ServerVariable Server 中 URL 模板中对应的参数变量值
type ServerVariable struct {
	Enum        []string `json:"enum,omitempty" yaml:"enum,omitempty"`
	Default     string   `json:"default" yaml:"default"`
	Description string   `json:"description,omitempty" yaml:"description,omitempty"`
}

// PathItem 每一条路径的详细描述信息
type PathItem struct {
	Ref         string       `json:"ref,omitempty" yaml:"ref,omitempty"`
	Summary     string       `json:"summary,omitempty" yaml:"summary,omitempty"`
	Description string       `json:"description,omitempty" yaml:"description,omitempty"`
	Get         *Operation   `json:"get,omitempty" yaml:"get,omitempty"`
	Put         *Operation   `json:"put,omitempty" yaml:"put,omitempty"`
	Post        *Operation   `json:"post,omitempty" yaml:"post,omitempty"`
	Delete      *Operation   `json:"delete,omitempty" yaml:"delete,omitempty"`
	Options     *Operation   `json:"options,omitempty" yaml:"options,omitempty"`
	Head        *Operation   `json:"head,omitempty" yaml:"head,omitempty"`
	Patch       *Operation   `json:"patch,omitempty" yaml:"patch,omitempty"`
	Trace       *Operation   `json:"trace,omitempty" yaml:"trace,omitempty"`
	Servers     []*Server    `json:"servers,omitempty" yaml:"servers,omitempty"`
	Parameters  []*Parameter `json:"parameters,omitempty" yaml:"parameters,omitempty"`
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

// Operation 描述对某一个资源的操作具体操作
type Operation struct {
	Tags         []string               `json:"tags,omitempty" yaml:"tags,omitempty"`
	Summary      string                 `json:"summary,omitempty" yaml:"summary,omitempty"`
	Description  string                 `json:"description,omitempty" yaml:"description,omitempty"`
	ExternalDocs *ExternalDocumentation `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`
	OperationID  string                 `json:"operationId,omitempty" yaml:"operationId,omitempty" `
	Parameters   []*Parameter           `json:"parameters,omitempty" yaml:"parameters,omitempty"`
	RequestBody  *RequestBody           `json:"requestBody,omitempty" yaml:"requestBody,omitempty"`
	Responses    map[string]*Response   `json:"responses" yaml:"responses"`
	Callbacks    map[string]*Callback   `json:"callbacks,omitempty" yaml:"callbacks,omitempty"`
	Deprecated   bool                   `json:"deprecated,omitempty" yaml:"deprecated,omitempty"`
	Security     []*SecurityRequirement `json:"security,omitempty" yaml:"security,omitempty"`
	Servers      []*Server              `json:"servers,omitempty" yaml:"servers,omitempty"`
}

// Parameter 参数信息
// 可同时作用于路径参数、请求参数、报头内容和 Cookie 值。
type Parameter struct {
	Style
	Name            string                `json:"name,omitempty" yaml:"name,omitempty"`
	IN              string                `json:"in" yaml:"in"`
	Description     string                `json:"description,omitempty" yaml:"description,omitempty"`
	Required        bool                  `json:"required,omitempty" yaml:"required,omitempty"`
	Deprecated      bool                  `json:"deprecated,omitempty" yaml:"deprecated,omitempty"`
	AllowEmptyValue bool                  `json:"allowEmptyValue,omitempty" yaml:"allowEmptyValue,omitempty"`
	Schema          *Schema               `json:"schema,omitemtpy" yaml:"schema,omitemtpy"`
	Example         ExampleValue          `json:"example,omitempty" yaml:"example,omitempty"`
	Examples        map[string]*Example   `json:"examples,omitempty" yaml:"examples,omitempty"`
	Content         map[string]*MediaType `json:"content,omitempty" yaml:"content,omitempty"`
}

// RequestBody 请求内容
type RequestBody struct {
	Description string                `json:"description,omitempty" yaml:"description,omitempty"`
	Content     map[string]*MediaType `json:"content" yaml:"content"`
	Required    bool                  `json:"required,omitempty" yaml:"required,omitempty" `
}

// MediaType 媒体类型
type MediaType struct {
	Schema   *Schema              `json:"schema,omitempty" yaml:"schema,omitempty"`
	Example  ExampleValue         `json:"example,omitempty" yaml:"example,omitemtpy"`
	Examples map[string]*Example  `json:"examples,omitempty" yaml:"examples,omitempty"`
	Encoding map[string]*Encoding `json:"encoding,omitempty" yaml:"encoding,omitempty"`
}

// Schema 定义了输出和输出的数据类型
type Schema struct {
	Type        string      `json:"type,omitempty" yaml:"type,omitempty"`
	Items       *Schema     `json:"items,omitempty" yaml:"items,omitempty"`
	Properties  *Schema     `json:"properties,omitempty" yaml:"properties,omitempty"`
	Default     interface{} `json:"default,omitempty" yaml:"default,omitempty"`
	Description string      `json:"description,omitempty" yaml:"description,omitempty"`

	// NOTE: 仅声明了部分使用到的变量

	Nullable      bool                   `json:"nullable,omitempty" yaml:"nullable,omitempty"`
	Discriminator *Discriminator         `json:"discriminator,omitempty" yaml:"discriminator,omitempty"`
	ReadOnly      bool                   `json:"readOnly,omitempty" yaml:"readOnly,omitempty"`
	WriteOnly     bool                   `json:"writeOnly,omitempty" yaml:"writeOnly,omitempty"`
	XML           *XML                   `json:"xml,omitempty" yaml:"xml,omitempty"`
	ExternalDocs  *ExternalDocumentation `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`
	Example       ExampleValue           `json:"example,omitempty" yaml:"example,omitempty"`
	Deprecated    bool                   `json:"deprecated,omitempty" yaml:"deprecated,omitempty"`
}

// XML 将 Schema 转换为 XML 的相关声明
type XML struct {
	Name      string `json:"name,omitempty" yaml:"name,omitempty"`
	Namespace string `json:"namespace,omitempty" yaml:"namespace,omitempty"`
	Prefix    string `json:"prefix,omitempty" yaml:"prefix,omitempty"`
	Attribute bool   `json:"attribute,omitempty" yaml:"attribute,omitempty"`
	Wrapped   bool   `json:"wrapped,omitempty" yaml:"wrapped,omitempty"`
}

// Discriminator Object
//
// NOTE: 暂时未用到。
type Discriminator struct {
	PropertyName string            `json:"propertyName" yaml:"propertyName"`
	Mapping      map[string]string `json:"mapping,omitempty" yaml:"mapping,omitempty"`
}

// Encoding 定义编码
type Encoding struct {
	Style
	ContentType string             `json:"contentType,omitempty" yaml:"contentType,omitempty"`
	Headers     map[string]*Header `json:"headers,omitempty" yaml:"headers,omitempty"`
}

// Callback Object
//
// NOTE: 暂时未用到
type Callback PathItem

// Response 每个 API 的返回信息
type Response struct {
	Description string                `json:"description" yaml:"description"`
	Headers     map[string]*Header    `json:"headers,omitempty" yaml:"headers,omitempty"`
	Content     map[string]*MediaType `json:"content,omitempty" yaml:"content,omitempty"`
	Links       map[string]*Link      `json:"links,omitempty" yaml:"links,omitempty"`
}

// Link 链接信息
type Link struct {
	OperationRef string            `json:"operationRef,omitempty" yaml:"operationRef,omitempty"`
	OperationID  string            `json:"operationId,omitempty" yaml:"operationId,omitempty"`
	Parameters   map[string]string `json:"parameters,omitempty" yaml:"parameters,omitempty"`
	RequestBody  map[string]string `json:"requestBody,omitempty" yaml:"requestBody,omitempty"`
	Description  string            `json:"description,omitempty" yaml:"description,omitempty"`
	Server       *Server           `json:"server,omitempty" yaml:"server,omitempty"`
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
}

// ExampleValue 表示示例的内容类型。
type ExampleValue string

// Header 即 Parameter 的别名，但 Name 字段必须存在。
type Header Parameter

// Reference 引用类型
type Reference struct {
	Ref string `json:"$ref" yaml:"$ref"`
}
