// SPDX-License-Identifier: MIT

package mock

import (
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/issue9/is"
	"github.com/issue9/qheader"

	"github.com/caixw/apidoc/v5/doc"
	"github.com/caixw/apidoc/v5/internal/locale"
	"github.com/caixw/apidoc/v5/internal/vars"
	"github.com/caixw/apidoc/v5/message"
)

func (m *Mock) buildAPI(api *doc.API) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.h.Message(message.Succ, locale.RequestAPI, r.Method, r.URL.Path)
		if api.Deprecated != "" {
			m.h.Message(message.Warn, locale.DeprecatedWarn, r.Method, r.URL.Path, api.Deprecated)
		}

		for _, query := range api.Path.Queries {
			if err := validParam(query, r.FormValue(query.Name)); err != nil {
				m.handleError(w, r, "queries["+query.Name+"]", err)
				return
			}
		}

		for _, header := range api.Headers {
			if err := validParam(header, r.Header.Get(header.Name)); err != nil {
				m.handleError(w, r, "headers["+header.Name+"]", err)
				return
			}
		}

		ct := r.Header.Get("Content-Type")
		if ct == "" || ct == "*/*" || strings.HasSuffix(ct, "/*") { // 用户提交的 content-type 必须是明确的值
			m.handleError(w, r, "headers[content-type]", locale.Errorf(locale.ErrInvalidValue))
			return
		}
		req := findRequestByContentType(api.Requests, ct)
		if req == nil {
			m.handleError(w, r, "headers[content-type]", locale.Errorf(locale.ErrInvalidValue))
			return
		}

		if err := validRequest(req, r, ct); err != nil {
			m.handleError(w, r, "request.body.", err)
			return
		}

		accepts, err := qheader.Accept(r)
		if err != nil {
			m.handleError(w, r, "request.headers[Accept]", err)
			return
		}
		resp, accept := findRequestByAccept(m.doc.Mimetypes, api.Responses, accepts)
		if resp == nil {
			m.handleError(w, r, "headers[Accept]", locale.Errorf(locale.ErrInvalidValue))
			return
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
			case doc.Bool:
				w.Header().Set(item.Name, strconv.FormatBool(generateBool()))
			case doc.Number:
				w.Header().Set(item.Name, strconv.FormatInt(generateNumber(item), 10))
			case doc.String:
				w.Header().Set(item.Name, generateString(item))
			default:
				m.handleError(w, r, "response.headers", locale.Errorf(locale.ErrInvalidFormat))
				return
			}
		}

		w.WriteHeader(int(resp.Status))
		if _, err := w.Write(data); err != nil {
			m.h.Error(message.Erro, err) // 此时状态码已经输出
		}
	})
}

func findRequestByContentType(requests []*doc.Request, ct string) *doc.Request {
	var none *doc.Request
	for _, req := range requests {
		if req.Mimetype == ct {
			return req
		}

		if none == nil && req.Mimetype == "" {
			none = req
		}
	}

	if none != nil {
		return none
	}

	return nil
}

// 查看 ct1 是否与 ct2 匹配，ct2 可以是模糊匹配，比如 text/* 或是 */*
func matchContentType(ct1, ct2 string) string {
	switch {
	case ct1 == ct2:
		return ct1
	case ct2 == "*/*":
		return ct1
	case strings.HasSuffix(ct2, "/*"):
		ct2 = ct2[:len(ct2)-1]
		if strings.HasPrefix(ct1, ct2) {
			return ct1
		}
	}
	return ""
}

// accepts 必须是已经按权重进行排序的。
func findRequestByAccept(mimetypes []string, requests []*doc.Request, accepts []*qheader.Header) (*doc.Request, string) {
	if len(requests) == 0 {
		return nil, ""
	}

	var none *doc.Request // 表示 requests 中 mimetype 值为空的第一个子项

	// 从 requests 中查找是否有符合 accepts 的内容，
	// 同时也获取 requests 中 mimetype 为空值项赋予 none，
	// 以及根据 accepts 内容是否可以匹配任意项。
	for _, req := range requests {
		for _, accept := range accepts {
			ct := matchContentType(req.Mimetype, accept.Value)
			if ct != "" {
				return req, ct
			}
		}

		if req.Mimetype == "" {
			none = req
		}
	}

	// 如果存在 none，则从 doc.mimetypes 中查找同时存在于 doc.mimetypes
	// 与 accepts 的 content-type 作为 none 的 content-type 值返回。
	if none != nil {
		for _, mt := range mimetypes {
			for _, accept := range accepts {
				ct := matchContentType(mt, accept.Value)
				if ct != "" {
					return none, ct
				}
			}
		}
	}

	return nil, ""
}

// 处理 serveHTTP 中的错误
func (m *Mock) handleError(w http.ResponseWriter, r *http.Request, field string, err error) {
	file := r.Method + " " + r.URL.Path

	if serr, ok := err.(*message.SyntaxError); ok {
		if serr.Field == "" {
			serr.Field = field
		} else {
			serr.Field = field + serr.Field
		}

		serr.File = file
	} else {
		err = message.WithError(file, field, 0, err)
	}

	m.h.Error(message.Erro, err)
	w.WriteHeader(http.StatusBadRequest)
}

// 验证单个参数
func validParam(p *doc.Param, val string) error {
	if p == nil {
		return nil
	}

	if val == "" && p.Type != doc.String { // 字符串的默认值可以为 “”
		if p.Optional || p.Default != "" {
			return nil
		}

		return message.NewLocaleError("", "", 0, locale.ErrRequired)
	}

	// TODO 如何验证数组的值？

	switch p.Type {
	case doc.Bool:
		if _, err := strconv.ParseBool(val); err != nil {
			return message.NewLocaleError("", "", 0, locale.ErrInvalidFormat)
		}
	case doc.Number:
		if !is.Number(val) {
			return message.NewLocaleError("", "", 0, locale.ErrInvalidFormat)
		}
	case doc.String:
	case doc.Object:
	case doc.None:
		if val != "" {
			return message.NewLocaleError("", "", 0, locale.ErrInvalidValue)
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
			return message.NewLocaleError("", "", 0, locale.ErrInvalidValue)
		}
	}

	return nil
}

func validRequest(p *doc.Request, r *http.Request, contentType string) error {
	if p == nil {
		return nil
	}

	for _, header := range p.Headers {
		if err := validParam(header, r.Header.Get(header.Name)); err != nil {
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

	switch contentType {
	case "application/json":
		return validJSON(p, content)
	case "application/xml", "text/xml":
		return validXML(p, content)
	default:
		return message.NewLocaleError("", "headers[content-type]", 0, locale.ErrInvalidValue)
	}
}

func buildResponse(p *doc.Request, r *http.Request) ([]byte, error) {
	if p == nil {
		return nil, nil
	}

	for _, header := range p.Headers {
		if err := validParam(header, r.Header.Get(header.Name)); err != nil {
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
		return nil, message.NewLocaleError("", "headers[accept]", 0, locale.ErrInvalidValue)
	}
}
