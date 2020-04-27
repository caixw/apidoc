// SPDX-License-Identifier: MIT

package ast

import (
	"sort"

	"github.com/caixw/apidoc/v6/core"
	"github.com/caixw/apidoc/v6/internal/locale"
	"github.com/caixw/apidoc/v6/internal/token"
)

// ErrNoDocFormat 表示内容并不是一个文档格式
//
// 比如普通的注释内容等。
var ErrNoDocFormat = locale.Errorf(locale.ErrIsNotAPIDoc)

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
		token.Base

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

		Created     *Attribute        `apidoc:"created,attr,usage-apidoc-created,omitempty"` // 文档的生成时间
		Version     *VersionAttribute `apidoc:"version,attr,usage-apidoc-version,omitempty"` // 文档的版本
		Title       *Element          `apidoc:"title,elem,usage-apidoc-title"`
		Description *Richtext         `apidoc:"description,elem,usage-apidoc-description,omitempty"`
		Contact     *Contact          `apidoc:"contact,elem,usage-apidoc-contact,omitempty"`
		License     *Link             `apidoc:"license,elem,usage-apidoc-license,omitempty"` // 版本信息
		Tags        []*Tag            `apidoc:"tag,elem,usage-apidoc-tags,omitempty"`        // 所有的标签
		Servers     []*Server         `apidoc:"server,elem,usage-apidoc-servers,omitempty"`
		Apis        []*API            `apidoc:"api,elem,usage-apidoc-apis,omitempty"`

		// 表示所有 API 都有可能返回的内容
		Responses []*Request `apidoc:"response,elem,usage-apidoc-responses,omitempty"`

		// 表示所有接口都支持的文档类型
		Mimetypes []*Element `apidoc:"mimetype,elem,usage-apidoc-mimetypes"`

		Block *core.Block `apidoc:"-"`
	}

	// API 表示 <api> 顶层元素
	API struct {
		token.Base
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

		Block *core.Block `apidoc:"-"`
		doc   *APIDoc
	}

	// Link 表示一个链接
	Link struct {
		token.Base
		Text *Attribute `apidoc:"text,attr,usage-link-text"`
		URL  *Attribute `apidoc:"url,attr,usage-link-url"`
	}

	// Contact 描述联系方式
	Contact struct {
		token.Base
		Name  *Attribute `apidoc:"name,attr,usage-contact-name"`
		URL   *Element   `apidoc:"url,elem,usage-contact-url,omitempty"`
		Email *Element   `apidoc:"email,elem,usage-contact-email,omitempty"`
	}

	// Callback 描述回调信息
	Callback struct {
		token.Base
		Method      *MethodAttribute  `apidoc:"method,attr,usage-callback-method"`
		Path        *Path             `apidoc:"path,elem,usage-callback-path,omitempty"`
		Summary     *Attribute        `apidoc:"summary,attr,usage-callback-summary,omitempty"`
		Description *Richtext         `apidoc:"description,elem,usage-callback-description,omitempty"`
		Deprecated  *VersionAttribute `apidoc:"deprecated,attr,usage-callback-deprecated,omitempty"`
		Reference   *Attribute        `apidoc:"ref,attr,usage-callback-reference,omitempty"`
		Responses   []*Request        `apidoc:"response,elem,usage-callback-responses,omitempty"`
		Requests    []*Request        `apidoc:"request,elem,usage-callback-requests"` // 至少一个
		Headers     []*Param          `apidoc:"header,elem,usage-callback-headers,omitempty"`
	}

	// Enum 表示枚举值
	Enum struct {
		token.Base
		Deprecated  *VersionAttribute `apidoc:"deprecated,attr,usage-enum-deprecated,omitempty"`
		Value       *Attribute        `apidoc:"value,attr,usage-enum-value"`
		Summary     *Attribute        `apidoc:"summary,attr,usage-enum-summary,omitempty"`
		Description *Richtext         `apidoc:"description,elem,usage-enum-description,omitempty"`
	}

	// Example 示例代码
	Example struct {
		token.Base
		Mimetype *Attribute `apidoc:"mimetype,attr,usage-example-mimetype"`
		Content  *CData     `apidoc:",cdata"`
		Summary  *Attribute `apidoc:"summary,attr,usage-example-summary,omitempty"`
	}

	// Param 表示参数类型
	Param struct {
		token.Base
		XML
		Name        *Attribute        `apidoc:"name,attr,usage-param-name"`
		Type        *TypeAttribute    `apidoc:"type,attr,usage-param-type"`
		Deprecated  *VersionAttribute `apidoc:"deprecated,attr,usage-param-deprecated,omitempty"`
		Default     *Attribute        `apidoc:"default,attr,usage-param-default,omitempty"`
		Optional    *BoolAttribute    `apidoc:"optional,attr,usage-param-optional,omitempty"`
		Array       *BoolAttribute    `apidoc:"array,attr,usage-parm-array,omitempty"`
		Items       []*Param          `apidoc:"param,elem,usage-param-items,omitempty"`
		Reference   *Attribute        `apidoc:"ref,attr,usage-param-reference,omitempty"`
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
		token.Base
		Path      *Attribute `apidoc:"path,attr,usage-path-path"`
		Reference *Attribute `apidoc:"ref,attr,usage-path-reference,omitempty"`
		Params    []*Param   `apidoc:"param,elem,usage-path-params,omitempty"`
		Queries   []*Param   `apidoc:"query,elem,usage-path-queries,omitempty"`
	}

	// Request 请求内容
	Request struct {
		token.Base
		XML

		// 一般无用，但是用于描述 XML 对象时，可以用来表示顶层元素的名称
		Name *Attribute `apidoc:"name,attr,usage-request-name,omitempty"`

		Type        *TypeAttribute    `apidoc:"type,attr,usage-request-type,omitempty"`
		Deprecated  *VersionAttribute `apidoc:"deprecated,attr,usage-request-deprecated,omitempty"`
		Enums       []*Enum           `apidoc:"enum,elem,usage-request-enums,omitempty"`
		Array       *BoolAttribute    `apidoc:"array,attr,usage-request-array,omitempty"`
		Items       []*Param          `apidoc:"param,elem,usage-request-items,omitempty"`
		Reference   *Attribute        `apidoc:"ref,attr,usage-request-reference,omitempty"`
		Summary     *Attribute        `apidoc:"summary,attr,usage-request-summary,omitempty"`
		Status      *StatusAttribute  `apidoc:"status,attr,usage-request-status,omitempty"`
		Mimetype    *Attribute        `apidoc:"mimetype,attr,usage-request-mimetype,omitempty"`
		Examples    []*Example        `apidoc:"example,elem,usage-request-examples,omitempty"`
		Headers     []*Param          `apidoc:"header,elem,usage-request-headers,omitempty"` // 当前独有的报头，公用的可以放在 API 中
		Description *Richtext         `apidoc:"description,elem,usage-request-description,omitempty"`
	}

	// Richtext 富文本内容
	Richtext struct {
		token.Base
		Type *Attribute `apidoc:"type,attr,usage-richtext-type,omitempty"` // 文档类型，可以是 html 或是 markdown
		Text *CData     `apidoc:",cdata"`
	}

	// Tag 标签内容
	Tag struct {
		token.Base
		Name       *Attribute        `apidoc:"name,attr,usage-tag-name"`   // 标签的唯一 ID
		Title      *Attribute        `apidoc:"title,attr,usage-tag-title"` // 显示的名称
		Deprecated *VersionAttribute `apidoc:"deprecated,attr,usage-tag-deprecated,omitempty"`
	}

	// Server 服务信息
	Server struct {
		token.Base
		Name        *Attribute        `apidoc:"name,attr,usage-server-name"` // 字面名称，需要唯一
		URL         *Attribute        `apidoc:"url,attr,usage-server-url"`
		Deprecated  *VersionAttribute `apidoc:"deprecated,attr,usage-server-deprecated,omitempty"`
		Summary     *Attribute        `apidoc:"summary,attr,usage-server-summary,omitempty"`
		Description *Richtext         `apidoc:"description,elem,usage-server-description,omitempty"`
	}

	// XML 仅作用于 XML 的几个属性
	XML struct {
		XMLAttr     *BoolAttribute `apidoc:"xml-attr,attr,usage-xml-xml-attr,omitempty"`        // 作为父元素的 XML 属性存在
		XMLExtract  *BoolAttribute `apidoc:"xml-extract,attr,usage-xml-xml-extract,omitempty"`  // 提取当前内容作为父元素的内容
		XMLNS       *Attribute     `apidoc:"xml-ns,attr,usage-xml-xml-ns,omitempty"`            // 命名空间
		XMLNSPrefix *Attribute     `apidoc:"xml-ns-prefix,attr,usage-xml-xml-prefix,omitempty"` // 命名空间前缀
		XMLWrapped  *Attribute     `apidoc:"xml-wrapped,attr,usage-xml-xml-wrapped,omitempty"`  // 如果当前元素是数组，是否将其包含在 wrapped 中
	}
)

