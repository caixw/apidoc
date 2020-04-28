// SPDX-License-Identifier: MIT

package apidoc

import (
	"net/http"
	"testing"

	"github.com/issue9/assert"
	"github.com/issue9/assert/rest"
	"github.com/issue9/version"

	"github.com/caixw/apidoc/v6/core"
	"github.com/caixw/apidoc/v6/core/messagetest"
	"github.com/caixw/apidoc/v6/internal/ast/asttest"
	"github.com/caixw/apidoc/v6/internal/docs"
	"github.com/caixw/apidoc/v6/internal/vars"
)

func TestVersion(t *testing.T) {
	a := assert.New(t)

	a.Equal(Version(), vars.Version())
	a.True(version.SemVerValid(Version()))
}

func TestValid(t *testing.T) {
	a := assert.New(t)

	data, err := asttest.URI(a).ReadAll(nil)
	a.NotError(err).NotNil(data)
	a.NotError(Valid(core.Block{Data: data}))
}

func TestStatic(t *testing.T) {
	srv := rest.NewServer(t, Static(docs.Dir(), false), nil)
	defer srv.Close()

	srv.Get("/icon.svg").Do().Status(http.StatusOK)
}

func TestView(t *testing.T) {
	a := assert.New(t)

	data := asttest.XML(a)
	h := View(http.StatusCreated, "/test/apidoc.xml", data, "text/xml", "", false)
	srv := rest.NewServer(t, h, nil)
	srv.Get("/test/apidoc.xml").Do().
		Status(http.StatusCreated).
		Header("content-type", "text/xml")

	srv.Get("/index.xml").Do().
		Status(http.StatusOK)

	srv.Get("/v6/apidoc.xsl").Do().
		Status(http.StatusOK)

	srv.Close()

	// 能正确覆盖 Static 中的 index.xml
	h = View(http.StatusCreated, "/index.xml", data, "text/css", "", false)
	srv = rest.NewServer(t, h, nil)
	srv.Get("/index.xml").Do().
		Status(http.StatusCreated).
		Header("content-type", "text/css")

	srv.Get("/v6/apidoc.xsl").Do().
		Status(http.StatusOK)

	srv.Close()
}

func TestViewFile(t *testing.T) {
	a := assert.New(t)

	h, err := ViewFile(http.StatusAccepted, "/apidoc.xml", asttest.URI(a), "text/xml", "", false)
	a.NotError(err).NotNil(h)
	srv := rest.NewServer(t, h, nil)
	srv.Get("/apidoc.xml").Do().
		Status(http.StatusAccepted).
		Header("content-type", "text/xml")
	srv.Close()

	h, err = ViewFile(http.StatusAccepted, "/apidoc.xml", asttest.URI(a), "", "", false)
	a.NotError(err).NotNil(h)
	srv = rest.NewServer(t, h, nil)
	srv.Get("/apidoc.xml").Do().
		Status(http.StatusAccepted)
	srv.Close()

	// 覆盖现有的 index.xml
	h, err = ViewFile(http.StatusAccepted, "", asttest.URI(a), "", "", false)
	a.NotError(err).NotNil(h)
	srv = rest.NewServer(t, h, nil)
	srv.Get("/index.xml").Do().
		Status(http.StatusAccepted)
	srv.Close()
}

func TestMock(t *testing.T) {
	a := assert.New(t)

	_, _, h := messagetest.MessageHandler()
	mock, err := Mock(h, asttest.XML(a), map[string]string{"admin": "/admin"})
	a.NotError(err).NotNil(h)
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

	h.Stop()
	srv.Close()

}

func TestMockFile(t *testing.T) {
	a := assert.New(t)

	_, _, h := messagetest.MessageHandler()
	mock, err := MockFile(h, asttest.URI(a), map[string]string{"admin": "/admin"})
	a.NotError(err).NotNil(h)
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

	h.Stop()
	srv.Close()

	// 测试多个 servers 值
	_, _, h = messagetest.MessageHandler()
	mock, err = MockFile(h, asttest.URI(a), map[string]string{"admin": "/admin", "client": "/c"})
	a.NotError(err).NotNil(h)
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

	h.Stop()
	srv.Close()
}

func TestAddStylesheet(t *testing.T) {
	a := assert.New(t)

	data := []*struct {
		input  string
		output string
	}{
		{
			input: "",
			output: `
<?xml-stylesheet type="text/xsl" href="./v6/apidoc.xsl"?>`,
		},
		{
			input: `<?xml version="1.0"?>`,
			output: `<?xml version="1.0"?>
<?xml-stylesheet type="text/xsl" href="./v6/apidoc.xsl"?>`,
		},
		{
			input: `<?xml version="1.0"?>
<?xml-stylesheet href="xxx"?>`,
			output: `<?xml version="1.0"?>
<?xml-stylesheet type="text/xsl" href="./v6/apidoc.xsl"?>
<?xml-stylesheet href="xxx"?>`,
		},
	}

	for index, item := range data {
		output := string(addStylesheet([]byte(item.input)))
		a.Equal(output, item.output, "not equal at %d,want: %s,get:%s", index, item.output, output)
	}
}
