// SPDX-License-Identifier: MIT

package xmlenc

import (
	"errors"
	"fmt"
	"io"
	"reflect"

	"github.com/issue9/is"

	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/internal/locale"
	"github.com/caixw/apidoc/v7/internal/node"
)

type (
	// Decoder 实现从 p 中解码内容到当前对象的值
	Decoder interface {
		// 从 p 中读取内容并实例化到当前对象中
		//
		// 必须要同时从 p 中读取相应的 EndElement 才能返回。
		// end 表示 EndElement.End 的值。
		//
		// NOTE: 如果是自闭合标签，则不会调用该接口。
		//
		// 接口应该只返回 *core.Error 作为错误对象。
		DecodeXML(p *Parser, start *StartElement) (end *EndElement, err error)
	}

	// AttrDecoder 实现从 attr 中解码内容到当前对象的值
	AttrDecoder interface {
		// 解析属性值
		//
		// 接口应该只返回 *core.Error 作为错误对象。
		DecodeXMLAttr(p *Parser, attr *Attribute) error
	}

	// Sanitizer 用于验证和修改对象中的数据
	Sanitizer interface {
		// 验证数据是否正确
		//
		// 可以通过 p.NewError 和 p.WithError 生成信息向 p.MessageHandler 输出
		Sanitize(p *Parser)
	}

	attributeSetter interface {
		setAttribute(string, *Attribute)
	}

	tagSetter interface {
		setTag(string, *StartElement, *EndElement)
	}

	decoder struct {
		p      *Parser
		prefix string // 命名空间所表示的前缀
	}
)

var (
	attrDecoderType = reflect.TypeOf((*AttrDecoder)(nil)).Elem()
	decoderType     = reflect.TypeOf((*Decoder)(nil)).Elem()
	sanitizerType   = reflect.TypeOf((*Sanitizer)(nil)).Elem()
)

func (b *BaseAttribute) setAttribute(usage string, attr *Attribute) {
	b.UsageKey = usage
	b.Range = attr.Range
	b.AttributeName = attr.Name
}

func (b *BaseTag) setTag(usage string, start *StartElement, end *EndElement) {
	b.UsageKey = usage
	b.StartTag = start.Name

	if end == nil {
		b.Range = start.Range
	} else {
		b.Range.Start = start.Start
		b.Range.End = end.End
		b.EndTag = end.Name
	}
}

// Decode 将 p 中的 XML 内容解码至 v 对象中
//
// namespace 如果不为空表示当前的 xml 所在的命名空间，只有该命名空间的元素才会被正确识别。
func Decode(p *Parser, v interface{}, namespace string) {
	d := &decoder{
		p: p,
	}

	var hasRoot bool
	for {
		t, r, err := p.Token()
		if errors.Is(err, io.EOF) {
			return
		} else if err != nil {
			d.p.Error(err)
			return
		}

		switch elem := t.(type) {
		case *StartElement:
			if hasRoot { // 多个根元素
				_ = p.endElement(elem) // 找到对应的结束标签，忽略错误
				msg := p.NewError(elem.Start, p.Current().Position, elem.Name.String(), locale.ErrMultipleRootTag).
					AddTypes(core.ErrorTypeUnused)
				d.p.Warning(msg)
				return
			}
			hasRoot = true

			d.prefix = findPrefix(elem, namespace)
			if !d.decodeElement(elem, node.ParseValue(reflect.ValueOf(v))) {
				return
			}
		case *Comment, *String, *Instruction: // 忽略注释和普通的文本内容
		default:
			d.p.Error(p.NewError(r.Start, r.End, "", locale.ErrInvalidXML))
			return
		}
	}
}

func findPrefix(start *StartElement, namespace string) string {
	for _, attr := range start.Attributes {
		if attr.Name.Prefix.Value == "xmlns" && attr.Value.Value == namespace {
			return attr.Name.Local.Value
		}
	}
	return ""
}

func (d *decoder) decode(n *node.Node, start *StartElement) *EndElement {
	d.decodeAttributes(n, start)

	if start.SelfClose {
		d.checkOmitempty(n, start.Start, start.End, start.Name.String())
		return nil
	}

	if end, ok := d.decodeElements(n); ok {
		d.checkOmitempty(n, start.Start, end.End, start.Name.String())
		return end
	}
	return nil
}

