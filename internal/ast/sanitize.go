// SPDX-License-Identifier: MIT

package ast

import (
	"sort"
	"strconv"

	"github.com/issue9/is"

	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/internal/locale"
	"github.com/caixw/apidoc/v7/internal/token"
)

// Sanitize token.Sanitizer
func (api *API) Sanitize(p *token.Parser) error {
	for _, header := range api.Headers { // 报头不能为 object
		if header.Type.V() == TypeObject {
			return p.NewError(header.Type.Start, header.Type.End, "header", locale.ErrInvalidValue)
		}
	}
	return nil
}

// Sanitize token.Sanitizer
func (e *Enum) Sanitize(p *token.Parser) error {
	if e.Description.V() == "" && e.Summary.V() == "" {
		return p.NewError(e.Start, e.End, "summary", locale.ErrRequired)
	}
	return nil
}

// Sanitize token.Sanitizer
func (p *Path) Sanitize(pp *token.Parser) error {
	if p.Path == nil || p.Path.V() == "" {
		return pp.NewError(p.Start, p.End, "path", locale.ErrRequired)
	}

	params, err := parsePath(p.Path.V())
	if err != nil {
		return pp.NewError(p.Path.Start, p.Path.End, "path", locale.ErrInvalidFormat)
	}
	if len(params) != len(p.Params) {
		return pp.NewError(p.Start, p.End, "path", locale.ErrPathNotMatchParams)
	}
	for _, param := range p.Params {
		if _, found := params[param.Name.V()]; !found {
			return pp.NewError(param.Start, param.End, "path", locale.ErrPathNotMatchParams)
		}
	}

	// 路径参数和查询参数不能为 object
	for _, item := range p.Params {
		if item.Type.V() == TypeObject {
			return pp.NewError(item.Start, item.End, "type", locale.ErrInvalidValue)
		}
	}
	for _, q := range p.Queries {
		if q.Type.V() == TypeObject {
			return pp.NewError(q.Start, q.End, "type", locale.ErrInvalidValue)
		}
	}

	return nil
}

func parsePath(path string) (params map[string]struct{}, err error) {
	start := -1
	for i, b := range path {
		switch b {
		case '{':
			if start != -1 {
				return nil, locale.NewError(locale.ErrInvalidFormat)
			}

			start = i + 1
		case '}':
			if start == -1 {
				return nil, locale.NewError(locale.ErrInvalidFormat)
			}

			if params == nil {
				params = make(map[string]struct{}, 3)
			}
			params[path[start:i]] = struct{}{}
			start = -1
		default:
		}
	}

	if start != -1 { // 没有结束符号
		return nil, locale.NewError(locale.ErrInvalidFormat)
	}

	return params, nil
}

// Sanitize token.Sanitizer
func (r *Request) Sanitize(p *token.Parser) error {
	if r.Type.V() == TypeObject && len(r.Items) == 0 {
		return p.NewError(r.Start, r.End, "param", locale.ErrRequired)
	}
	if r.Type.V() == TypeNone && len(r.Items) > 0 {
		return p.NewError(r.Start, r.End, "type", locale.ErrInvalidValue)
	}

	// 判断 enums 的值是否相同
	if rng, found := getDuplicateEnum(r.Enums); found {
		return p.NewError(rng.Start, rng.End, "enum", locale.ErrDuplicateValue)
	}

	if err := chkEnumsType(r.Type, r.Enums, p); err != nil {
		return err
	}

	if err := checkXML(r.Array.V(), len(r.Items) > 0, &r.XML, p); err != nil {
		return err
	}

	if r.Mimetype.V() != "" {
		for _, exp := range r.Examples {
			if exp.Mimetype.V() != r.Mimetype.V() {
				return p.NewError(r.Mimetype.Start, r.Mimetype.End, "mimetype", locale.ErrInvalidValue)
			}
		}
	}

	// 报头不能为 object
	for _, header := range r.Headers {
		if header.Type.V() == TypeObject {
			return p.NewError(header.Type.Start, header.Type.End, "type", locale.ErrInvalidValue)
		}
	}

	// 判断 items 的值是否相同
	if rng, found := getDuplicateItems(r.Items); found {
		return p.NewError(rng.Start, rng.End, "param", locale.ErrDuplicateValue)
	}

	return nil
}

