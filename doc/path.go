// SPDX-License-Identifier: MIT

package doc

import (
	"encoding/xml"

	"github.com/caixw/apidoc/v5/internal/locale"
)

// Path 路径信息
//  <path path="/users/{id}">
//      <param name="id" type="number" summary="summary" />
//      <query name="page" type="number" summary="page" default="1" />
//  </path>
type Path struct {
	Path      string         `xml:"path,attr"`
	Params    []*SimpleParam `xml:"param,omitempty"`
	Queries   []*SimpleParam `xml:"query,omitempty"`
	Reference string         `xml:"ref,attr,omitempty"`
}

type shadowPath Path

// UnmarshalXML xml.Unmarshaler
func (p *Path) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	field := "/" + start.Name.Local
	shadow := (*shadowPath)(p)
	if err := d.DecodeElement(shadow, &start); err != nil {
		return fixedSyntaxError(err, "", field, 0)
	}

	if shadow.Path == "" {
		return newSyntaxError(field+"#path", locale.ErrRequired)
	}

	params, err := parsePath(shadow.Path)
	if err != nil {
		return fixedSyntaxError(err, "", field+"#path", 0)
	}
	if len(params) != len(shadow.Params) {
		return newSyntaxError(field+"#path", locale.ErrPathNotMatchParams)
	}
	for _, param := range shadow.Params {
		if _, found := params[param.Name]; !found {
			return newSyntaxError(field+"#path", locale.ErrPathNotMatchParams)
		}
	}

	// queries 不作判断，同名会被当作数据处理

	return nil
}

func parsePath(path string) (params map[string]struct{}, err error) {
	start := -1
	for i, b := range path {
		switch b {
		case '{':
			if start != -1 {
				return nil, locale.Errorf(locale.ErrInvalidFormat)
			}

			start = i + 1
		case '}':
			if start == -1 {
				return nil, locale.Errorf(locale.ErrInvalidFormat)
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
		return nil, locale.Errorf(locale.ErrInvalidFormat)
	}

	return params, nil
}
