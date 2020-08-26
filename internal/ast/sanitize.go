// SPDX-License-Identifier: MIT

package ast

import (
	"strconv"

	"github.com/issue9/is"
	"github.com/issue9/sliceutil"

	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/internal/locale"
	"github.com/caixw/apidoc/v7/internal/xmlenc"
)

// Sanitize token.Sanitizer
func (api *API) Sanitize(p *xmlenc.Parser) {
	for _, header := range api.Headers { // 报头不能为 object
		if header.Type.V() == TypeObject {
			p.Error(header.Type.Location.NewError(locale.ErrInvalidValue).WithField("header"))
		}
	}

	// 对 Servers 和 Tags 查重
	indexes := sliceutil.Dup(api.Servers, func(i, j int) bool { return api.Servers[i].V() == api.Servers[j].V() })
	if len(indexes) > 0 {
		err := api.Servers[indexes[0]].Location.NewError(locale.ErrDuplicateValue).WithField("server")
		for _, srv := range indexes[1:] {
			err.Relate(api.Servers[srv].Location, locale.Sprintf(locale.ErrDuplicateValue))
		}
		p.Error(err)
	}
	indexes = sliceutil.Dup(api.Tags, func(i, j int) bool { return api.Tags[i].V() == api.Tags[j].V() })
	if len(indexes) > 0 {
		err := api.Tags[indexes[0]].Location.NewError(locale.ErrDuplicateValue).WithField("server")
		for _, tag := range indexes[1:] {
			err.Relate(api.Tags[tag].Location, locale.Sprintf(locale.ErrDuplicateValue))
		}
		p.Error(err)
	}
}

// Sanitize token.Sanitizer
func (e *Enum) Sanitize(p *xmlenc.Parser) {
	if e.Description.V() == "" && e.Summary.V() == "" {
		p.Error(e.Location.NewError(locale.ErrIsEmpty, "summary").WithField("summary"))
	}
}

