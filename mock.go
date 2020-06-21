// SPDX-License-Identifier: MIT

package apidoc

import (
	"math/rand"
	"net/http"

	"github.com/issue9/rands"
	"github.com/issue9/utils"

	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/internal/ast"
	"github.com/caixw/apidoc/v7/internal/mock"
)

// Range 表示数值的范围
type Range struct {
	Min, Max int
}

// MockOptions mock 的一些随机设置项
type MockOptions struct {
	Indent    string            // 缩进字符串
	Servers   map[string]string // 为文档中所有 server 以及对应的路由前缀。
	SliceSize Range             // 指定用于生成数组大小范围的数值

	NumberSize  Range // 指定用于生成数值数据的范围
	EnableFloat bool  // 是否允许生成浮点数

	StringSize  Range  // 指定生成随机字符串的长度范围
	StringAlpha []byte // 指定生成字符串可用的字符

	URLDomains        []string // 指定生成 url 类型数据时可用的域名，默认为 example.com
	EmailDomains      []string // 指定生成 email 类型数据时可用的域名，默认为 example.com
	EmailUsernameSize Range    // 指定生成 email 类型数据的用户名长度范围，默认 [3,8]

	ImageBasePrefix string // 图片的基地址
}

var defaultMockOptions = &MockOptions{
	Indent:    "\t",
	SliceSize: Range{Min: 5, Max: 50},

	NumberSize:  Range{Min: 100, Max: 10000},
	EnableFloat: false,

	StringSize:  Range{Min: 50, Max: 1024},
	StringAlpha: rands.AlphaNumber,

	URLDomains:        []string{"https://example.com/"},
	EmailDomains:      []string{"example.com"},
	EmailUsernameSize: Range{Min: 3, Max: 8},

	ImageBasePrefix: "/__images__",
}

func (o *MockOptions) gen() *mock.GenOptions {
	return &mock.GenOptions{
		Number: func(p *ast.Param) interface{} {
			switch p.Type.V() {
			case ast.TypeFloat:
				return o.float()
			case ast.TypeInt:
				return o.integer()
			}

			if !o.EnableFloat {
				return o.integer()
			}

			if rand.Int()%2 == 0 {
				return o.integer()
			}
			return o.float()
		},

		String: func(p *ast.Param) string {
			switch p.Type.V() {
			case ast.TypeEmail:
				return o.email()
			case ast.TypeURL:
				return o.url()
			case ast.TypeImage:
				return o.image()
			}
			return rands.String(o.StringSize.Min, o.StringSize.Max, o.StringAlpha)
		},

		Bool: func() bool {
			return rand.Int()%2 == 0
		},

		SliceSize: func() int {
			return rand.Intn(o.SliceSize.Max-o.SliceSize.Min) + o.SliceSize.Min
		},

		Index: func(max int) int {
			return rand.Intn(max)
		},
	}
}

func (o *MockOptions) integer() int {
	return rand.Intn(o.NumberSize.Max-o.NumberSize.Min) + o.NumberSize.Min
}

func (o *MockOptions) float() float32 {
	return float32(o.NumberSize.Min) + rand.Float32()*float32(o.NumberSize.Max-o.NumberSize.Min)
}

func (o *MockOptions) url() string {
	url := o.URLDomains[rand.Intn(len(o.URLDomains))]
	if url[len(url)-1] != '/' {
		url += "/"
	}

	size := rand.Intn(4)
	for i := 0; i < size; i++ {
		url += rands.String(1, 5, rands.AlphaNumber) + "/"
	}
	return url
}

func (o *MockOptions) email() string {
	domain := o.EmailDomains[rand.Intn(len(o.EmailDomains))]
	username := rands.String(o.EmailUsernameSize.Min, o.EmailUsernameSize.Max, rands.AlphaNumber)
	return username + "@" + domain
}

func (o *MockOptions) image() string {
	path := o.ImageBasePrefix
	if path[len(path)-1] != '/' {
		path += "/"
	}
	return path + rands.String(1, 5, rands.AlphaNumber)
}

// Mock 生成 Mock 中间件
//
// data 为文档内容；
// options 用于生成 Mock 数据的随机项，如果为 nil，则会使用一些默认值；
func Mock(h *core.MessageHandler, data []byte, o *MockOptions) (http.Handler, error) {
	o, err := mergeMockOptions(o)
	if err != nil {
		return nil, err
	}

	d := &ast.APIDoc{}
	d.Parse(h, core.Block{Data: data})
	return mock.New(h, d, o.Indent, o.ImageBasePrefix, o.Servers, o.gen())
}

// MockFile 生成 Mock 中间件
//
// path 为文档路径，可以是本地路径也可以是 URL，根据是否为 http 或是 https 开头做判断；
// options 用于生成 Mock 数据的随机项，如果为 nil，则会使用一些默认值；
func MockFile(h *core.MessageHandler, path core.URI, o *MockOptions) (http.Handler, error) {
	o, err := mergeMockOptions(o)
	if err != nil {
		return nil, err
	}

	return mock.Load(h, path, o.Indent, o.ImageBasePrefix, o.Servers, o.gen())
}

func mergeMockOptions(options *MockOptions) (*MockOptions, error) {
	if options == nil {
		return defaultMockOptions, nil
	}

	ret := &MockOptions{}
	*ret = *defaultMockOptions
	err := utils.Merge(false, ret, options)
	return ret, err
}