// Parse 将注释块的内容添加到当前文档
//
// 分析注释块内容，如果正确，则添加到当前文档中，
// 或是在出错时，返回错误信息。
//
// 如果内容不是文档内容，刚将返回 ErrNoDocFormat
func (doc *APIDoc) Parse(b core.Block) error {
	p, err := token.NewParser(b)
	if err != nil {
		return err
	}

	name, err := getTagName(p)
	if err != nil {
		return err
	}
	switch name {
	case "api":
		api := &API{doc: doc}
		if doc.Apis == nil {
			doc.Apis = make([]*API, 0, 100)
		}
		doc.Apis = append(doc.Apis, api)
		return token.Decode(p, api)
	case "apidoc":
		if doc.Title != nil { // 多个 apidoc 标签
			return core.NewLocaleError(b.Location, "", locale.ErrDuplicateValue)
		}
		return token.Decode(p, doc)
	default:
		return core.NewLocaleError(b.Location, "", locale.ErrInvalidXML)
	}
}

// 获取根标签的名称
func getTagName(p *token.Parser) (string, error) {
	start := p.Position()
	for {
		t, err := p.Token()
		if err != nil {
			return "", err
		}
		if t == nil {
			return "", nil
		}

		switch elem := t.(type) {
		case *token.StartElement:
			p.Move(start)
			return elem.Name.Value, nil
		case *token.EndElement, *token.CData:
			return "", ErrNoDocFormat
		default: // 其它标签忽略
		}
	}
}

