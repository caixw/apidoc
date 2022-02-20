// SPDX-License-Identifier: MIT

// Package mock 根据 doc 生成 mock 数据
package mock

import (
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/issue9/mux/v6"
	"github.com/issue9/qheader"
	"github.com/issue9/version"

	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/internal/ast"
	"github.com/caixw/apidoc/v7/internal/lexer"
	"github.com/caixw/apidoc/v7/internal/locale"
)

type mock struct {
	msgHandler *core.MessageHandler
	doc        *ast.APIDoc
	router     *mux.Router
	h          http.Handler
	servers    map[string]string
	indent     string
	gen        *GenOptions
}

// New 声明 Mock 对象
//
// msg 用于处理各类输出消息，仅在 ServeHTTP 中的消息才输出到 msg；
// d doc.APIDoc 实例，调用方需要保证该数据类型的正确性；
// indent 缩进字符串；
// servers 用于指定 d.Servers 中每一个服务对应的路由前缀；
// gen 生成随机数据的函数；
func New(msg *core.MessageHandler, d *ast.APIDoc, indent, imageURL string, servers map[string]string, gen *GenOptions) (http.Handler, error) {
	c, err := version.SemVerCompatible(d.APIDoc.V(), ast.Version)
	if err != nil {
		return nil, err
	}
	if !c {
		return nil, locale.NewError(locale.VersionInCompatible)
	}

	mu := mux.NewRouter("apidoc mock server", nil)
	router := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := w.Header()
		h.Set("Access-Control-Allow-Origin", "*")
		h.Set("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization")
		h.Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, HEAD, OPTIONS")
		mu.ServeHTTP(w, r)
	})

	m := &mock{
		msgHandler: msg,
		doc:        d,
		router:     mu,
		h:          router,
		indent:     indent,
		servers:    servers,
		gen:        gen,
	}

	if imageURL != "" {
		if imageURL[0] != '/' || imageURL[len(imageURL)-1] == '/' {
			panic("参数 imageURL 必须以 / 开头且不能以 / 结尾")
		}
		m.router.Get(imageURL+"/{path}", http.HandlerFunc(m.getImage))
	}

	m.parse()

	return m, nil
}

// Load 从本地或是远程加载文档内容
func Load(h *core.MessageHandler, path core.URI, indent, imageURL string, servers map[string]string, gen *GenOptions) (http.Handler, error) {
	data, err := path.ReadAll(nil)
	if err != nil {
		return nil, err
	}

	b := core.Block{Data: data, Location: core.Location{URI: path}}
	p, err := lexer.BlockEndPosition(b)
	if err != nil {
		return nil, err
	}
	b.Location.Range.End = p.Position

	// 加载并验证
	d := &ast.APIDoc{}
	d.Parse(h, b)
	return New(h, d, indent, imageURL, servers, gen)
}

func (m *mock) parse() {
	for _, api := range m.doc.APIs {
		handler := m.buildAPI(api)
		method := api.Method.V()
		path := api.Path.Path.V()

		if len(api.Servers) == 0 {
			m.router.Handle(path, handler, method)
			continue
		}

		for _, srv := range api.Servers {
			prefix, found := m.servers[srv.V()]
			if !found {
				prefix = "/" + srv.V()
			}
			m.router.Prefix(prefix).Handle(path, handler, method)
		}
	}

	routers := m.router.Routes()
	for path, methods := range routers {
		m.msgHandler.Locale(core.Info, locale.LoadAPI, "["+strings.Join(methods, ",")+"]", path)
	}
}

func (m *mock) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.h.ServeHTTP(w, r)
}

func (m *mock) getImage(w http.ResponseWriter, r *http.Request) {
	width, height := 500, 500
	var err error

	if ww := r.FormValue("width"); ww != "" {
		if width, err = strconv.Atoi(ww); err != nil {
			http.Error(w, locale.Sprintf(locale.ErrInvalidValue), http.StatusBadRequest)
			return
		}
	}

	if hh := r.FormValue("height"); hh != "" {
		if height, err = strconv.Atoi(hh); err != nil {
			http.Error(w, locale.Sprintf(locale.ErrInvalidValue), http.StatusBadRequest)
			return
		}
	}

	img := image.NewNRGBA(image.Rect(0, 0, width, height))
	headers := qheader.Accept(r)

	for _, h := range headers {
		switch strings.ToLower(h.Value) {
		case "image/jpeg":
			w.Header().Add("Content-Type", "image/jpeg")
			if err = jpeg.Encode(w, img, nil); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		case "image/png", "*/*":
			w.Header().Add("Content-Type", "image/png")
			if err = png.Encode(w, img); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		case "image/gif":
			w.Header().Add("Content-Type", "image/gif")
			if err = gif.Encode(w, img, nil); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
	}

	http.Error(w, http.StatusText(http.StatusNotAcceptable), http.StatusNotAcceptable)
}

func isValidRFC3339Date(val string) bool {
	_, err := time.Parse(ast.DateFormat, val)
	return err == nil
}

func isValidRFC3339Time(val string) bool {
	_, err := time.Parse(ast.TimeFormat, val)
	return err == nil
}

func isValidRFC3339DateTime(val string) bool {
	_, err := time.Parse(time.RFC3339, val)
	return err == nil
}
