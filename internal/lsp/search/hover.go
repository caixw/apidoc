// SPDX-License-Identifier: MIT

package search

import (
	"reflect"
	"unicode"

	"github.com/issue9/sliceutil"

	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/internal/ast"
	"github.com/caixw/apidoc/v7/internal/lsp/protocol"
	"github.com/caixw/apidoc/v7/internal/node"
	"github.com/caixw/apidoc/v7/internal/xmlenc"
)

var baseType = reflect.TypeOf(xmlenc.Base{})

// Hover 从 doc 查找最符合 uri 和 pos 条件的元素并赋值给 hover
//
// 返回值表示是否找到了相应在的元素。
func Hover(doc *ast.APIDoc, uri core.URI, pos core.Position, hover *protocol.Hover) {
	setHover := func(b *xmlenc.Base) {
		hover.Range = b.Range
		hover.Contents = protocol.MarkupContent{
			Kind:  protocol.MarkupKinMarkdown,
			Value: b.Usage(),
		}
	}

	if doc.URI == uri {
		if b := usage(reflect.ValueOf(doc), pos, "APIs"); b != nil {
			setHover(b)
		}
	}

	for _, api := range doc.APIs {
		matched := api.URI == uri || (api.URI == "" && doc.URI == uri)
		if !matched {
			continue
		}

		if b := usage(reflect.ValueOf(api), pos); b != nil {
			setHover(b)
		}
	}
}

// 从 v 中查找最匹配 pos 位置的元素，如果找到匹配项，还会查找其子项，是不是匹配度更高。
func usage(v reflect.Value, pos core.Position, exclude ...string) (b *xmlenc.Base) {
	v = node.RealValue(v)
	if b = getBase(v, pos); b == nil {
		return nil
	}

	// 查询 b 的子元素中是否有更精确的匹配
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		tf := t.Field(i)
		if tf.Anonymous || // 不考虑匿名字段，因为如果匿名字段正好是 Base 实例的话，应该在函开始处已经被处理。
			unicode.IsLower(rune(tf.Name[0])) ||
			sliceutil.Count(exclude, func(i int) bool { return exclude[i] == tf.Name }) > 0 { // 需要过滤的字段
			continue
		}

		vf := node.RealValue(v.Field(i))
		if vf.Kind() == reflect.Array || vf.Kind() == reflect.Slice {
			for j := 0; j < vf.Len(); j++ {
				if b2 := usage(vf.Index(j), pos); b2 != nil {
					return b2
				}
			}
		} else {
			if b2 := usage(vf, pos); b2 != nil {
				return b2
			}
		}
	}

	return b // 没有更精确的匹配，则返回 b
}

func getBase(v reflect.Value, pos core.Position) *xmlenc.Base {
	t := v.Type()
	switch {
	case t == baseType:
		if b := v.Interface().(xmlenc.Base); b.Contains(pos) && b.UsageKey != nil {
			return &b
		}
	case t.Kind() == reflect.Struct:
		for i := 0; i < t.NumField(); i++ {
			if !t.Field(i).Anonymous {
				continue
			}

			vf := node.RealValue(v.Field(i))
			if b := getBase(vf, pos); b != nil {
				return b
			}
		}
	}

	return nil
}
