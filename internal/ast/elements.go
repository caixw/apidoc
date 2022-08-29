// SPDX-License-Identifier: MIT

package ast

import "github.com/caixw/apidoc/v7/internal/xmlenc"

type (
	// APIDoc 对应 apidoc 元素
	APIDoc struct {
		xmlenc.BaseTag
		RootName struct{} `apidoc:"apidoc,meta,usage-apidoc"`

		APIDoc        *APIDocVersionAttribute `apidoc:"apidoc,attr,usage-apidoc-apidoc,omitempty"` // 文档格式的版本号
		Lang          *Attribute              `apidoc:"lang,attr,usage-apidoc-lang,omitempty"`     // 区域信息，应该使用 BCP47 指定的格式
		Logo          *Attribute              `apidoc:"logo,attr,usage-apidoc-logo,omitempty"`
		XMLNamespaces []*XMLNamespace         `apidoc:"xml-namespace,elem,usage-apidoc-xml-namespaces,omitempty"`
		Created       *DateAttribute          `apidoc:"created,attr,usage-apidoc-created,omitempty"` // 生成时间
		Version       *VersionAttribute       `apidoc:"version,attr,usage-apidoc-version,omitempty"`
		Title         *Element                `apidoc:"title,elem,usage-apidoc-title"`
		Description   *Richtext               `apidoc:"description,elem,usage-apidoc-description,omitempty"` // 说明内容
		Contact       *Contact                `apidoc:"contact,elem,usage-apidoc-contact,omitempty"`         // 团队的联系方式
		License       *Link                   `apidoc:"license,elem,usage-apidoc-license,omitempty"`         // 版权信息
		Tags          []*Tag                  `apidoc:"tag,elem,usage-apidoc-tags,omitempty"`                // 标签列表
		Servers       []*Server               `apidoc:"server,elem,usage-apidoc-servers,omitempty"`          // 服务器列表
		APIs          []*API                  `apidoc:"api,elem,usage-apidoc-apis,omitempty"`                // API 列表
		Headers       []*Param                `apidoc:"header,elem,usage-apidoc-headers,omitempty"`          // 公共报头
		Responses     []*Request              `apidoc:"response,elem,usage-apidoc-responses,omitempty"`      // 所有 API 都有可能的返回内容
		Mimetypes     []*Element              `apidoc:"mimetype,elem,usage-apidoc-mimetypes"`                // 所有接口都支持的 mimetypes
	}

	// XMLNamespace 定义命名空间的相关属性
	XMLNamespace struct {
		xmlenc.BaseTag
		RootName struct{}   `apidoc:"xml-namespace,meta,usage-xml-namespace"`
		Prefix   *Attribute `apidoc:"prefix,attr,usage-xml-namespace-prefix,omitempty"`
		URN      *Attribute `apidoc:"urn,attr,usage-xml-namespace-urn"`
	}

	// API 表示 <api> 顶层元素
	API struct {
		xmlenc.BaseTag
		RootName struct{} `apidoc:"api,meta,usage-api"`
		doc      *APIDoc

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
		Tags        []*TagValue       `apidoc:"tag,elem,usage-api-tags,omitempty"`
		Servers     []*ServerValue    `apidoc:"server,elem,usage-api-servers,omitempty"`
	}

	// Link 表示一个链接
	Link struct {
		xmlenc.BaseTag
		RootName struct{} `apidoc:"link,meta,usage-link"`

		Text *Attribute `apidoc:"text,attr,usage-link-text"`
		URL  *Attribute `apidoc:"url,attr,usage-link-url"`
	}

	// Contact 描述联系方式
	Contact struct {
		xmlenc.BaseTag
		RootName struct{} `apidoc:"contact,meta,usage-contact"`

		Name  *Attribute `apidoc:"name,attr,usage-contact-name"`
		URL   *Element   `apidoc:"url,elem,usage-contact-url,omitempty"`
		Email *Element   `apidoc:"email,elem,usage-contact-email,omitempty"`
	}

	// Callback 描述回调信息
	Callback struct {
		xmlenc.BaseTag
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
		xmlenc.BaseTag
		RootName struct{} `apidoc:"enum,meta,usage-enum"`

		Deprecated  *VersionAttribute `apidoc:"deprecated,attr,usage-enum-deprecated,omitempty"`
		Value       *Attribute        `apidoc:"value,attr,usage-enum-value"`
		Summary     *Attribute        `apidoc:"summary,attr,usage-enum-summary,omitempty"`
		Description *Richtext         `apidoc:"description,elem,usage-enum-description,omitempty"`
	}

	// Example 示例代码
	Example struct {
		xmlenc.BaseTag
		RootName struct{} `apidoc:"example,meta,usage-example"`

		Mimetype *Attribute    `apidoc:"mimetype,attr,usage-example-mimetype"`
		Content  *ExampleValue `apidoc:",cdata,usage-example-content"`
		Summary  *Attribute    `apidoc:"summary,attr,usage-example-summary,omitempty"`
	}

	// Param 表示参数类型
	Param struct {
		xmlenc.BaseTag
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
		xmlenc.BaseTag
		RootName struct{} `apidoc:"path,meta,usage-path"`

		Path    *Attribute `apidoc:"path,attr,usage-path-path"`
		Params  []*Param   `apidoc:"param,elem,usage-path-params,omitempty"`
		Queries []*Param   `apidoc:"query,elem,usage-path-queries,omitempty"`
	}

	// Request 请求内容
	Request struct {
		xmlenc.BaseTag
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
		xmlenc.BaseTag
		RootName struct{} `apidoc:"richtext,meta,usage-richtext"`

		Type *Attribute `apidoc:"type,attr,usage-richtext-type"` // 文档类型，可以是 html 或是 markdown
		Text *CData     `apidoc:",cdata,usage-richtext-text"`
	}

	// Tag 标签内容
	Tag struct {
		xmlenc.BaseTag
		RootName struct{} `apidoc:"tag,meta,usage-tag"`

		Name       *Attribute        `apidoc:"name,attr,usage-tag-name"`   // 标签的唯一 ID
		Title      *Attribute        `apidoc:"title,attr,usage-tag-title"` // 显示的名称
		Deprecated *VersionAttribute `apidoc:"deprecated,attr,usage-tag-deprecated,omitempty"`

		references []*Reference
	}

	// Server 服务信息
	Server struct {
		xmlenc.BaseTag
		RootName struct{} `apidoc:"server,meta,usage-server"`

		Name        *Attribute        `apidoc:"name,attr,usage-server-name"` // 字面名称，需要唯一
		URL         *Attribute        `apidoc:"url,attr,usage-server-url"`
		Deprecated  *VersionAttribute `apidoc:"deprecated,attr,usage-server-deprecated,omitempty"`
		Summary     *Attribute        `apidoc:"summary,attr,usage-server-summary,omitempty"`
		Description *Richtext         `apidoc:"description,elem,usage-server-description,omitempty"`

		references []*Reference
	}

	// XML 仅作用于 XML 的几个属性
	XML struct {
		XMLAttr     *BoolAttribute `apidoc:"xml-attr,attr,usage-xml-attr,omitempty"`        // 作为父元素的 XML 属性存在
		XMLExtract  *BoolAttribute `apidoc:"xml-extract,attr,usage-xml-extract,omitempty"`  // 提取当前内容作为父元素的内容
		XMLCData    *BoolAttribute `apidoc:"xml-cdata,attr,usage-xml-cdata,omitempty"`      // 内容为 CDATA
		XMLNSPrefix *Attribute     `apidoc:"xml-ns-prefix,attr,usage-xml-prefix,omitempty"` // 命名空间前缀
		XMLWrapped  *Attribute     `apidoc:"xml-wrapped,attr,usage-xml-wrapped,omitempty"`  // 如果当前元素是数组，是否将其包含在 wrapped 中
	}

	// Element 定义不包含子元素和属性的基本的 XML 元素
	Element struct {
		xmlenc.BaseTag
		Content  Content  `apidoc:",content"`
		RootName struct{} `apidoc:"string,meta,usage-string"`
	}

	// TagValue api.tag 的类型
	TagValue struct {
		xmlenc.BaseTag
		Content  Content  `apidoc:",content"`
		RootName struct{} `apidoc:"string,meta,usage-string"`

		definition *Definition
	}

	// ServerValue api.server 的类型
	ServerValue struct {
		xmlenc.BaseTag
		Content  Content  `apidoc:",content"`
		RootName struct{} `apidoc:"string,meta,usage-string"`

		definition *Definition
	}

	// CData 表示 XML 的 CDATA 数据
	CData struct {
		xmlenc.BaseTag
		Value    xmlenc.String `apidoc:"-"`
		RootName struct{}      `apidoc:"string,meta,usage-string"`
	}

	// ExampleValue 示例代码的内容
	ExampleValue CData

	// Content 表示一段普通的 XML 元素内容
	Content struct {
		xmlenc.Base
		Value    string   `apidoc:"-"`
		RootName struct{} `apidoc:"string,meta,usage-string"`
	}
)

