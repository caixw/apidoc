// SPDX-License-Identifier: MIT

package doc

import "github.com/caixw/apidoc/v5/internal/locale"

// XML 仅作用于 XML 的几个属性
type XML struct {
	XMLAttr     bool   `xml:"xml-attr,attr,omitempty"`      // 作为父元素的 XML 属性存在
	XMLExtract  bool   `xml:"xml-extract,attr,omitempty"`   // 提取当前内容作为父元素的内容
	XMLNS       string `xml:"xml-ns,attr,omitempty"`        // 命名空间
	XMLNSPrefix string `xml:"xml-ns-prefix,attr,omitempty"` // 命名空间前缀
}

func checkXML(xml *XML, field string) error {
	if xml.XMLAttr {
		if xml.XMLExtract {
			return newSyntaxError(field+"/@xml-extract", locale.ErrInvalidValue)
		}

		if xml.XMLNS != "" {
			return newSyntaxError(field+"/@xml-ns", locale.ErrInvalidValue)
		}

		if xml.XMLNSPrefix != "" {
			return newSyntaxError(field+"/@xml-ns-prefix", locale.ErrInvalidValue)
		}
	}

	if xml.XMLExtract {
		if xml.XMLNS != "" {
			return newSyntaxError(field+"/@xml-ns", locale.ErrInvalidValue)
		}

		if xml.XMLNSPrefix != "" {
			return newSyntaxError(field+"/@xml-ns-prefix", locale.ErrInvalidValue)
		}
	}

	// 有命名空间，必须要有前缀
	if xml.XMLNS != "" && xml.XMLNSPrefix == "" {
		return newSyntaxError(field+"/@xml-ns-prefix", locale.ErrInvalidValue)
	}

	return nil
}
