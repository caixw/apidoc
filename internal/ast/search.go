// SPDX-License-Identifier: MIT

package ast

import (
	"reflect"
	"unicode"

	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/internal/node"
)

var searcherType = reflect.TypeOf((*core.Searcher)(nil)).Elem()

// Search 搜索符合条件的对象并返回
//
// 从 doc 中查找符合符合 pos 定位的最小对象，且该对象必须实现了 t 类型。
// 如果不存在则返回 nil。t 必须是一个接口。
func (doc *APIDoc) Search(uri core.URI, pos core.Position, t reflect.Type) (r core.Searcher) {
	r = search(reflect.ValueOf(doc), uri, pos, t)
	if r == nil { // apidoc 的 uri 可以与 api 的 uri 不同
		for _, api := range doc.APIs {
			if rr := search(reflect.ValueOf(api), uri, pos, t); rr != nil {
				return rr
			}
		}
	}

	return r
}

func search(v reflect.Value, uri core.URI, pos core.Position, t reflect.Type) (r core.Searcher) {
	if v.IsZero() {
		return
	}

	v = node.RealValue(v)

	if v.CanInterface() && v.Type().Implements(searcherType) && (t == nil || v.Type().Implements(t)) {
		if rr := v.Interface().(core.Searcher); rr.Contains(uri, pos) {
			r = rr
		}
	} else if v.CanAddr() {
		if pv := v.Addr(); pv.CanInterface() && pv.Type().Implements(searcherType) && (t == nil || pv.Type().Implements(t)) {
			if rr := pv.Interface().(core.Searcher); rr.Contains(uri, pos) {
				r = rr
			}
		}
	}

	if r == nil && t == nil { // 不匹配当前元素，也不需要搜查子元素是否实现 t，则直接返回 nil。
		return nil
	}

	if v.Kind() == reflect.Struct {
		for vt, i := v.Type(), 0; i < vt.NumField(); i++ {
			ft := vt.Field(i)
			if ft.Anonymous || unicode.IsLower(rune(ft.Name[0])) {
				continue
			}

			fv := v.Field(i)
			if fv.Kind() == reflect.Array || fv.Kind() == reflect.Slice {
				for j := 0; j < fv.Len(); j++ {
					if rr := search(fv.Index(j), uri, pos, t); rr != nil {
						return rr
					}
				}
				continue
			} else if rr := search(fv, uri, pos, t); rr != nil {
				return rr
			}
		}
	}

	return r
}