// V 返回当前富文本中的内容
func (r *Richtext) V() string {
	if r == nil || r.Text == nil {
		return ""
	}
	return r.Text.Value.Value
}

// V 返回当前属性实际表示的值
func (s *Element) V() string {
	if s == nil {
		return ""
	}
	return s.Content.Value
}

// V 返回当前属性实际表示的值
func (s *TagValue) V() string {
	if s == nil {
		return ""
	}
	return s.Content.Value
}

// V 返回当前属性实际表示的值
func (s *ServerValue) V() string {
	if s == nil {
		return ""
	}
	return s.Content.Value
}

// Definition Definitioner.Definition
func (s *TagValue) Definition() *Definition {
	return s.definition
}

// Definition Definitioner.Definition
func (s *ServerValue) Definition() *Definition {
	return s.definition
}

// EncodeXML Encoder.EncodeXML
func (cdata *CData) EncodeXML() (string, error) {
	return cdata.Value.Value, nil
}

// EncodeXML Encoder.EncodeXML
func (s *Content) EncodeXML() (string, error) {
	return s.Value, nil
}

// EncodeXML Encoder.EncodeXML
//
// 示例代码的内容，会在此处去掉其前导的空格
func (v *ExampleValue) EncodeXML() (string, error) {
	return trimLeftSpace(v.Value.Value), nil
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

// XMLNamespace 获取指定前缀名称的命名空间
func (doc *APIDoc) XMLNamespace(prefix string) *XMLNamespace {
	for _, ns := range doc.XMLNamespaces {
		if ns.Prefix.V() == prefix {
			return ns
		}
	}
	return nil
}

// References impl Referencer
func (tag *Tag) References() []*Reference {
	return tag.references
}

// References impl Referencer
func (srv *Server) References() []*Reference {
	return srv.references
}
