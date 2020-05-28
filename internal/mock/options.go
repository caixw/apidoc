// SPDX-License-Identifier: MIT

package mock

import (
	"strconv"

	"github.com/caixw/apidoc/v7/internal/ast"
)

// GenOptions 生成随机数据的选项
type GenOptions struct {
	Number    func() interface{}
	String    func() string
	Bool      func() bool
	SliceSize func() int
	Index     func(max int) int
}

func isEnum(p *ast.Param) bool {
	return len(p.Enums) > 0
}

func (g *GenOptions) generateBool() bool {
	return g.Bool()
}

func (g *GenOptions) generateNumber(p *ast.Param) interface{} {
	if isEnum(p) {
		index := g.Index(len(p.Enums))
		v, err := strconv.ParseInt(p.Enums[index].Value.V(), 10, 32)
		if err != nil { // 这属于文档定义错误，直接 panic
			panic(err)
		}
		return v
	}
	return g.Number()
}

func (g *GenOptions) generateString(p *ast.Param) string {
	if isEnum(p) {
		return p.Enums[g.Index(len(p.Enums))].Value.V()
	}
	return g.String()
}

// 生成随机的数组长度
func (g *GenOptions) generateSliceSize() int {
	return g.SliceSize()
}
