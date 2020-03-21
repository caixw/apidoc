// SPDX-License-Identifier: MIT

package mock

import (
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/issue9/is"
	"github.com/issue9/qheader"

	"github.com/caixw/apidoc/v6/core"
	"github.com/caixw/apidoc/v6/internal/locale"
	"github.com/caixw/apidoc/v6/internal/vars"
	"github.com/caixw/apidoc/v6/spec"
)

func (m *Mock) buildAPI(api *spec.API) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.h.Message(core.Succ, locale.RequestAPI, r.Method, r.URL.Path)
		if api.Deprecated != "" {
			m.h.Message(core.Warn, locale.DeprecatedWarn, r.Method, r.URL.Path, api.Deprecated)
		}

		if err := validQueryArrayParam(api.Path.Queries, r); err != nil {
			m.handleError(w, r, "", err)
			return
		}

		for _, header := range api.Headers {
			if err := validSimpleParam(header, r.Header.Get(header.Name)); err != nil {
				m.handleError(w, r, "headers["+header.Name+"]", err)
				return
			}
		}

		if len(api.Requests) > 0 { // GET、OPTIONS 之类的可能没有 body
			if err := validRequest(api.Requests, r); err != nil {
				m.handleError(w, r, "request.body.", err)
				return
			}
		}

		m.renderResponse(api, w, r)
	})
}

func validRequest(requests []*spec.Request, r *http.Request) error {
	ct := r.Header.Get("Content-Type")
	if ct == "" || ct == "*/*" || strings.HasSuffix(ct, "/*") { // 用户提交的 content-type 必须是明确的值
		return core.NewLocaleError(core.Location{}, "headers[content-type]", locale.ErrInvalidValue)
	}
	req := findRequestByContentType(requests, ct)
	if req == nil {
		return core.NewLocaleError(core.Location{}, "headers[content-type]", locale.ErrInvalidValue)
	}

	for _, header := range req.Headers {
		if err := validSimpleParam(header, r.Header.Get(header.Name)); err != nil {
			return err
		}
	}

	content, err := ioutil.ReadAll(r.Body)
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
		return validXML(req, content)
	default:
		return core.NewLocaleError(core.Location{}, "headers[content-type]", locale.ErrInvalidValue)
	}
}

func (m *Mock) renderResponse(api *spec.API, w http.ResponseWriter, r *http.Request) {
	accepts, err := qheader.Accept(r)
	if err != nil {
		m.handleError(w, r, "request.headers[Accept]", err)
		return
	}

	resp, accept := findResponseByAccept(m.doc.Mimetypes, api.Responses, accepts)
	if resp == nil {
		// 仅在 api.Responses 无法匹配任何内容的时候，才从 doc.Responses 中查找内容
		resp, accept = findResponseByAccept(m.doc.Mimetypes, m.doc.Responses, accepts)
		if resp == nil {
			m.handleError(w, r, "headers[Accept]", locale.Errorf(locale.ErrInvalidValue))
			return
		}
	}

	data, err := buildResponse(resp, r)
	if err != nil {
		m.handleError(w, r, "response.body.", err)
		return
	}

	w.Header().Set("Content-Type", accept)
	w.Header().Set("Server", vars.Name)
	for _, item := range resp.Headers {
		switch item.Type {
		case spec.Bool:
			w.Header().Set(item.Name, strconv.FormatBool(generateBool()))
		case spec.Number:
			w.Header().Set(item.Name, strconv.FormatInt(generateNumber(item), 10))
		case spec.String:
			w.Header().Set(item.Name, generateString(item))
		default:
			m.handleError(w, r, "response.headers", locale.Errorf(locale.ErrInvalidFormat))
			return
		}
	}

	w.WriteHeader(int(resp.Status))
	if _, err := w.Write(data); err != nil {
		m.h.Error(core.Erro, err) // 此时状态码已经输出
	}
}

// 需要保证 ct 的值不能为空
func findRequestByContentType(requests []*spec.Request, ct string) *spec.Request {
	var none *spec.Request
	for _, req := range requests {
		if req.Mimetype == ct {
			return req
		} else if none == nil && req.Mimetype == "" {
			none = req
		}
	}

	if none != nil {
		return none
	}
	return nil
}