// Sanitize token.Sanitizer
func (p *Param) Sanitize(pp *token.Parser) error {
	if p.Type.V() == TypeNone {
		return pp.NewError(p.Start, p.End, "type", locale.ErrRequired)
	}
	if p.Type.V() == TypeObject && len(p.Items) == 0 {
		return pp.NewError(p.Start, p.End, "param", locale.ErrRequired)
	}

	if p.Type.V() != TypeObject && len(p.Items) > 0 {
		return pp.NewError(p.Type.Value.Start, p.Type.Value.End, "type", locale.ErrInvalidValue)
	}

	// 判断 enums 的值是否相同
	if r, found := getDuplicateEnum(p.Enums); found {
		return pp.NewError(r.Start, r.End, "enum", locale.ErrDuplicateValue)
	}

	if err := chkEnumsType(p.Type, p.Enums, pp); err != nil {
		return err
	}

	// 判断 items 的值是否相同
	if r, found := getDuplicateItems(p.Items); found {
		return pp.NewError(r.Start, r.End, "param", locale.ErrDuplicateValue)
	}

	if err := checkXML(p.Array.V(), len(p.Items) > 0, &p.XML, pp); err != nil {
		return err
	}

	if p.Summary.V() == "" && p.Description.V() == "" {
		return pp.NewError(p.Start, p.End, "summary", locale.ErrRequired)
	}

	return nil
}

// 检测 enums 中的类型是否符合 t 的标准，比如 Number 要求枚举值也都是数值
func chkEnumsType(t *TypeAttribute, enums []*Enum, p *token.Parser) error {
	if len(enums) == 0 {
		return nil
	}

	switch t.V() {
	case TypeNumber:
		for _, enum := range enums {
			if !is.Number(enum.Value.V()) {
				return p.NewError(enum.Start, enum.End, enum.StartTag.String(), locale.ErrInvalidFormat)
			}
		}
	case TypeBool:
		for _, enum := range enums {
			if _, err := strconv.ParseBool(enum.Value.V()); err != nil {
				return p.NewError(enum.Start, enum.End, enum.StartTag.String(), locale.ErrInvalidFormat)
			}
		}
	case TypeObject, TypeNone:
		return p.NewError(t.Start, t.End, t.AttributeName.String(), locale.ErrInvalidValue)
	}

	return nil
}

// 返回重复枚举的值
func getDuplicateEnum(enums []*Enum) (core.Range, bool) {
	if len(enums) == 0 {
		return core.Range{}, false
	}

	es := make([]*Enum, 0, len(enums))
	for _, e := range enums {
		es = append(es, e)
	}
	sort.SliceStable(es, func(i, j int) bool { return es[i].Value.V() > es[j].Value.V() })

	for i := 1; i < len(es); i++ {
		if es[i].Value.V() == es[i-1].Value.V() {
			return es[i].Range, true
		}
	}

	return core.Range{}, false
}

func getDuplicateItems(items []*Param) (core.Range, bool) {
	if len(items) == 0 {
		return core.Range{}, false
	}

	es := make([]*Param, 0, len(items))
	for _, e := range items {
		es = append(es, e)
	}
	sort.SliceStable(es, func(i, j int) bool { return es[i].Name.V() > es[j].Name.V() })

	for i := 1; i < len(es); i++ {
		if es[i].Name.V() == es[i-1].Name.V() {
			return es[i].Range, true
		}
	}

	return core.Range{}, false
}

