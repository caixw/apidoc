// SPDX-License-Identifier: MIT

package apidoc

import (
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/issue9/assert"
	"github.com/issue9/assert/rest"
	"github.com/issue9/is"

	"github.com/caixw/apidoc/v7/core/messagetest"
	"github.com/caixw/apidoc/v7/internal/ast"
	"github.com/caixw/apidoc/v7/internal/ast/asttest"
	"github.com/caixw/apidoc/v7/internal/token"
)

func TestMockOptions_url(t *testing.T) {
	a := assert.New(t)

	o := &MockOptions{
		URLDomains: []string{"https://apidoc.tools/"},
	}
	url := o.url()
	a.True(strings.HasPrefix(url, o.URLDomains[0])).
		True(is.URL(url))

	o.URLDomains[0] = "https://apidoc.tools"
	url = o.url()
	a.True(strings.HasPrefix(url, o.URLDomains[0])).
		True(is.URL(url))
}

func TestMockOptions_email(t *testing.T) {
	a := assert.New(t)

	o := &MockOptions{
		EmailDomains:      []string{"apidoc.tools"},
		EmailUsernameSize: Range{Min: 5, Max: 11},
	}
	email := o.email()
	a.True(strings.HasSuffix(email, o.EmailDomains[0])).
		True(is.Email(email))
	index := strings.IndexByte(email, '@')
	username := email[:index]
	a.True(len(username) >= o.EmailUsernameSize.Min).
		True(len(username) <= o.EmailUsernameSize.Max)
}

func TestMockOptions_gen(t *testing.T) {
	a := assert.New(t)

	g, err := defaultMockOptions.gen()
	a.NotError(err)
	a.NotNil(g)
	count := 500

	for i := 0; i < count; i++ {
		size := g.SliceSize()
		a.True(size >= defaultMockOptions.SliceSize.Min).
			True(size <= defaultMockOptions.SliceSize.Max)
	}

	for i := 0; i < count; i++ {
		size := g.Index(500)
		a.True(size <= 500).
			True(size >= 0)
	}

	// String
	for i := 0; i < count; i++ {
		str := g.String(&ast.Param{})
		a.True(len(str) >= defaultMockOptions.StringSize.Min).
			True(len(str) <= defaultMockOptions.StringSize.Max)
	}

	// String.Email
	for i := 0; i < count; i++ {
		str := g.String(&ast.Param{Type: &ast.TypeAttribute{Value: token.String{Value: "string.email"}}})
		a.True(is.Email(str))
	}

	// String.URL
	for i := 0; i < count; i++ {
		str := g.String(&ast.Param{Type: &ast.TypeAttribute{Value: token.String{Value: "string.url"}}})
		a.True(is.URL(str))
	}

	// String.Image
	for i := 0; i < count; i++ {
		str := g.String(&ast.Param{Type: &ast.TypeAttribute{Value: token.String{Value: "string.image"}}})
		a.True(strings.HasPrefix(str, defaultMockOptions.ImageBasePrefix))
	}

	// String.Date
	for i := 0; i < count; i++ {
		str := g.String(&ast.Param{Type: &ast.TypeAttribute{Value: token.String{Value: "string.date"}}})
		t, err := time.Parse(ast.DateFormat, str)
		a.NotError(err).
			True(t.After(defaultMockOptions.DateStart)).
			True(t.Before(defaultMockOptions.DateEnd))
	}

	// String.DateTime
	for i := 0; i < count; i++ {
		str := g.String(&ast.Param{Type: &ast.TypeAttribute{Value: token.String{Value: "string.date-time"}}})
		t, err := time.Parse(time.RFC3339, str)
		a.NotError(err).
			True(t.After(defaultMockOptions.DateStart)).
			True(t.Before(defaultMockOptions.DateEnd))
	}

	// String.Time
	for i := 0; i < count; i++ {
		str := g.String(&ast.Param{Type: &ast.TypeAttribute{Value: token.String{Value: "string.time"}}})
		_, err := time.Parse(ast.TimeFormat, str)
		a.NotError(err)
	}

	// Number
	defaultMockOptions.EnableFloat = false
	g, err = defaultMockOptions.gen()
	a.NotError(err)
	for i := 0; i < count; i++ {
		numInterface := g.Number(&ast.Param{})
		num, ok := numInterface.(int)
		a.True(ok).
			True(num >= defaultMockOptions.NumberSize.Min).
			True(num <= defaultMockOptions.NumberSize.Max)
	}

	// Number enableFloat
	defaultMockOptions.EnableFloat = true
	g, err = defaultMockOptions.gen()
	a.NotError(err)
	for i := 0; i < count; i++ {
		numInterface := g.Number(&ast.Param{})
		num, ok := numInterface.(int)
		if ok {
			a.True(num >= defaultMockOptions.NumberSize.Min).
				True(num <= defaultMockOptions.NumberSize.Max)
		} else {
			f, ok := numInterface.(float32)
			a.True(ok).
				True(f >= float32(defaultMockOptions.NumberSize.Min)).
				True(f <= float32(defaultMockOptions.NumberSize.Max), "%f,%d", f, defaultMockOptions.NumberSize.Max)
		}
	}

	// Number.int
	defaultMockOptions.EnableFloat = false
	g, err = defaultMockOptions.gen()
	a.NotError(err)
	for i := 0; i < count; i++ {
		numInterface := g.Number(&ast.Param{Type: &ast.TypeAttribute{Value: token.String{Value: "number.int"}}})
		num, ok := numInterface.(int)
		a.True(ok).
			True(num >= defaultMockOptions.NumberSize.Min).
			True(num <= defaultMockOptions.NumberSize.Max)
	}

	// Number.float
	defaultMockOptions.EnableFloat = false
	g, err = defaultMockOptions.gen()
	a.NotError(err)
	for i := 0; i < count; i++ {
		numInterface := g.Number(&ast.Param{Type: &ast.TypeAttribute{Value: token.String{Value: "number.float"}}})
		num, ok := numInterface.(float32)
		a.True(ok).
			True(num >= float32(defaultMockOptions.NumberSize.Min)).
			True(num <= float32(defaultMockOptions.NumberSize.Max))
	}
}

