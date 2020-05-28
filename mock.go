// SPDX-License-Identifier: MIT

package apidoc

import (
	"math/rand"
	"net/http"

	"github.com/issue9/rands"

	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/internal/ast"
	"github.com/caixw/apidoc/v7/internal/mock"
)

// MockOptions mock 的一些随机设置项
type MockOptions struct {
	Indent string // 缩进字符串

	// 为文档中所有 server 以及对应的路由前缀。
	Servers map[string]string

	// 生成介于 [MinSliceSize, MaxSliceSize] 之间的数值，
	// 该数仁用于生成 mock 数据中的数组。
	MaxSliceSize int
	MinSliceSize int

	// 生成介于 [MinNumber, MaxNumber] 之间的数值，
	// 该数值用于生成 mock 数据中的 number 类型的值。
	// EnableFloat 表示是否可以返回浮点数。
	MaxNumber   int
	MinNumber   int
	EnableFloat bool

	// 生成介于 [MinString, MaxString] 之间的数值，
	// 该数值用于生成 mock 数据中的字符串类型的值。
	// StringAlpha 表示用于生成字符串的可用字符。
	MaxString   int
	MinString   int
	StringAlpha []byte
}

var defaultMockOptions = &MockOptions{
	Indent: "\t",

	MaxSliceSize: 50,
	MinSliceSize: 5,

	MaxNumber:   10000,
	MinNumber:   100,
	EnableFloat: false,

	MaxString:   1024,
	MinString:   50,
	StringAlpha: rands.AlphaNumber,
}

func (o *MockOptions) gen() *mock.GenOptions {
	return &mock.GenOptions{
		Number: func() interface{} {
			if !o.EnableFloat {
				return rand.Intn(o.MaxNumber-o.MinNumber) + o.MinNumber
			}

			if rand.Int()%2 == 0 {
				return rand.Intn(o.MaxNumber-o.MinNumber) + o.MinNumber
			}
			return float32(o.MinNumber) + rand.Float32()*float32(o.MaxNumber-o.MinNumber)
		},

		String: func() string {
			return rands.String(o.MinString, o.MaxString, o.StringAlpha)
		},

		Bool: func() bool {
			return rand.Int()%2 == 0
		},

		SliceSize: func() int {
			return rand.Intn(o.MaxSliceSize-o.MinSliceSize) + o.MinSliceSize
		},

		Index: func(max int) int {
			return rand.Intn(max)
		},
	}
}

// Mock 生成 Mock 中间件
//
// data 为文档内容；
// options 用于生成 Mock 数据的随机项，如果为 nil，则会使用一些默认值；
func Mock(h *core.MessageHandler, data []byte, options *MockOptions) (http.Handler, error) {
	d := &ast.APIDoc{}
	d.Parse(h, core.Block{Data: data})
	if options == nil {
		options = defaultMockOptions
	}
	return mock.New(h, d, options.Indent, options.Servers, options.gen())
}

// MockFile 生成 Mock 中间件
//
// path 为文档路径，可以是本地路径也可以是 URL，根据是否为 http 或是 https 开头做判断；
// options 用于生成 Mock 数据的随机项，如果为 nil，则会使用一些默认值；
func MockFile(h *core.MessageHandler, path core.URI, options *MockOptions) (http.Handler, error) {
	if options == nil {
		options = defaultMockOptions
	}
	return mock.Load(h, path, options.Indent, options.Servers, options.gen())
}
