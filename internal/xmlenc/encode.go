// SPDX-License-Identifier: MIT

package xmlenc

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"reflect"

	"github.com/issue9/validation/is"

	"github.com/caixw/apidoc/v7/internal/node"
)

// Encoder 将元素内容编码成 XML 内容
type Encoder interface {
	// EncodeXML 仅需要返回元素内容的 XML 编码，不需要包含本身的标签和属性。
	EncodeXML() (string, error)
}

// AttrEncoder 将属性值编码成符合 XML 规范的值
type AttrEncoder interface {
	// EncodeXMLAttr 仅需要返回属性的 XML 表示，不需要包含属性值的引号字符。
	EncodeXMLAttr() (string, error)
}

var (
	attrEncoderType = reflect.TypeOf((*AttrEncoder)(nil)).Elem()
	encoderType     = reflect.TypeOf((*Encoder)(nil)).Elem()
)

// Encode 将 v 转换成 XML 内容
//
// namespace 指定 XML 的命名空间；
// prefix 命名空间的前缀，如果为空表示使用默认命名空间；
func Encode(indent string, v interface{}, namespace, prefix string) ([]byte, error) {
	if prefix != "" && namespace == "" {
		panic("参数 prefix 不为空的情况下，必须指定 namespace")
	}

	buf := new(bytes.Buffer)
	e := xml.NewEncoder(buf)
	e.Indent("", indent)

	if err := encode(node.New("", reflect.ValueOf(v)), e, namespace, prefix, true); err != nil {
		return nil, err
	}

	if err := e.Flush(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func encode(n *node.Node, e *xml.Encoder, namespace, prefix string, root bool) error {
	start, err := buildStartElement(n, namespace, prefix, root)
	if err != nil {
		return err
	}

	if n.CData != nil && !isOmitempty(n.CData) {
		chardata, err := getContentValue(n.CData.Value)
		if err != nil {
			return err
		}

		return e.EncodeElement(struct {
			string `xml:",cdata"`
		}{chardata}, start)
	}

	if n.Content != nil && !isOmitempty(n.Content) {
		chardata, err := getContentValue(n.Content.Value)
		if err != nil {
			return err
		}
		return e.EncodeElement(xml.CharData(chardata), start)
	}

	return encodeElements(n, e, start, namespace, prefix)
}

func encodeElements(n *node.Node, e *xml.Encoder, start xml.StartElement, namespace, prefix string) (err error) {
	if err = e.EncodeToken(start); err != nil {
		return err
	}
	for _, v := range n.Elements {
		if err := encodeElement(e, v, namespace, prefix); err != nil {
			return err
		}
	}

	return e.EncodeToken(xml.EndElement{Name: buildXMLName(n.Value.Name, prefix)})
}

func encodeElement(e *xml.Encoder, v *node.Value, namespace, prefix string) (err error) {
	if isOmitempty(v) {
		return nil
	}

	var chardata string
	var found bool

	if v.CanInterface() && v.Type().Implements(encoderType) {
		if chardata, err = v.Interface().(Encoder).EncodeXML(); err != nil {
			return err
		}
		found = true
	} else if !found && v.CanAddr() {
		pv := v.Addr()
		if pv.CanInterface() && pv.Type().Implements(encoderType) {
			if chardata, err = pv.Interface().(Encoder).EncodeXML(); err != nil {
				return err
			}
			found = true
		}
	}
	if !found && node.IsPrimitive(v.Value) {
		chardata = fmt.Sprint(v.Interface())
		found = true
	}

	if found {
		start := xml.StartElement{Name: buildXMLName(v.Name, prefix)}
		if err := e.EncodeElement(xml.CharData(chardata), start); err != nil {
			return err
		}
		return nil
	}

	if v.Kind() == reflect.Array || v.Kind() == reflect.Slice {
		for i := 0; i < v.Len(); i++ {
			val := node.NewValue(v.Name, v.Index(i), v.Omitempty, v.Usage)
			if err := encodeElement(e, val, namespace, prefix); err != nil {
				return err
			}
		}
		return nil
	}

	return encode(node.New(v.Name, v.Value), e, namespace, prefix, false)
}

func buildStartElement(n *node.Node, namespace, prefix string, root bool) (xml.StartElement, error) {
	start := xml.StartElement{
		Name: buildXMLName(n.Value.Name, prefix),
		Attr: make([]xml.Attr, 0, len(n.Attributes)+1), // +1 表示 xmlns 的字段
	}

	for _, v := range n.Attributes {
		if isOmitempty(v) {
			continue
		}

		val, err := getAttributeValue(v.Value)
		if err != nil {
			return xml.StartElement{}, err
		}

		start.Attr = append(start.Attr, xml.Attr{
			Name:  buildXMLName(v.Name, prefix),
			Value: val,
		})
	}

	if root && namespace != "" {
		name := "xmlns"
		if prefix != "" {
			name += ":" + prefix
		}
		start.Attr = append(start.Attr, xml.Attr{
			Name:  xml.Name{Local: name},
			Value: namespace,
		})
	}

	return start, nil
}

func buildXMLName(name, prefix string) xml.Name {
	if prefix == "" {
		return xml.Name{Local: name}
	}
	return xml.Name{Local: prefix + ":" + name}
}

func getAttributeValue(elem reflect.Value) (string, error) {
	if elem.CanInterface() && elem.Type().Implements(attrEncoderType) {
		return elem.Interface().(AttrEncoder).EncodeXMLAttr()
	} else if elem.CanAddr() {
		pv := elem.Addr()
		if pv.CanInterface() && pv.Type().Implements(attrEncoderType) {
			return pv.Interface().(AttrEncoder).EncodeXMLAttr()
		}
	}

	return fmt.Sprint(elem.Interface()), nil
}

// 获取 cdata 和 content 节点的的内容
func getContentValue(elem reflect.Value) (string, error) {
	elem = node.RealValue(elem)
	if elem.CanInterface() && elem.Type().Implements(encoderType) {
		return elem.Interface().(Encoder).EncodeXML()
	} else if elem.CanAddr() {
		if pv := elem.Addr(); pv.CanInterface() && pv.Type().Implements(encoderType) {
			return pv.Interface().(Encoder).EncodeXML()
		}
	}

	return fmt.Sprint(elem.Interface()), nil
}

func isOmitempty(v *node.Value) bool {
	return v.Omitempty && is.Empty(v.Value.Interface(), true)
}