func checkXML(isArray, hasItems bool, xml *XML, p *token.Parser) error {
	if xml.XMLAttr.V() {
		if isArray || hasItems {
			return p.NewError(xml.XMLAttr.Start, xml.XMLAttr.End, xml.XMLAttr.AttributeName.String(), locale.ErrInvalidValue)
		}

		if xml.XMLWrapped.V() != "" {
			return p.NewError(xml.XMLWrapped.Start, xml.XMLWrapped.End, xml.XMLWrapped.AttributeName.String(), locale.ErrInvalidValue)
		}

		if xml.XMLExtract.V() {
			return p.NewError(xml.XMLExtract.Start, xml.XMLExtract.End, xml.XMLExtract.AttributeName.String(), locale.ErrInvalidValue)
		}

		if xml.XMLNSPrefix.V() != "" {
			return p.NewError(xml.XMLNSPrefix.Start, xml.XMLNSPrefix.End, xml.XMLNSPrefix.AttributeName.String(), locale.ErrInvalidValue)
		}

		if xml.XMLCData.V() {
			return p.NewError(xml.XMLCData.Start, xml.XMLCData.End, xml.XMLCData.AttributeName.String(), locale.ErrInvalidValue)
		}
	}

	if xml.XMLWrapped.V() != "" && !isArray {
		return p.NewError(xml.XMLWrapped.Start, xml.XMLWrapped.End, xml.XMLWrapped.AttributeName.String(), locale.ErrInvalidValue)
	}

	if xml.XMLExtract.V() {
		if xml.XMLNSPrefix.V() != "" {
			return p.NewError(xml.XMLNSPrefix.Start, xml.XMLNSPrefix.End, xml.XMLNSPrefix.AttributeName.String(), locale.ErrInvalidValue)
		}
	}

	return nil
}

// Sanitize 检测内容是否合法
func (doc *APIDoc) Sanitize(p *token.Parser) error {
	// 按 URN 查重
	sort.SliceStable(doc.XMLNamespaces, func(i, j int) bool {
		return doc.XMLNamespaces[i].URN.V() > doc.XMLNamespaces[j].URN.V()
	})
	for i := 1; i < len(doc.XMLNamespaces); i++ {
		curr := doc.XMLNamespaces[i].URN
		if doc.XMLNamespaces[i-1].URN.V() == curr.V() {
			return p.NewError(curr.Start, curr.End, "@urn", locale.ErrDuplicateValue)
		}
	}

	// 按 prefix 查重
	sort.SliceStable(doc.XMLNamespaces, func(i, j int) bool {
		return doc.XMLNamespaces[i].Prefix.V() > doc.XMLNamespaces[j].Prefix.V()
	})
	var auto bool
	for i := 1; i < len(doc.XMLNamespaces); i++ {
		curr := doc.XMLNamespaces[i]
		if doc.XMLNamespaces[i-1].Prefix.V() == curr.Prefix.V() {
			return p.NewError(curr.Start, curr.End, "@prefix", locale.ErrDuplicateValue)
		}
		if curr.Auto.V() {
			if auto {
				return p.NewError(curr.Start, curr.End, "@auto", locale.ErrInvalidValue)
			}
			auto = true
		}
	}

	for _, api := range doc.APIs {
		if api.doc == nil {
			api.doc = doc // 保证单文件的文档能正常解析
			api.URI = doc.URI
		}
		if err := api.sanitizeTags(); err != nil {
			return err
		}
	}

	return nil
}

func (doc *APIDoc) sortAPIs() {
	sort.SliceStable(doc.APIs, func(i, j int) bool {
		ii := doc.APIs[i]
		jj := doc.APIs[j]

		var iip string
		if ii.Path != nil && ii.Path.Path != nil {
			iip = ii.Path.Path.V()
		}

		var jjp string
		if jj.Path != nil && jj.Path.Path != nil {
			jjp = jj.Path.Path.V()
		}

		var iim string
		if ii.Method != nil {
			iim = ii.Method.V()
		}

		var jjm string
		if jj.Method != nil {
			jjm = jj.Method.V()
		}

		if iip == jjp {
			return iim < jjm
		}
		return iip < jjp
	})
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

func (api *API) sanitizeTags() error {
	if api.doc == nil {
		panic("api.doc 未获取正确的值")
	}

	for _, tag := range api.Tags {
		if !api.doc.tagExists(tag.Content.Value) {
			loc := core.Location{URI: api.URI,
				Range: core.Range{
					Start: tag.Content.Start,
					End:   tag.Content.End,
				}}
			return core.NewSyntaxError(loc, "content", locale.ErrInvalidValue)
		}
	}

	for _, srv := range api.Servers {
		if !api.doc.serverExists(srv.Content.Value) {
			loc := core.Location{URI: api.URI,
				Range: core.Range{
					Start: srv.Content.Start,
					End:   srv.Content.End,
				}}
			return core.NewSyntaxError(loc, "content", locale.ErrInvalidValue)
		}
	}

	return nil
}