// Sanitize 检测内容是否合法
func (doc *APIDoc) Sanitize() error {
	// doc.Apis 是多线程导入的，无法保证其顺序，
	// 此处可以保证输出内容是按一定顺序排列的。
	sort.SliceStable(doc.Apis, func(i, j int) bool {
		ii := doc.Apis[i]
		jj := doc.Apis[j]

		if ii.Path.Path == jj.Path.Path {
			return ii.Method.Value.Value < jj.Method.Value.Value
		}
		return ii.Path.Path.Value.Value < jj.Path.Path.Value.Value
	})

	for _, api := range doc.Apis { // 查看 API 中的标签是否都存在
		if err := api.sanitize(); err != nil {
			return err
		}
	}

	return nil
}

func (doc *APIDoc) tagExists(tag string) bool {
	for _, s := range doc.Tags {
		if s.Name.Value.Value == tag {
			return true
		}
	}
	return false
}

func (doc *APIDoc) serverExists(srv string) bool {
	for _, s := range doc.Servers {
		if s.Name.Value.Value == srv {
			return true
		}
	}
	return false
}

// 检测和修复 api 对象，无法修复返回错误。
//
// NOTE: 需要保证 doc 已经初始化
func (api *API) sanitize() error {
	if api.doc == nil {
		panic("api.doc 未获取正确的值")
	}

	for _, tag := range api.Tags {
		if !api.doc.tagExists(tag.Content.Value) {
			loc := core.Location{URI: api.Block.Location.URI, Range: tag.Content.Range}
			return core.NewLocaleError(loc, "", locale.ErrInvalidValue)
		}
	}

	if len(api.Servers) == 0 {
		return core.NewLocaleError(api.Block.Location, "", locale.ErrRequired)
	}

	for _, srv := range api.Servers {
		if !api.doc.serverExists(srv.Content.Value) {
			loc := core.Location{URI: api.Block.Location.URI, Range: srv.Content.Range}
			return core.NewLocaleError(loc, "", locale.ErrInvalidValue)
		}
	}

	return nil
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
		Reference:   r.Reference,
		Summary:     r.Summary,
		Enums:       r.Enums,
		Description: r.Description,
	}
}
