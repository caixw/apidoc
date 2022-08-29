// SPDX-License-Identifier: MIT

package mock

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/issue9/qheader"
	"github.com/issue9/validation/is"

	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/internal/ast"
	"github.com/caixw/apidoc/v7/internal/locale"
)

func (m *mock) buildAPI(api *ast.API) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.msgHandler.Locale(core.Succ, locale.RequestAPI, r.Method, r.URL.Path)
		if api.Deprecated != nil {
			m.msgHandler.Locale(core.Warn, locale.DeprecatedWarn, r.Method, r.URL.Path, api.Deprecated.V())
		}

		if err := validQueries(api.Path.Queries, r); err != nil {
			m.handleError(w, r, "", err)
			return
		}

		for _, header := range api.Headers {
			field := "headers[" + header.Name.V() + "]"
			if err := validSimpleParam(header, field, r.Header.Get(header.Name.V())); err != nil {
				m.handleError(w, r, field, err)
				return
			}
		}

		if len(api.Requests) > 0 { // GET、OPTIONS 之类的可能没有 body
			if err := validRequest(m.doc.XMLNamespaces, api.Requests, r); err != nil {
				m.handleError(w, r, "request.body.", err)
				return
			}
		}

		m.renderResponse(api, w, r)
	})
}

func validRequest(ns []*ast.XMLNamespace, requests []*ast.Request, r *http.Request) error {
	ct := r.Header.Get("Content-Type")
	if ct == "" || ct == "*/*" || strings.HasSuffix(ct, "/*") { // 用户提交的 content-type 必须是明确的值
		return core.NewError(locale.ErrInvalidValue).WithField("headers[content-type]")
	}
	req := findRequestByContentType(requests, ct)
	if req == nil {
		return core.NewError(locale.ErrInvalidValue).WithField("headers[content-type]")
	}

	for _, header := range req.Headers {
		if err := validSimpleParam(header, "headers["+header.Name.V()+"]", r.Header.Get(header.Name.V())); err != nil {
			return err
		}
	}

	content, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	if err = r.Body.Close(); err != nil {
		return err
	}

	switch ct {
	case "application/json":
		return validJSON(req, content)
	case "application/xml", "text/xml":
		return validXML(ns, req, content)
	default:
		return core.NewError(locale.ErrInvalidValue).WithField("headers[content-type]")
	}
}

func (m *mock) renderResponse(api *ast.API, w http.ResponseWriter, r *http.Request) {
	accepts := qheader.Accept(r)

	resp, accept := findResponseByAccept(m.doc.Mimetypes, api.Responses, accepts.Items)
	if resp == nil {
		// 仅在 api.Responses 无法匹配任何内容的时候，才从 doc.Responses 中查找内容
		resp, accept = findResponseByAccept(m.doc.Mimetypes, m.doc.Responses, accepts.Items)
		if resp == nil {
			m.handleError(w, r, "headers[Accept]", locale.NewError(locale.ErrInvalidValue))
			return
		}
	}

	data, err := m.buildResponse(resp, r)
	if err != nil {
		m.handleError(w, r, "response.body.", err)
		return
	}

	w.Header().Set("Content-Type", accept)
	w.Header().Set("Server", core.Name)
	for _, item := range resp.Headers {
		switch primitive, _ := ast.ParseType(item.Type.V()); primitive {
		case ast.TypeBool:
			w.Header().Set(item.Name.V(), strconv.FormatBool(m.gen.generateBool()))
		case ast.TypeNumber:
			w.Header().Set(item.Name.V(), fmt.Sprint(m.gen.generateNumber(item)))
		case ast.TypeString:
			w.Header().Set(item.Name.V(), m.gen.generateString(item))
		default:
			m.handleError(w, r, "response.headers", locale.NewError(locale.ErrInvalidFormat))
			return
		}
	}

	w.WriteHeader(resp.Status.V())
	if _, err := w.Write(data); err != nil {
		m.msgHandler.Error(err) // 此时状态码已经输出
	}
}

// 需要保证 ct 的值不能为空
func findRequestByContentType(requests []*ast.Request, ct string) *ast.Request {
	var none *ast.Request
	for _, req := range requests {
		if req.Mimetype.V() == ct {
			return req
		} else if none == nil && req.Mimetype.V() == "" {
			none = req
		}
	}

	if none != nil {
		return none
	}
	return nil
}

