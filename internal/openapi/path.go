// SPDX-License-Identifier: MIT

package openapi

import (
	"github.com/caixw/apidoc/v5/internal/locale"
	"github.com/caixw/apidoc/v5/message"
)

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

// RequestBody 请求内容
type RequestBody struct {
	Description string                `json:"description,omitempty" yaml:"description,omitempty"`
	Content     map[string]*MediaType `json:"content" yaml:"content"`
	Required    bool                  `json:"required,omitempty" yaml:"required,omitempty" `

	Ref string `json:"$ref,omitempty" yaml:"$ref,omitempty"`
}

// MediaType 媒体类型
type MediaType struct {
	Schema   *Schema              `json:"schema,omitempty" yaml:"schema,omitempty"`
	Example  ExampleValue         `json:"example,omitempty" yaml:"example,omitempty"`
	Examples map[string]*Example  `json:"examples,omitempty" yaml:"examples,omitempty"`
	Encoding map[string]*Encoding `json:"encoding,omitempty" yaml:"encoding,omitempty"`
}

// Encoding 定义编码
//
// 对父对象中的 Schema 中的一些字段的特殊定义
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

	Ref string `json:"$ref,omitempty" yaml:"$ref,omitempty"`
}

func (path *PathItem) sanitize() *message.SyntaxError {
	var o *Operation
	switch {
	case path.Get != nil:
		o = path.Get
	case path.Put != nil:
		o = path.Put
	case path.Post != nil:
		o = path.Post
	case path.Delete != nil:
		o = path.Delete
	case path.Options != nil:
		o = path.Options
	case path.Head != nil:
		o = path.Head
	case path.Patch != nil:
		o = path.Patch
	case path.Trace != nil:
		o = path.Trace
	}

	if o == nil {
		return message.NewLocaleError("", "operation", 0, locale.ErrRequired)

	}

	if err := o.sanitize(); err != nil {
		err.Field = "method." + err.Field
		return err
	}
	return nil
}

func (o *Operation) sanitize() *message.SyntaxError {
	if len(o.Responses) == 0 {
		return message.NewLocaleError("", "responses", 0, locale.ErrRequired)
	}
	return nil
}

func (req *RequestBody) sanitize() *message.SyntaxError {
	if len(req.Content) == 0 {
		return message.NewLocaleError("", "content", 0, locale.ErrRequired)
	}

	for key, mt := range req.Content {
		if err := mt.sanitize(); err != nil {
			err.Field = "content[" + key + "]." + err.Field
			return err
		}
	}

	return nil
}

func (resp *Response) sanitize() *message.SyntaxError {
	if resp.Description == "" {
		return message.NewLocaleError("", "description", 0, locale.ErrRequired)
	}

	for key, header := range resp.Headers {
		if err := header.sanitize(); err != nil {
			err.Field = "headers[" + key + "]." + err.Field
			return err
		}
	}

	for key, mt := range resp.Content {
		if err := mt.sanitize(); err != nil {
			err.Field = "content[" + key + "]." + err.Field
			return err
		}
	}

	for key, link := range resp.Links {
		if err := link.sanitize(); err != nil {
			err.Field = "links[" + key + "]." + err.Field
			return err
		}
	}

	return nil
}

func (mt *MediaType) sanitize() *message.SyntaxError {
	if mt.Schema != nil {
		if err := mt.sanitize(); err != nil {
			err.Field = "schema." + err.Field
			return err
		}
	}

	for key, en := range mt.Encoding {
		if err := en.sanitize(); err != nil {
			err.Field = "encoding[" + key + "]." + err.Field
			return err
		}
	}
	return nil
}

func (en *Encoding) sanitize() *message.SyntaxError {
	if err := en.Style.sanitize(); err != nil {
		return err
	}

	for key, header := range en.Headers {
		if err := header.sanitize(); err != nil {
			err.Field = "headers[" + key + "]." + err.Field
			return err
		}
	}

	return nil
}
