// SPDX-License-Identifier: MIT

package token

import (
	"reflect"

	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/internal/locale"
	"github.com/caixw/apidoc/v7/internal/node"
	"github.com/issue9/sliceutil"
)

var tiperType = reflect.TypeOf((*tiper)(nil)).Elem()

// Tip 定义了 LSP 查找功能返回的提示内容
type Tip struct {
	core.Range
	Usage string
}

type tiper interface {
	contains(core.Position) bool
	tip() *Tip
}

func (b *Base) contains(pos core.Position) bool {
	return b.Contains(pos)
}

func (b *Base) tip() *Tip {
	return &Tip{
		Range: b.Range,
		Usage: locale.Sprintf(b.UsageKey),
	}
}

// SearchUsage 根据 r 从 v 中查找相应的 usage 字段内容
func SearchUsage(v reflect.Value, pos core.Position, exclude ...string) (tip *Tip) {
	v = node.GetRealValue(v)
	if tip = getUsage(v, pos); tip == nil {
		return
	}

	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		tf := t.Field(i)
		if tf.Anonymous || // 不考虑匿名字段，因为如果有实现接口也已经被当前对象使用。
			sliceutil.Count(exclude, func(i int) bool { return exclude[i] == tf.Name }) > 0 { // 需要过滤的字段
			continue
		}

		vf := node.GetRealValue(v.Field(i))
		if vf.Kind() == reflect.Array || vf.Kind() == reflect.Slice {
			for j := 0; j < vf.Len(); j++ {
				if tip2 := SearchUsage(vf.Index(j), pos); tip2 != nil {
					return tip2
				}
			}
		} else {
			if tip2 := SearchUsage(vf, pos); tip2 != nil {
				return tip2
			}
		}
	}

	return tip
}

func getUsage(v reflect.Value, pos core.Position) *Tip {
	if v.Type().Implements(tiperType) && v.CanInterface() {
		if tip := v.Interface().(tiper); tip.contains(pos) {
			return tip.tip()
		}
		return nil
	} else if v.CanAddr() {
		if pv := v.Addr(); pv.Type().Implements(tiperType) && pv.CanInterface() {
			if tip := pv.Interface().(tiper); tip.contains(pos) {
				return tip.tip()
			}
			return nil
		}
	}
	return nil
}
