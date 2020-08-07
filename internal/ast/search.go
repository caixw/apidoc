// SPDX-License-Identifier: MIT

package ast

import (
	"reflect"
	"unicode"

	"github.com/issue9/sliceutil"

	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/internal/node"
)

var rangerType = reflect.TypeOf((*core.Ranger)(nil)).Elem()

// Search 搜索符合条件的对象并返回
//
// 从 doc 中查找符合符合 pos 定位的最小对象，且该对象必须实现了 t 类型。
// 如果不存在则返回 nil。
func (doc *APIDoc) Search(uri core.URI, pos core.Position, t reflect.Type) (r core.Ranger) {
	if doc.URI == uri {
		r = search(reflect.ValueOf(doc), pos, t, "APIs")
	}

	for _, api := range doc.APIs {
		matched := api.URI == uri || (api.URI == "" && doc.URI == uri)
		if !matched {
			continue
		}

		if rr := search(reflect.ValueOf(api), pos, t); rr != nil {
			return rr
		}
	}

	return r
}

func search(v reflect.Value, pos core.Position, t reflect.Type, exclude ...string) (r core.Ranger) {
	if v.IsZero() {
		return
	}

	v = node.RealValue(v)

	if v.CanInterface() && v.Type().Implements(rangerType) && (t == nil || v.Type().Implements(t)) {
		if rr := v.Interface().(core.Ranger); rr.Contains(pos) {
			r = rr
		}
	} else if v.CanAddr() {
		if pv := v.Addr(); pv.CanInterface() && pv.Type().Implements(rangerType) && (t == nil || pv.Type().Implements(t)) {
			if rr := pv.Interface().(core.Ranger); rr.Contains(pos) {
				r = rr
			}
		}
	}

	if v.Kind() == reflect.Struct {
		for vt, i := v.Type(), 0; i < vt.NumField(); i++ {
			ft := vt.Field(i)
			if ft.Anonymous ||
				unicode.IsLower(rune(ft.Name[0])) ||
				sliceutil.Count(exclude, func(i int) bool { return exclude[i] == ft.Name }) > 0 { // 需要过滤的字段
				continue
			}

			fv := v.Field(i)
			if fv.Kind() == reflect.Array || fv.Kind() == reflect.Slice {
				for j := 0; j < fv.Len(); j++ {
					if rr := search(fv.Index(j), pos, t, exclude...); rr != nil {
						return rr
					}
				}
				continue
			} else if rr := search(fv, pos, t, exclude...); rr != nil {
				return rr
			}
		}
	}

	return r
}