// 判断 omitempty 属性
func (d *decoder) checkOmitempty(n *node.Node, start, end core.Position, field string) {
	for _, attr := range n.Attributes {
		if canNotEmpty(attr) {
			d.p.Error(d.p.NewError(start, end, attr.Name, locale.ErrIsEmpty, attr.Name))
		}
	}
	for _, elem := range n.Elements {
		if canNotEmpty(elem) {
			d.p.Error(d.p.NewError(start, end, elem.Name, locale.ErrIsEmpty, elem.Name))
		}
	}
	if n.CData != nil && canNotEmpty(n.CData) {
		d.p.Error(d.p.NewError(start, end, "cdata", locale.ErrIsEmpty, field))
	}
	if n.Content != nil && canNotEmpty(n.Content) {
		d.p.Error(d.p.NewError(start, end, "content", locale.ErrIsEmpty, field))
	}
}

// 当前表示的值必须是一个非空值
func canNotEmpty(v *node.Value) bool {
	return v.Name != "" && // cdata 和 content 在未初始化时 name 字段为空值
		!v.Omitempty &&
		(!v.CanInterface() || is.Empty(v.Interface(), true))
}

// 将 start 的属性内容解码到 obj.Attributes 之中
func (d *decoder) decodeAttributes(n *node.Node, start *StartElement) {
	if start == nil {
		return
	}

	for _, attr := range start.Attributes {
		item, found := n.Attribute(attr.Name.Local.Value)
		if !found || d.prefix != attr.Name.Prefix.Value { // 未找到或是命名空间不匹配
			continue
		}
		v := node.RealValue(item.Value)
		v.Set(reflect.New(v.Type()).Elem())

		var impl bool
		if item.CanInterface() && item.Type().Implements(attrDecoderType) {
			if err := item.Interface().(AttrDecoder).DecodeXMLAttr(d.p, attr); err != nil {
				d.p.Error(err)
			}
			impl = true
		} else if item.CanAddr() {
			pv := item.Addr()
			if pv.CanInterface() && pv.Type().Implements(attrDecoderType) {
				if err := pv.Interface().(AttrDecoder).DecodeXMLAttr(d.p, attr); err != nil {
					d.p.Error(err)
				}
				impl = true
			}
		}

		if !impl {
			panic(fmt.Sprintf("当前属性 %s 未实现 AttrDecoder 接口", attr.Name))
		}

		v.Addr().Interface().(attributeSetter).setAttribute(item.Usage, attr)
		callSanitizer(v, d.p) // Sanitize 在最后调用，可以确保能取到 v.Range
	}
}

func (d *decoder) decodeElements(n *node.Node) (end *EndElement, ok bool) {
	for {
		t, r, err := d.p.Token()
		if errors.Is(err, io.EOF) {
			// 应该只有 EndElement 才能返回，否则就不完整的 XML
			d.p.Error(d.p.NewError(d.p.Current().Position, d.p.Current().Position, "", locale.ErrNotFoundEndTag))
			return nil, false
		} else if err != nil {
			d.p.Error(err)
			return nil, false
		}

		switch elem := t.(type) {
		case *EndElement: // 找到当前对象的结束标签
			if (elem.Name.Local.Value == n.Value.Name) && (elem.Name.Prefix.Value == d.prefix) {
				return elem, true
			}
			d.p.Error(d.p.NewError(elem.Start, elem.End, n.Value.Name, locale.ErrNotFoundEndTag))
			return nil, false
		case *CData:
			if n.CData != nil {
				copyContentValue(n.CData.Value, reflect.ValueOf(elem))
			}
		case *String:
			if n.Content != nil {
				copyContentValue(n.Content.Value, reflect.ValueOf(elem))
			}
		case *StartElement:
			item, found := n.Element(elem.Name.Local.Value)
			if !found || d.prefix != elem.Name.Prefix.Value {
				if err := d.p.endElement(elem); err != nil {
					d.p.Error(err)
					return nil, false
				}

				e := d.p.NewError(elem.Start, d.p.Current().Position, elem.Name.String(), locale.ErrInvalidTag).
					AddTypes(core.ErrorTypeUnused)
				d.p.Warning(e)
				break // 忽略不存在的子元素
			}
			if !d.decodeElement(elem, item) {
				return nil, false
			}
		case *Comment: // 忽略注释内容
		default:
			d.p.Error(d.p.NewError(r.Start, r.End, "", locale.ErrInvalidXML))
			return nil, false
		}
	}
}