// accepts 必须是已经按权重进行排序的。
func findResponseByAccept(mimetypes []*ast.Element, requests []*ast.Request, accepts []*qheader.Item) (*ast.Request, string) {
	if len(requests) == 0 {
		return nil, ""
	}

	var none *ast.Request // 表示 requests 中 mimetype 值为空的第一个子项

	// 从 requests 中查找是否有符合 accepts 的内容
	for _, req := range requests {
		if none == nil && req.Mimetype.V() == "" {
			none = req
		}
		if req.Mimetype.V() != "" && matchContentType(req.Mimetype.V(), accepts) {
			return req, req.Mimetype.V()
		}
	}

	// 如果存在 none，则从 doc.mimetypes 中查找是否有与 none.Mimetype 相匹配的
	if none != nil {
		for _, item := range mimetypes {
			mt := item.Content.Value
			if mt != "" && matchContentType(mt, accepts) {
				return none, mt
			}
		}
	}

	return nil, ""
}

// 查看 ct 是否有与 accepts 相匹配的项，必须保证 ct 的值不为空。
func matchContentType(ct string, accepts []*qheader.Item) bool {
	for _, a := range accepts {
		if (strings.HasSuffix(a.Value, "/*") && strings.HasPrefix(ct, a.Value[:len(a.Value)-1])) ||
			ct == a.Value ||
			a.Value == "*/*" {
			return true
		}
	}
	return false
}

// 处理 serveHTTP 中的错误
func (m *mock) handleError(w http.ResponseWriter, r *http.Request, field string, err error) {
	// 这并不是一个真实存在的 URI
	file := core.URI(r.Method + ": " + r.URL.Path)

	if serr, ok := err.(*core.Error); ok {
		if serr.Field == "" {
			serr.Field = field
		} else {
			serr.Field = field + serr.Field
		}

		serr.Location.URI = file
	} else {
		err = (core.Location{URI: file}).WithError(err).WithField(field)
	}

	m.msgHandler.Error(err)
	w.WriteHeader(http.StatusBadRequest)
}

func validQueries(queries []*ast.Param, r *http.Request) error {
	for _, query := range queries {
		field := "queries[" + query.Name.V() + "]."

		valid := func(p *ast.Param, v string) error {
			err := validSimpleParam(p, field, v)
			if serr, ok := err.(*core.Error); ok {
				serr.Field = field + serr.Field
			}
			return err
		}

		if !query.Array.V() {
			if err := valid(query, r.FormValue(query.Name.V())); err != nil {
				return err
			}
		} else if !query.ArrayStyle.V() { // 默认的 form 格式
			if err := r.ParseForm(); err != nil {
				return err
			}
			for _, v := range r.Form[query.Name.V()] {
				if err := valid(query, v); err != nil {
					return err
				}
			}
		} else {
			values := strings.Split(r.FormValue(query.Name.V()), ",")
			for _, v := range values {
				if err := valid(query, v); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// 验证单个参数，仅支持对 query、header 等简单类型的参数验证
func validSimpleParam(p *ast.Param, name, val string) error {
	if p == nil {
		return nil
	}

	if val == "" && p.Type.V() != ast.TypeString { // 字符串的默认值可以为 “”
		if (p.Optional != nil && p.Optional.V()) ||
			(p.Default != nil && p.Default.V() != "") {
			return nil
		}
		return core.NewError(locale.ErrIsEmpty, name)
	}

	switch p.Type.V() {
	case ast.TypeBool:
		if _, err := strconv.ParseBool(val); err != nil {
			return core.NewError(locale.ErrInvalidFormat)
		}
	case ast.TypeNumber:
		if !is.Number(val) {
			return core.NewError(locale.ErrInvalidFormat)
		}
	case ast.TypeString:
	case ast.TypeObject:
	case ast.TypeNone:
		if val != "" {
			return core.NewError(locale.ErrInvalidValue)
		}
	}

	if isEnum(p) {
		found := false
		for _, e := range p.Enums {
			if e.Value.V() == val {
				found = true
				break
			}
		}

		if !found {
			return core.NewError(locale.ErrInvalidValue)
		}
	}

	return nil
}

func (m *mock) buildResponse(p *ast.Request, r *http.Request) ([]byte, error) {
	if p == nil {
		return nil, nil
	}

	for _, header := range p.Headers {
		if err := validSimpleParam(header, "headers["+header.Name.V()+"]", r.Header.Get(header.Name.V())); err != nil {
			return nil, err
		}
	}

	headers := qheader.Accept(r)
	if headers != nil {
		for _, h := range headers.Items {
			switch strings.ToLower(h.Value) {
			case "application/json", "*/*":
				return buildJSON(p, m.indent, m.gen)
			case "application/xml", "text/xml":
				return buildXML(m.doc.XMLNamespaces, p, m.indent, m.gen)
			}
		}
	}
	return nil, core.NewError(locale.ErrInvalidValue).WithField("headers[accept]")
}