// Sanitize token.Sanitizer
func (p *Path) Sanitize(pp *xmlenc.Parser) {
	if p.Path == nil || p.Path.V() == "" {
		pp.Error(p.Location.NewError(locale.ErrIsEmpty, "path").WithField("path"))
	}

	params, err := parsePath(p.Path.V())
	if err != nil {
		pp.Error(p.Path.Location.NewError(locale.ErrInvalidFormat).WithField("path"))
	}
	if len(params) != len(p.Params) {
		pp.Error(p.Location.NewError(locale.ErrPathNotMatchParams).WithField("path"))
	}
	for _, param := range p.Params {
		if _, found := params[param.Name.V()]; !found {
			pp.Error(param.Location.NewError(locale.ErrPathNotMatchParams).WithField("path"))
		}
	}

	// 路径参数和查询参数不能为 object
	for _, item := range p.Params {
		if item.Type.V() == TypeObject {
			pp.Error(item.Location.NewError(locale.ErrInvalidValue).WithField("type"))
		}
	}
	for _, q := range p.Queries {
		if q.Type.V() == TypeObject {
			pp.Error(q.Location.NewError(locale.ErrInvalidValue).WithField("type"))
		}
	}
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
func (r *Request) Sanitize(p *xmlenc.Parser) {
	if r.Type.V() == TypeObject && len(r.Items) == 0 {
		p.Error(r.Location.NewError(locale.ErrIsEmpty, "param").WithField("param"))
	}
	if r.Type.V() == TypeNone && len(r.Items) > 0 {
		p.Error(r.Location.NewError(locale.ErrInvalidValue).WithField("type"))
	}

	checkDuplicateEnum(r.Enums, p)

	if err := chkEnumsType(r.Type, r.Enums, p); err != nil {
		p.Error(err)
	}

	if err := checkXML(r.Array.V(), len(r.Items) > 0, &r.XML, p); err != nil {
		p.Error(err)
	}

	if r.Mimetype.V() != "" {
		for _, exp := range r.Examples {
			if exp.Mimetype.V() != r.Mimetype.V() {
				p.Error(r.Mimetype.Location.NewError(locale.ErrInvalidValue).WithField("mimetype"))
			}
		}
	}

	// 报头不能为 object
	for _, header := range r.Headers {
		if header.Type.V() == TypeObject {
			p.Error(header.Type.Location.NewError(locale.ErrInvalidValue).WithField("type"))
		}
	}

	checkDuplicateItems(r.Items, p)
}

// Sanitize token.Sanitizer
func (p *Param) Sanitize(pp *xmlenc.Parser) {
	if p.Type.V() == TypeNone {
		pp.Error(p.Location.NewError(locale.ErrIsEmpty, "type").WithField("type"))
	}
	if p.Type.V() == TypeObject && len(p.Items) == 0 {
		pp.Error(p.Location.NewError(locale.ErrIsEmpty, "param").WithField("param"))
	}

	if p.Type.V() != TypeObject && len(p.Items) > 0 {
		pp.Error(p.Type.Value.Location.NewError(locale.ErrInvalidValue).WithField("type"))
	}

	checkDuplicateEnum(p.Enums, pp)

	if err := chkEnumsType(p.Type, p.Enums, pp); err != nil {
		pp.Error(err)
	}

	checkDuplicateItems(p.Items, pp)

	if err := checkXML(p.Array.V(), len(p.Items) > 0, &p.XML, pp); err != nil {
		pp.Error(err)
	}

	if p.Summary.V() == "" && p.Description.V() == "" {
		pp.Error(p.Location.NewError(locale.ErrIsEmpty, "summary").WithField("summary"))
	}
}

// 检测 enums 中的类型是否符合 t 的标准，比如 Number 要求枚举值也都是数值
func chkEnumsType(t *TypeAttribute, enums []*Enum, p *xmlenc.Parser) error {
	if len(enums) == 0 {
		return nil
	}

	switch t.V() {
	case TypeNumber:
		for _, enum := range enums {
			if !is.Number(enum.Value.V()) {
				return enum.Location.NewError(locale.ErrInvalidFormat).WithField(enum.StartTag.String())
			}
		}
	case TypeBool:
		for _, enum := range enums {
			if _, err := strconv.ParseBool(enum.Value.V()); err != nil {
				return enum.Location.NewError(locale.ErrInvalidFormat).WithField(enum.StartTag.String())
			}
		}
	case TypeObject, TypeNone:
		return t.Location.NewError(locale.ErrInvalidValue).WithField(t.AttributeName.String())
	}

	return nil
}

func checkDuplicateEnum(enums []*Enum, p *xmlenc.Parser) {
	indexes := sliceutil.Dup(enums, func(i, j int) bool { return enums[i].Value.V() == enums[j].Value.V() })
	if len(indexes) > 0 {
		err := enums[indexes[0]].Location.NewError(locale.ErrDuplicateValue).WithField("enum")
		for _, i := range indexes[1:] {
			err.Relate(enums[i].Location, locale.Sprintf(locale.ErrDuplicateValue))
		}
		p.Error(err)
	}
}

func checkDuplicateItems(items []*Param, p *xmlenc.Parser) {
	indexes := sliceutil.Dup(items, func(i, j int) bool { return items[i].Name.V() == items[j].Name.V() })
	if len(indexes) > 0 {
		err := items[indexes[0]].Location.NewError(locale.ErrDuplicateValue).WithField("param")
		for _, i := range indexes[1:] {
			err.Relate(items[i].Location, locale.Sprintf(locale.ErrDuplicateValue))
		}
		p.Error(err)
	}
}

func checkXML(isArray, hasItems bool, xml *XML, p *xmlenc.Parser) error {
	if xml.XMLAttr.V() {
		if isArray || hasItems {
			return xml.XMLAttr.Location.NewError(locale.ErrInvalidValue).WithField(xml.XMLAttr.AttributeName.String())
		}

		if xml.XMLWrapped.V() != "" {
			return xml.XMLWrapped.Location.NewError(locale.ErrInvalidValue).WithField(xml.XMLWrapped.AttributeName.String())
		}

		if xml.XMLExtract.V() {
			return xml.XMLExtract.NewError(locale.ErrInvalidValue).WithField(xml.XMLExtract.AttributeName.String())
		}

		if xml.XMLCData.V() {
			return xml.XMLCData.NewError(locale.ErrInvalidValue).WithField(xml.XMLCData.AttributeName.String())
		}
	}

	if xml.XMLWrapped.V() != "" && !isArray {
		return xml.XMLWrapped.NewError(locale.ErrInvalidValue).WithField(xml.XMLWrapped.AttributeName.String())
	}

	if xml.XMLExtract.V() {
		if xml.XMLNSPrefix.V() != "" {
			return xml.XMLNSPrefix.NewError(locale.ErrInvalidValue).WithField(xml.XMLNSPrefix.AttributeName.String())
		}
	}

	return nil
}

// Sanitize 检测内容是否合法
func (doc *APIDoc) Sanitize(p *xmlenc.Parser) {
	if err := doc.checkXMLNamespaces(p); err != nil {
		p.Error(err)
	}
	doc.URI = p.Location.URI

	for _, api := range doc.APIs {
		if api.doc == nil {
			api.doc = doc // 保证单文件的文档能正常解析
			api.URI = doc.URI
		}
		api.sanitizeTags(p)
	}
}

// Sanitize 检测内容是否合法
func (ns *XMLNamespace) Sanitize(p *xmlenc.Parser) {
	if ns.URN.V() == "" {
		p.Error(ns.Location.NewError(locale.ErrIsEmpty, "@urn").WithField("@urn"))
	}
}

func (doc *APIDoc) checkXMLNamespaces(p *xmlenc.Parser) error {
	if len(doc.XMLNamespaces) == 0 {
		return nil
	}

	// 按 URN 查重
	indexes := sliceutil.Dup(doc.XMLNamespaces, func(i, j int) bool {
		return doc.XMLNamespaces[i].URN.V() == doc.XMLNamespaces[j].URN.V()
	})
	if len(indexes) > 0 {
		err := doc.XMLNamespaces[indexes[0]].URN.Location.NewError(locale.ErrDuplicateValue).WithField("@urn")
		for _, i := range indexes[1:] {
			err.Relate(doc.XMLNamespaces[i].Location, locale.Sprintf(locale.ErrDuplicateValue))
		}
		return err
	}

	// 按 prefix 查重
	indexes = sliceutil.Dup(doc.XMLNamespaces, func(i, j int) bool {
		return doc.XMLNamespaces[i].Prefix.V() == doc.XMLNamespaces[j].Prefix.V()
	})
	if len(indexes) > 0 {
		err := doc.XMLNamespaces[indexes[0]].URN.Location.NewError(locale.ErrDuplicateValue).WithField("@prefix")
		for _, i := range indexes[1:] {
			err.Relate(doc.XMLNamespaces[i].Location, locale.Sprintf(locale.ErrDuplicateValue))
		}
		return err
	}

	return nil
}

func (doc *APIDoc) findTag(tag string) *Tag {
	for _, t := range doc.Tags {
		if t.Name.V() == tag {
			return t
		}
	}
	return nil
}

func (doc *APIDoc) findServer(srv string) *Server {
	for _, s := range doc.Servers {
		if s.Name.V() == srv {
			return s
		}
	}
	return nil
}

func (api *API) sanitizeTags(p *xmlenc.Parser) {
	if api.doc == nil {
		panic("api.doc 未获取正确的值")
	}
	api.checkDup(p)

	apiURI := api.URI
	if apiURI == "" {
		apiURI = api.doc.URI
	}

	for _, tag := range api.Tags {
		t := api.doc.findTag(tag.Content.Value)
		if t == nil {
			p.Warning(tag.Content.Location.NewError(locale.ErrInvalidValue).AddTypes(core.ErrorTypeUnused))
			continue
		}

		tag.definition = &Definition{
			Location: t.Location,
			Target:   t,
		}
		t.references = append(t.references, &Reference{
			Location: tag.Location,
			Target:   tag,
		})
	}

	for _, srv := range api.Servers {
		s := api.doc.findServer(srv.Content.Value)
		if s == nil {
			p.Warning(srv.Content.Location.NewError(locale.ErrInvalidValue).AddTypes(core.ErrorTypeUnused))
			continue
		}

		srv.definition = &Definition{
			Location: s.Location,
			Target:   s,
		}
		s.references = append(s.references, &Reference{
			Location: srv.Location,
			Target:   srv,
		})
	}
}

// 检测当前 api 是否与 apidoc.APIs 中存在相同的值
func (api *API) checkDup(p *xmlenc.Parser) {
	err := api.Location.NewError(locale.ErrDuplicateValue)

	for _, item := range api.doc.APIs {
		if item == api {
			continue
		}

		if api.Method.V() != item.Method.V() {
			continue
		}

		p := ""
		if api.Path != nil {
			p = api.Path.Path.V()
		}
		iip := ""
		if item.Path != nil {
			iip = item.Path.Path.V()
		}
		if p != iip {
			continue
		}

		// 默认服务器
		if len(api.Servers) == 0 && len(item.Servers) == 0 {
			err.Relate(item.Location, locale.Sprintf(locale.ErrDuplicateValue))
			continue
		}

		// 判断是否拥有相同的 server 字段
		for _, srv := range api.Servers {
			s := sliceutil.Count(item.Servers, func(i int) bool { return srv.V() == item.Servers[i].V() })
			if s > 0 {
				err.Relate(item.Location, locale.Sprintf(locale.ErrDuplicateValue))
				continue
			}
		}
	}

	if len(err.Related) > 0 {
		p.Error(err)
	}
}