func copyContentValue(target, source reflect.Value) {
	target = node.RealValue(target)
	source = node.RealValue(source)

	st := source.Type()
	num := st.NumField()
	for i := 0; i < num; i++ {
		name := st.Field(i).Name
		target.FieldByName(name).
			Set(source.Field(i))
	}
}

func (d *decoder) decodeElement(start *StartElement, v *node.Value) (ok bool) {
	v.Value = node.RealValue(v.Value)
	k := v.Kind()
	switch {
	case k == reflect.Ptr, k == reflect.Func, k == reflect.Chan, k == reflect.Array, node.IsPrimitive(v.Value):
		panic(fmt.Sprintf("%s 是无效的类型", v.Value.Type()))
	case k == reflect.Slice:
		return d.decodeSlice(start, v)
	}

	end, impl, err := callDecodeXML(v.Value, d.p, start)
	if !impl {
		end = d.decode(node.New(start.Name.Local.Value, v.Value), start)
	}
	if err != nil {
		d.p.Error(err)
		return false
	}

	d.setTagValue(v.Value, v.Usage, start, end)
	return true
}

func (d *decoder) decodeSlice(start *StartElement, slice *node.Value) (ok bool) {
	// 不相配，表示当前元素找不到与之相配的元素，需要忽略这个元素，
	// 所以要过滤与 start 想匹配的结束符号才算结束。
	if !start.SelfClose && (start.Name.Local.Value != slice.Name) {
		if err := d.p.endElement(start); err != nil {
			d.p.Error(err)
			return false
		}
	}

	elem := reflect.New(slice.Type().Elem()).Elem()
	if elem.Kind() == reflect.Ptr && elem.IsNil() {
		elem.Set(reflect.New(elem.Type().Elem()))
	}

	end, impl, err := callDecodeXML(elem, d.p, start)
	if !impl {
		if node.IsPrimitive(elem) {
			panic(fmt.Sprintf("%s:%s 必须实现 Decoder 接口", slice.Name, elem.Type()))
		}
		end = d.decode(node.New(start.Name.Local.Value, elem), start)
	}
	if err != nil {
		d.p.Error(err)
		return false
	}

	d.setTagValue(node.RealValue(elem), slice.Usage, start, end)
	slice.Value.Set(reflect.Append(slice.Value, elem))
	return true
}

// 调用 v 的 DecodeXML 接口方法
//
// 当 impl 为 true 时，err 表示的是 DecodeXML 接口返回的错误，否则 err 永远为 nil
func callDecodeXML(v reflect.Value, p *Parser, start *StartElement) (end *EndElement, impl bool, err error) {
	if v.CanInterface() && v.Type().Implements(decoderType) {
		if !start.SelfClose {
			end, err = v.Interface().(Decoder).DecodeXML(p, start)
		}
		return end, true, err
	} else if v.CanAddr() {
		pv := v.Addr()
		if pv.CanInterface() && pv.Type().Implements(decoderType) {
			if !start.SelfClose {
				end, err = pv.Interface().(Decoder).DecodeXML(p, start)
			}
			return end, true, err
		}
	}
	return nil, false, nil
}

func (d *decoder) setTagValue(v reflect.Value, usage string, start *StartElement, end *EndElement) {
	v.Addr().Interface().(tagSetter).setTag(usage, start, end)
	callSanitizer(v, d.p) // Sanitize 在最后调用，可以确保能取到 v.Range
}

func callSanitizer(v reflect.Value, p *Parser) {
	if v.CanInterface() && v.Type().Implements(sanitizerType) {
		v.Interface().(Sanitizer).Sanitize(p)
	} else if v.CanAddr() {
		pv := v.Addr()
		if pv.CanInterface() && pv.Type().Implements(sanitizerType) {
			pv.Interface().(Sanitizer).Sanitize(p)
		}
	}
}
