// SPDX-License-Identifier: MIT

package apidoc

import (
	"net/http"
	"testing"

	"github.com/issue9/assert"
	"github.com/issue9/assert/rest"

	"github.com/caixw/apidoc/v7/core/messagetest"
	"github.com/caixw/apidoc/v7/internal/ast/asttest"
)

func TestMockOptions_gen(t *testing.T) {
	a := assert.New(t)

	g := defaultMockOptions.gen()
	a.NotNil(g)
	count := 500

	for i := 0; i < count; i++ {
		size := g.SliceSize()
		a.True(size >= defaultMockOptions.MinSliceSize).
			True(size <= defaultMockOptions.MaxSliceSize)
	}

	for i := 0; i < count; i++ {
		size := g.Index(500)
		a.True(size <= 500).
			True(size >= 0)
	}

	for i := 0; i < count; i++ {
		str := g.String()
		a.True(len(str) >= defaultMockOptions.MinString).
			True(len(str) <= defaultMockOptions.MaxString)
	}

	defaultMockOptions.EnableFloat = false
	g = defaultMockOptions.gen()
	for i := 0; i < count; i++ {
		numInterface := g.Number()
		num, ok := numInterface.(int)
		a.True(ok).
			True(num >= defaultMockOptions.MinNumber).
			True(num <= defaultMockOptions.MaxNumber)
	}

	defaultMockOptions.EnableFloat = true
	g = defaultMockOptions.gen()
	for i := 0; i < count; i++ {
		numInterface := g.Number()
		num, ok := numInterface.(int)
		if ok {
			a.True(num >= defaultMockOptions.MinNumber).
				True(num <= defaultMockOptions.MaxNumber)
		} else {
			f, ok := numInterface.(float32)
			a.True(ok).
				True(f >= float32(defaultMockOptions.MinNumber)).
				True(f <= float32(defaultMockOptions.MaxNumber), "%f,%d", f, defaultMockOptions.MaxNumber)
		}
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

	srv.Get("/admin/users").
		Header("authorization", "xxx").
		Header("content-type", "application/json").
		Header("Accept", "application/json").
		Do().Status(http.StatusOK)

	srv.Get("/c/users").
		Header("authorization", "xxx").
		Header("content-type", "application/json").
		Header("Accept", "application/json").
		Do().Status(http.StatusOK)

	srv.Post("/admin/users", nil).Do().Status(http.StatusBadRequest)    // 未指定报头
	srv.Delete("/admin/users").Do().Status(http.StatusMethodNotAllowed) // 不存在

	srv.Post("/c/users", nil).Do().Status(http.StatusMethodNotAllowed) // POST /users 未指定 client
	srv.Delete("/c/users").Do().Status(http.StatusMethodNotAllowed)    // 不存在

	rslt.Handler.Stop()
	srv.Close()
}
