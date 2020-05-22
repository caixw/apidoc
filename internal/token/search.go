// SPDX-License-Identifier: MIT

package token

import (
	"reflect"

	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/internal/locale"
	"github.com/caixw/apidoc/v7/internal/node"
)

var baseType = reflect.TypeOf((*baser)(nil)).Elem()

type baser interface {
	contains(core.Range) bool
	usage() string
}

func (b *Base) contains(r core.Range) bool {
	return b.Contains(r.Start) && b.Contains(r.End)
}

func (b *Base) usage() string {
	return locale.Sprintf(b.UsageKey)
}

// SearchUsage 根据 r 从 v 中查找相应的 usage 字段内容
func SearchUsage(v reflect.Value, r core.Range, exclude ...string) (usage string, contains bool) {
	v = node.GetRealValue(v)
	if usage, contains = getUsage(v, r); !contains {
		return
	}

	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		tf := t.Field(i)
		if tf.Anonymous || // 不考虑匿名字段，因为如果有实现接口也已经被当前对象使用。
			inStringSlice(exclude, tf.Name) { // 需要过滤的字段
			continue
		}

		vf := node.GetRealValue(v.Field(i))
		if vf.Kind() == reflect.Array || vf.Kind() == reflect.Slice {
			for j := 0; j < vf.Len(); j++ {
				if u, c := SearchUsage(vf.Index(j), r); c {
					return u, c
				}
			}
		} else {
			if u, c := SearchUsage(vf, r); c {
				return u, c
			}
		}
	}

	return usage, true
}

func getUsage(v reflect.Value, r core.Range) (string, bool) {
	if v.Type().Implements(baseType) && v.CanInterface() {
		base := v.Interface().(baser)
		return base.usage(), base.contains(r)
	} else if v.CanAddr() {
		if pv := v.Addr(); pv.Type().Implements(baseType) && pv.CanInterface() {
			base := pv.Interface().(baser)
			return base.usage(), base.contains(r)
		}
	}
	return "", false
}

func inStringSlice(slice []string, key string) bool {
	for _, v := range slice {
		if v == key {
			return true
		}
	}
	return false
}