// accepts 必须是已经按权重进行排序的。
func findResponseByAccept(mimetypes []string, requests []*spec.Request, accepts []*qheader.Header) (*spec.Request, string) {
	if len(requests) == 0 {
		return nil, ""
	}

	var none *spec.Request // 表示 requests 中 mimetype 值为空的第一个子项

	// 从 requests 中查找是否有符合 accepts 的内容
	for _, req := range requests {
		if none == nil && req.Mimetype == "" {
			none = req
		}
		if req.Mimetype != "" && matchContentType(req.Mimetype, accepts) {
			return req, req.Mimetype
		}
	}

	// 如果存在 none，则从 doc.mimetypes 中查找是否有与 none.Mimetype 相匹配的
	if none != nil {
		for _, mt := range mimetypes {
			if mt != "" && matchContentType(mt, accepts) {
				return none, mt
			}
		}
	}

	return nil, ""
}

// 查看 ct 是否有与 accepts 相匹配的项，必须保证 ct 的值不为空。
func matchContentType(ct string, accepts []*qheader.Header) bool {
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
func (m *Mock) handleError(w http.ResponseWriter, r *http.Request, field string, err error) {
	// 这并不是一个真实存在的 URI
	file := core.URI(r.Method + ": " + r.URL.Path)

	if serr, ok := err.(*core.SyntaxError); ok {
		if serr.Field == "" {
			serr.Field = field
		} else {
			serr.Field = field + serr.Field
		}

		serr.Location.URI = file
	} else {
		err = core.WithError(core.Location{URI: file}, field, err)
	}

	m.h.Error(core.Erro, err)
	w.WriteHeader(http.StatusBadRequest)
}

func validQueryArrayParam(queries []*spec.Param, r *http.Request) error {
	for _, query := range queries {
		field := "queries[" + query.Name + "]."

		valid := func(p *spec.Param, v string) error {
			err := validSimpleParam(p, v)
			if serr, ok := err.(*core.SyntaxError); ok {
				serr.Field = field + serr.Field
			}
			return err
		}

		if !query.Array {
			if err := valid(query, r.FormValue(query.Name)); err != nil {
				return err
			}
		} else if !query.ArrayStyle { // 默认的 form 格式
			if err := r.ParseForm(); err != nil {
				return err
			}
			for _, v := range r.Form[query.Name] {
				if err := valid(query, v); err != nil {
					return err
				}
			}
		} else {
			values := strings.Split(r.FormValue(query.Name), ",")
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
func validSimpleParam(p *spec.Param, val string) error {
	if p == nil {
		return nil
	}

	if val == "" && p.Type != spec.String { // 字符串的默认值可以为 “”
		if p.Optional || p.Default != "" {
			return nil
		}
		return core.NewLocaleError(core.Location{}, "", locale.ErrRequired)
	}

	switch p.Type {
	case spec.Bool:
		if _, err := strconv.ParseBool(val); err != nil {
			return core.NewLocaleError(core.Location{}, "", locale.ErrInvalidFormat)
		}
	case spec.Number:
		if !is.Number(val) {
			return core.NewLocaleError(core.Location{}, "", locale.ErrInvalidFormat)
		}
	case spec.String:
	case spec.Object:
	case spec.None:
		if val != "" {
			return core.NewLocaleError(core.Location{}, "", locale.ErrInvalidValue)
		}
	}

	if p.IsEnum() {
		found := false
		for _, e := range p.Enums {
			if e.Value == val {
				found = true
				break
			}
		}

		if !found {
			return core.NewLocaleError(core.Location{}, "", locale.ErrInvalidValue)
		}
	}

	return nil
}

func buildResponse(p *spec.Request, r *http.Request) ([]byte, error) {
	if p == nil {
		return nil, nil
	}

	for _, header := range p.Headers {
		if err := validSimpleParam(header, r.Header.Get(header.Name)); err != nil {
			return nil, err
		}
	}

	contentType := r.Header.Get("Accept")
	switch contentType {
	case "application/json":
		return buildJSON(p)
	case "application/xml", "text/xml":
		return buildXML(p)
	default:
		return nil, core.NewLocaleError(core.Location{}, "headers[accept]", locale.ErrInvalidValue)
	}
}