func TestMock(t *testing.T) {
	a := assert.New(t)

	rslt := messagetest.NewMessageHandler()
	opt := &MockOptions{}
	*opt = *defaultMockOptions
	opt.Servers = map[string]string{"admin": "/admin"}
	mock, err := Mock(rslt.Handler, asttest.XML(a), opt)
	a.NotError(err).NotNil(mock)
	srv := rest.NewServer(t, mock, nil)

	srv.Get("/admin/users").
		Header("authorization", "xxx").
		Header("content-type", "application/json").
		Header("Accept", "application/json").
		Do().Status(http.StatusOK)

	// 不存在 client
	srv.Get("/client/users").Do().Status(http.StatusNotFound)

	srv.Post("/admin/users", nil).Do().Status(http.StatusBadRequest)    // 未指定报头
	srv.Delete("/admin/users").Do().Status(http.StatusMethodNotAllowed) // 不存在

	rslt.Handler.Stop()
	srv.Close()
}

func TestMockFile(t *testing.T) {
	a := assert.New(t)

	rslt := messagetest.NewMessageHandler()
	opt := &MockOptions{}
	*opt = *defaultMockOptions
	opt.Servers = map[string]string{"admin": "/admin"}
	mock, err := MockFile(rslt.Handler, asttest.URI(a), opt)
	a.NotError(err).NotNil(mock)
	srv := rest.NewServer(t, mock, nil)

	srv.Get("/admin/users").
		Header("authorization", "xxx").
		Header("content-type", "application/json").
		Header("Accept", "application/json").
		Do().Status(http.StatusOK)

	// 不存在 client
	srv.Get("/client/users").Do().Status(http.StatusNotFound)

	srv.Post("/admin/users", nil).Do().Status(http.StatusBadRequest)    // 未指定报头
	srv.Delete("/admin/users").Do().Status(http.StatusMethodNotAllowed) // 不存在

	rslt.Handler.Stop()
	srv.Close()

	// 测试多个 servers 值
	rslt = messagetest.NewMessageHandler()
	*opt = *defaultMockOptions
	opt.Servers = map[string]string{"admin": "/admin", "client": "/c"}
	mock, err = MockFile(rslt.Handler, asttest.URI(a), opt)
	a.NotError(err).NotNil(mock)
	srv = rest.NewServer(t, mock, nil)

	srv.Post("/admin/users", []byte(`{"id":1,"name":"name"}`)).
		Header("authorization", "xxx").
		Header("content-type", "application/json").
		Header("Accept", "application/json").
		Do().Status(http.StatusCreated)

	srv.Post("/c/users", []byte(`{"id":1,"name":"name"}`)).
		Header("authorization", "xxx").
		Header("content-type", "application/json").
		Header("Accept", "application/json").
		Do().Status(http.StatusCreated)

	srv.Post("/admin/users", nil).Do().Status(http.StatusBadRequest)    // 未指定报头
	srv.Delete("/admin/users").Do().Status(http.StatusMethodNotAllowed) // 不存在

	srv.Get("/c/users").Do().Status(http.StatusMethodNotAllowed)    // POST /users 未指定 client
	srv.Delete("/c/users").Do().Status(http.StatusMethodNotAllowed) // 不存在

	rslt.Handler.Stop()
	srv.Close()
}

func TestMergeMockOptions(t *testing.T) {
	a := assert.New(t)

	o := &MockOptions{SliceSize: Range{Min: 0, Max: 5}}
	err := o.mergeMockOptions()
	a.NotError(err).
		Equal(o.EmailDomains, []string{"example.com"}).
		Equal(o.SliceSize, Range{Min: 0, Max: 5})
}
