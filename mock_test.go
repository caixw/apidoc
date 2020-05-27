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
