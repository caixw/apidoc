// SPDX-License-Identifier: MIT

package ast

import (
	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/internal/token"
)

// 表示支持的各种数据类型
const (
	TypeNone   string = "" // 空值表示不输出任何内容，仅用于 Request
	TypeBool          = "bool"
	TypeObject        = "object"
	TypeNumber        = "number"
	TypeString        = "string"
)

// 富文本可用的类型
const (
	RichtextTypeHTML     = "html"
	RichtextTypeMarkdown = "markdown"
)

type (
	// APIDoc 对应 apidoc 元素
	APIDoc struct {
		token.BaseTag
		URI      core.URI `apidoc:"-"`
		RootName struct{} `apidoc:"apidoc,meta,usage-apidoc"`

		// 程序的版本号
		//
		// 同时也作为文档格式的版本号。客户端可以依此值确定文档格式。
		// 仅用于输出，文档中不需要指定此值。
		APIDoc *APIDocVersionAttribute `apidoc:"apidoc,attr,usage-apidoc-apidoc,omitempty"`

		// 文档内容的区域信息
		// 如果存在此值，客户端应该尽量根据此值显示相应的界面语言。
		Lang *Attribute `apidoc:"lang,attr,usage-apidoc-lang,omitempty"`

		// 文档的图标
		//
		// 如果采用默认的 xsl 转换，会替换掉页面上的图标和 favicon 图标
		Logo *Attribute `apidoc:"logo,attr,usage-apidoc-logo,omitempty"`

		Created     *DateAttribute    `apidoc:"created,attr,usage-apidoc-created,omitempty"` // 文档的生成时间
		Version     *VersionAttribute `apidoc:"version,attr,usage-apidoc-version,omitempty"` // 文档的版本
		Title       *Element          `apidoc:"title,elem,usage-apidoc-title"`
		Description *Richtext         `apidoc:"description,elem,usage-apidoc-description,omitempty"`
		Contact     *Contact          `apidoc:"contact,elem,usage-apidoc-contact,omitempty"`
		License     *Link             `apidoc:"license,elem,usage-apidoc-license,omitempty"` // 版本信息
		Tags        []*Tag            `apidoc:"tag,elem,usage-apidoc-tags,omitempty"`        // 所有的标签
		Servers     []*Server         `apidoc:"server,elem,usage-apidoc-servers,omitempty"`
		APIs        []*API            `apidoc:"api,elem,usage-apidoc-apis,omitempty"`

		// 公共的报头，所有 API 默认都采用此报头
		Headers []*Param `apidoc:"header,elem,usage-apidoc-headers,omitempty"`

		// 表示所有 API 都有可能返回的内容
		Responses []*Request `apidoc:"response,elem,usage-apidoc-responses,omitempty"`

		// 表示所有接口都支持的文档类型
		Mimetypes []*Element `apidoc:"mimetype,elem,usage-apidoc-mimetypes"`

		// 指定命名空间的相关属性
		XMLNamespaces []*XMLNamespace `apidoc:"xml-namespace,elem,usage-apidoc-namespaces,omitempty"`
	}

	// XMLNamespace 定义命名空间的相关属性
	XMLNamespace struct {
		token.BaseTag
		RootName struct{}   `apidoc:"xml-namespace,meta,usage-xml-namespace"`
		Prefix   *Attribute `apidoc:"prefix,attr,usage-xml-namespace-prefix,omitempty"`
		URN      *Attribute `apidoc:"urn,attr,usage-xml-namespace-urn"`
	}

	// API 表示 <api> 顶层元素
	API struct {
		token.BaseTag
		RootName struct{} `apidoc:"api,meta,usage-api"`

		Version     *VersionAttribute `apidoc:"version,attr,usage-api-version,omitempty"`
		Method      *MethodAttribute  `apidoc:"method,attr,usage-api-method"`
		ID          *Attribute        `apidoc:"id,attr,usage-api-id,omitempty"`
		Path        *Path             `apidoc:"path,elem,usage-api-path"`
		Summary     *Attribute        `apidoc:"summary,attr,usage-api-summary,omitempty"`
		Description *Richtext         `apidoc:"description,elem,usage-api-description,omitempty"`
		Requests    []*Request        `apidoc:"request,elem,usage-api-requests,omitempty"` // 不同的 mimetype 可能会定义不同
		Responses   []*Request        `apidoc:"response,elem,usage-api-responses,omitempty"`
		Callback    *Callback         `apidoc:"callback,elem,usage-api-callback,omitempty"`
		Deprecated  *VersionAttribute `apidoc:"deprecated,attr,usage-api-deprecated,omitempty"`
		Headers     []*Param          `apidoc:"header,elem,usage-api-headers,omitempty"`

		Tags    []*Element `apidoc:"tag,elem,usage-api-tags,omitempty"`
		Servers []*Element `apidoc:"server,elem,usage-api-servers,omitempty"`

		URI core.URI `apidoc:"-"`
		doc *APIDoc
	}

	// Link 表示一个链接
	Link struct {
		token.BaseTag
		RootName struct{} `apidoc:"link,meta,usage-link"`

		Text *Attribute `apidoc:"text,attr,usage-link-text"`
		URL  *Attribute `apidoc:"url,attr,usage-link-url"`
	}

	// Contact 描述联系方式
	Contact struct {
		token.BaseTag
		RootName struct{} `apidoc:"contact,meta,usage-contact"`

		Name  *Attribute `apidoc:"name,attr,usage-contact-name"`
		URL   *Element   `apidoc:"url,elem,usage-contact-url,omitempty"`
		Email *Element   `apidoc:"email,elem,usage-contact-email,omitempty"`
	}

	// Callback 描述回调信息
	Callback struct {
		token.BaseTag
		RootName struct{} `apidoc:"callback,meta,usage-callback"`

		Method      *MethodAttribute  `apidoc:"method,attr,usage-callback-method"`
		Path        *Path             `apidoc:"path,elem,usage-callback-path,omitempty"`
		Summary     *Attribute        `apidoc:"summary,attr,usage-callback-summary,omitempty"`
		Description *Richtext         `apidoc:"description,elem,usage-callback-description,omitempty"`
		Deprecated  *VersionAttribute `apidoc:"deprecated,attr,usage-callback-deprecated,omitempty"`
		Responses   []*Request        `apidoc:"response,elem,usage-callback-responses,omitempty"`
		Requests    []*Request        `apidoc:"request,elem,usage-callback-requests"` // 至少一个
		Headers     []*Param          `apidoc:"header,elem,usage-callback-headers,omitempty"`
	}

	// Enum 表示枚举值
	Enum struct {
		token.BaseTag
		RootName struct{} `apidoc:"enum,meta,usage-enum"`

		Deprecated  *VersionAttribute `apidoc:"deprecated,attr,usage-enum-deprecated,omitempty"`
		Value       *Attribute        `apidoc:"value,attr,usage-enum-value"`
		Summary     *Attribute        `apidoc:"summary,attr,usage-enum-summary,omitempty"`
		Description *Richtext         `apidoc:"description,elem,usage-enum-description,omitempty"`
	}

	// Example 示例代码
	Example struct {
		token.BaseTag
		RootName struct{} `apidoc:"example,meta,usage-example"`

		Mimetype *Attribute    `apidoc:"mimetype,attr,usage-example-mimetype"`
		Content  *ExampleValue `apidoc:",cdata,usage-example-content"`
		Summary  *Attribute    `apidoc:"summary,attr,usage-example-summary,omitempty"`
	}

	// Param 表示参数类型
	Param struct {
		token.BaseTag
		RootName struct{} `apidoc:"param,meta,usage-param"`

		XML
		Name        *Attribute        `apidoc:"name,attr,usage-param-name"`
		Type        *TypeAttribute    `apidoc:"type,attr,usage-param-type"`
		Deprecated  *VersionAttribute `apidoc:"deprecated,attr,usage-param-deprecated,omitempty"`
		Default     *Attribute        `apidoc:"default,attr,usage-param-default,omitempty"`
		Optional    *BoolAttribute    `apidoc:"optional,attr,usage-param-optional,omitempty"`
		Array       *BoolAttribute    `apidoc:"array,attr,usage-param-array,omitempty"`
		Items       []*Param          `apidoc:"param,elem,usage-param-items,omitempty"`
		Summary     *Attribute        `apidoc:"summary,attr,usage-param-summary,omitempty"`
		Enums       []*Enum           `apidoc:"enum,elem,usage-param-enums,omitempty"`
		Description *Richtext         `apidoc:"description,elem,usage-param-description,omitempty"`

		// 数组参数是否展开
		//
		// 数组可以有以下两种展示方式：
		//  1. k=1&k=2
		//  2. k=1,2
		// 1 为默认方式，ArrayStyle 为 true，则展示为第二种方式
		// 该参数目前仅在查询参数中启作用
		ArrayStyle *BoolAttribute `apidoc:"array-style,attr,usage-param-array-style,omitempty"`
	}

	// Path 路径信息
	Path struct {
		token.BaseTag
		RootName struct{} `apidoc:"path,meta,usage-path"`

		Path    *Attribute `apidoc:"path,attr,usage-path-path"`
		Params  []*Param   `apidoc:"param,elem,usage-path-params,omitempty"`
		Queries []*Param   `apidoc:"query,elem,usage-path-queries,omitempty"`
	}

	// Request 请求内容
	Request struct {
		token.BaseTag
		RootName struct{} `apidoc:"request,meta,usage-request"`

		XML
		// 一般无用，但是用于描述 XML 对象时，可以用来表示顶层元素的名称
		Name *Attribute `apidoc:"name,attr,usage-request-name,omitempty"`

		Type        *TypeAttribute    `apidoc:"type,attr,usage-request-type,omitempty"`
		Deprecated  *VersionAttribute `apidoc:"deprecated,attr,usage-request-deprecated,omitempty"`
		Enums       []*Enum           `apidoc:"enum,elem,usage-request-enums,omitempty"`
		Array       *BoolAttribute    `apidoc:"array,attr,usage-request-array,omitempty"`
		Items       []*Param          `apidoc:"param,elem,usage-request-items,omitempty"`
		Summary     *Attribute        `apidoc:"summary,attr,usage-request-summary,omitempty"`
		Status      *StatusAttribute  `apidoc:"status,attr,usage-request-status,omitempty"`
		Mimetype    *Attribute        `apidoc:"mimetype,attr,usage-request-mimetype,omitempty"`
		Examples    []*Example        `apidoc:"example,elem,usage-request-examples,omitempty"`
		Headers     []*Param          `apidoc:"header,elem,usage-request-headers,omitempty"` // 当前独有的报头，公用的可以放在 API 中
		Description *Richtext         `apidoc:"description,elem,usage-request-description,omitempty"`
	}

	// Richtext 富文本内容
	Richtext struct {
		token.BaseTag
		RootName struct{} `apidoc:"richtext,meta,usage-richtext"`

		Type *Attribute `apidoc:"type,attr,usage-richtext-type"` // 文档类型，可以是 html 或是 markdown
		Text *CData     `apidoc:",cdata,usage-richtext-text"`
	}

	// Tag 标签内容
	Tag struct {
		token.BaseTag
		RootName struct{} `apidoc:"tag,meta,usage-tag"`

		Name       *Attribute        `apidoc:"name,attr,usage-tag-name"`   // 标签的唯一 ID
		Title      *Attribute        `apidoc:"title,attr,usage-tag-title"` // 显示的名称
		Deprecated *VersionAttribute `apidoc:"deprecated,attr,usage-tag-deprecated,omitempty"`
	}

	// Server 服务信息
	Server struct {
		token.BaseTag
		RootName struct{} `apidoc:"server,meta,usage-server"`

		Name        *Attribute        `apidoc:"name,attr,usage-server-name"` // 字面名称，需要唯一
		URL         *Attribute        `apidoc:"url,attr,usage-server-url"`
		Deprecated  *VersionAttribute `apidoc:"deprecated,attr,usage-server-deprecated,omitempty"`
		Summary     *Attribute        `apidoc:"summary,attr,usage-server-summary,omitempty"`
		Description *Richtext         `apidoc:"description,elem,usage-server-description,omitempty"`
	}

	// XML 仅作用于 XML 的几个属性
	XML struct {
		XMLAttr     *BoolAttribute `apidoc:"xml-attr,attr,usage-xml-attr,omitempty"`        // 作为父元素的 XML 属性存在
		XMLExtract  *BoolAttribute `apidoc:"xml-extract,attr,usage-xml-extract,omitempty"`  // 提取当前内容作为父元素的内容
		XMLCData    *BoolAttribute `apidoc:"xml-cdata,attr,usage-xml-cdata,omitempty"`      // 内容为 CDATA
		XMLNSPrefix *Attribute     `apidoc:"xml-ns-prefix,attr,usage-xml-prefix,omitempty"` // 命名空间前缀
		XMLWrapped  *Attribute     `apidoc:"xml-wrapped,attr,usage-xml-wrapped,omitempty"`  // 如果当前元素是数组，是否将其包含在 wrapped 中
	}
)

