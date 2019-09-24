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
	Path      string   `xml:"path,attr"`
	Params    []*Param `xml:"param,omitempty"`
	Queries   []*Param `xml:"query,omitempty"`
	Reference string   `xml:"ref,attr,omitempty"`
}

type shadowPath Path

// UnmarshalXML xml.Unmarshaler
func (p *Path) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var pp shadowPath
	if err := d.DecodeElement(&pp, &start); err != nil {
		return err
	}

	if pp.Path == "" {
		return locale.Errorf(locale.ErrRequired, "path")
	}

	params, err := parsePath(pp.Path)
	if err != nil {
		return err
	}
	if len(params) != len(pp.Params) {
		return locale.Errorf(locale.ErrPathNotMatchParams)
	}
	for _, param := range pp.Params {
		if _, found := params[param.Name]; !found {
			return locale.Errorf(locale.ErrPathNotMatchParams, param.Name)
		}
	}

	// queries 不作判断，同名会被当作数据处理

	*p = Path(pp)
	return nil
}

func parsePath(path string) (params map[string]struct{}, err error) {
	start := -1
	for i, b := range path {
		switch b {
		case '{':
			if start != -1 {
				return nil, locale.Errorf(locale.ErrPathSyntaxError, path)
			}

			start = i + 1
		case '}':
			if start == -1 {
				return nil, locale.Errorf(locale.ErrPathSyntaxError, path)
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
		return nil, locale.Errorf(locale.ErrPathSyntaxError, path)
	}

	return params, nil
}
