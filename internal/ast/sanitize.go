// SPDX-License-Identifier: MIT

package ast

import (
	"strconv"

	"github.com/issue9/is"
	"github.com/issue9/sliceutil"

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

	// 对 Servers 和 Tags 查重
	i := sliceutil.Dup(api.Servers, func(i, j int) bool { return api.Servers[i].V() == api.Servers[j].V() })
	if i > -1 {
		return p.NewError(api.Servers[i].Start, api.Servers[i].End, "server", locale.ErrDuplicateValue)
	}
	i = sliceutil.Dup(api.Tags, func(i, j int) bool { return api.Tags[i].V() == api.Tags[j].V() })
	if i > -1 {
		return p.NewError(api.Tags[i].Start, api.Tags[i].End, "server", locale.ErrDuplicateValue)
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
	i := sliceutil.Dup(enums, func(i, j int) bool { return enums[i].Value.V() == enums[j].Value.V() })
	if i > -1 {
		return enums[i].Range, true
	}
	return core.Range{}, false
}

func getDuplicateItems(items []*Param) (core.Range, bool) {
	i := sliceutil.Dup(items, func(i, j int) bool { return items[i].Name.V() == items[j].Name.V() })
	if i > -1 {
		return items[i].Range, true
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
	if err := doc.checkXMLNamespaces(p); err != nil {
		return err
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

// Sanitize 检测内容是否合法
func (ns *XMLNamespace) Sanitize(p *token.Parser) error {
	if ns.URN.V() == "" {
		return p.NewError(ns.Start, ns.End, "@urn", locale.ErrRequired)
	}
	return nil
}

func (doc *APIDoc) checkXMLNamespaces(p *token.Parser) error {
	if len(doc.XMLNamespaces) == 0 {
		return nil
	}

	// 按 URN 查重
	i := sliceutil.Dup(doc.XMLNamespaces, func(i, j int) bool {
		return doc.XMLNamespaces[i].URN.V() == doc.XMLNamespaces[j].URN.V()
	})
	if i > -1 {
		curr := doc.XMLNamespaces[i].URN
		return p.NewError(curr.Start, curr.End, "@urn", locale.ErrDuplicateValue)
	}

	// 按 prefix 查重
	i = sliceutil.Dup(doc.XMLNamespaces, func(i, j int) bool {
		return doc.XMLNamespaces[i].Prefix.V() == doc.XMLNamespaces[j].Prefix.V()
	})
	if i > -1 {
		curr := doc.XMLNamespaces[i].URN
		return p.NewError(curr.Start, curr.End, "@prefix", locale.ErrDuplicateValue)
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
			return core.NewSyntaxError(loc, "", locale.ErrInvalidValue)
		}
	}

	for _, srv := range api.Servers {
		if !api.doc.serverExists(srv.Content.Value) {
			loc := core.Location{URI: api.URI,
				Range: core.Range{
					Start: srv.Content.Start,
					End:   srv.Content.End,
				}}
			return core.NewSyntaxError(loc, "", locale.ErrInvalidValue)
		}
	}

	return nil
}
