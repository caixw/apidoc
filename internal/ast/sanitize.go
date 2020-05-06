// SPDX-License-Identifier: MIT

package ast

import (
	"sort"
	"strconv"

	"github.com/issue9/is"

	"github.com/caixw/apidoc/v7/internal/locale"
	"github.com/caixw/apidoc/v7/internal/token"
)

// Sanitize token.Sanitizer
func (p *Param) Sanitize(pp *token.Parser) error {
	if p.Type.V() == TypeNone {
		return pp.NewError(p.Start, p.End, "type", locale.ErrRequired)
	}
	if p.Type.V() == TypeObject && len(p.Items) == 0 {
		return pp.NewError(p.Start, p.End, "param", locale.ErrRequired)
	}

	// 判断 enums 的值是否相同
	if key := getDuplicateEnum(p.Enums); key != "" {
		return pp.NewError(p.Start, p.End, "enum", locale.ErrDuplicateValue)
	}

	if err := chkEnumsType(p.Type, p.Enums, pp); err != nil {
		return err
	}

	// 判断 items 的值是否相同
	if key := getDuplicateItems(p.Items); key != "" {
		return pp.NewError(p.Start, p.End, "param", locale.ErrDuplicateValue)
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
				return p.NewError(enum.Start, enum.End, enum.XMLName.Value, locale.ErrInvalidFormat)
			}
		}
	case TypeBool:
		for _, enum := range enums {
			if _, err := strconv.ParseBool(enum.Value.V()); err != nil {
				return p.NewError(enum.Start, enum.End, enum.XMLName.Value, locale.ErrInvalidFormat)
			}
		}
	case TypeObject, TypeNone:
		return p.NewError(t.Start, t.End, t.XMLName.Value, locale.ErrInvalidValue)
	}

	return nil
}

// 返回重复枚举的值
func getDuplicateEnum(enums []*Enum) string {
	if len(enums) == 0 {
		return ""
	}

	es := make([]string, 0, len(enums))
	for _, e := range enums {
		es = append(es, e.Value.V())
	}
	sort.SliceStable(es, func(i, j int) bool { return es[i] > es[j] })

	for i := 1; i < len(es); i++ {
		if es[i] == es[i-1] {
			return es[i]
		}
	}

	return ""
}

func getDuplicateItems(items []*Param) string {
	if len(items) == 0 {
		return ""
	}

	es := make([]string, 0, len(items))
	for _, e := range items {
		es = append(es, e.Name.V())
	}
	sort.SliceStable(es, func(i, j int) bool { return es[i] > es[j] })

	for i := 1; i < len(es); i++ {
		if es[i] == es[i-1] {
			return es[i]
		}
	}

	return ""
}

func checkXML(isArray, hasItems bool, xml *XML, p *token.Parser) error {
	if xml.XMLAttr.V() {
		if isArray || hasItems {
			return p.NewError(xml.XMLAttr.Start, xml.XMLAttr.End, xml.XMLAttr.XMLName.Value, locale.ErrInvalidValue)
		}

		if xml.XMLWrapped.V() != "" {
			return p.NewError(xml.XMLWrapped.Start, xml.XMLWrapped.End, xml.XMLWrapped.XMLName.Value, locale.ErrInvalidValue)
		}

		if xml.XMLExtract.V() {
			return p.NewError(xml.XMLExtract.Start, xml.XMLExtract.End, xml.XMLExtract.XMLName.Value, locale.ErrInvalidValue)
		}

		if xml.XMLNS.V() != "" {
			return p.NewError(xml.XMLNS.Start, xml.XMLNS.End, xml.XMLNS.XMLName.Value, locale.ErrInvalidValue)
		}

		if xml.XMLNSPrefix.V() != "" {
			return p.NewError(xml.XMLNSPrefix.Start, xml.XMLNSPrefix.End, xml.XMLNSPrefix.XMLName.Value, locale.ErrInvalidValue)
		}
	}

	if xml.XMLWrapped.V() != "" && !isArray {
		return p.NewError(xml.XMLWrapped.Start, xml.XMLWrapped.End, xml.XMLWrapped.XMLName.Value, locale.ErrInvalidValue)
	}

	if xml.XMLExtract.V() {
		if xml.XMLNS.V() != "" {
			return p.NewError(xml.XMLNS.Start, xml.XMLNS.End, xml.XMLNS.XMLName.Value, locale.ErrInvalidValue)
		}

		if xml.XMLNSPrefix.V() != "" {
			return p.NewError(xml.XMLNSPrefix.Start, xml.XMLNSPrefix.End, xml.XMLNSPrefix.XMLName.Value, locale.ErrInvalidValue)
		}
	}

	// 有命名空间，必须要有前缀
	if xml.XMLNS.V() != "" && xml.XMLNSPrefix.V() == "" {
		return p.NewError(xml.XMLNSPrefix.Start, xml.XMLNSPrefix.End, xml.XMLNSPrefix.XMLName.Value, locale.ErrInvalidValue)
	}

	return nil
}

// Sanitize 检测内容是否合法
func (doc *APIDoc) Sanitize(p *token.Parser) error {
	// TODO

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

	for _, api := range doc.Apis {
		if api.doc == nil {
			api.doc = doc // 保证单文件的文档能正常解析
		}
		if err := api.sanitizeTags(p); err != nil {
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

func (api *API) sanitizeTags(p *token.Parser) error {
	if api.doc == nil {
		panic("api.doc 未获取正确的值")
	}

	for _, tag := range api.Tags {
		if !api.doc.tagExists(tag.Content.Value) {
			return p.NewError(tag.Content.Start, tag.Content.End, "content", locale.ErrInvalidValue)
		}
	}

	for _, srv := range api.Servers {
		if !api.doc.serverExists(srv.Content.Value) {
			return p.NewError(srv.Content.Start, srv.Content.End, "content", locale.ErrInvalidValue)
		}
	}

	return nil
}