// V 返回当前富文本中的内容
func (r *Richtext) V() string {
	if r == nil || r.Text == nil {
		return ""
	}
	return r.Text.Value.Value
}

// Param 转换成 Param 对象
//
// Request 可以说是 Param 的超级，两者在大部分情况下能用。
func (r *Request) Param() *Param {
	if r == nil {
		return nil
	}

	return &Param{
		XML:         r.XML,
		Name:        r.Name,
		Type:        r.Type,
		Deprecated:  r.Deprecated,
		Optional:    &BoolAttribute{Value: Bool{Value: true}},
		Array:       r.Array,
		Items:       r.Items,
		Summary:     r.Summary,
		Enums:       r.Enums,
		Description: r.Description,
	}
}

// DeleteURI 删除与 uri 相关的文档内容
func (doc *APIDoc) DeleteURI(uri core.URI) {
	for index, api := range doc.APIs {
		if api.URI == uri {
			doc.APIs = append(doc.APIs[:index], doc.APIs[index+1:]...)
		}
	}

	if doc.URI == uri {
		*doc = APIDoc{
			APIs: doc.APIs,
		}
	}
}

// XMLNamespace 获取指定前缀名称的命名空间
func (doc *APIDoc) XMLNamespace(prefix string) *XMLNamespace {
	for _, ns := range doc.XMLNamespaces {
		if ns.Prefix.V() == prefix {
			return ns
		}
	}
	return nil
}
