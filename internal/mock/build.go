// SPDX-License-Identifier: MIT

package mock

import (
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/issue9/is"

	"github.com/caixw/apidoc/v5/doc"
	"github.com/caixw/apidoc/v5/internal/locale"
	"github.com/caixw/apidoc/v5/message"
)

func (m *Mock) buildAPI(api *doc.API) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, query := range api.Path.Queries {
			if err := validParam(query, r.FormValue(query.Name)); err != nil {
				m.handleError(api, w, "queries["+query.Name+"]", err)
				return
			}
		}

		for _, header := range api.Headers {
			if err := validParam(header, r.Header.Get(header.Name)); err != nil {
				m.handleError(api, w, "headers["+header.Name+"]", err)
				return
			}
		}

		req := findRequestByContentType(api.Requests, r.Header.Get("Content-Type"))
		if req == nil {
			m.handleError(api, w, "headers[content-type]", locale.Errorf(locale.ErrInvalidValue))
			return
		}

		if err := validRequest(req, r); err != nil {
			m.handleError(api, w, "request.body", err)
			return
		}

		resp := findRequestByContentType(api.Responses, r.Header.Get("Accept"))
		if resp == nil {
			m.handleError(api, w, "headers[Accept]", locale.Errorf(locale.ErrInvalidValue))
			return
		}

		data, err := buildResponse(resp, r)
		if err != nil {
			m.handleError(api, w, "response.body", err)
			return
		}

		w.WriteHeader(int(resp.Status))
		if _, err := w.Write(data); err != nil {
			m.h.Error(message.Erro, err) // 此时状态码已经输出
		}
	})
}

// 查找第一个符合条件的 Request 实例，如果用户定义了多个相同 mimetype 的实例，也只返回第一符合要求的
func findRequestByContentType(request []*doc.Request, ct string) *doc.Request {
	for _, req := range request {
		if req.Mimetype == ct {
			return req
		}
	}

	return nil
}

// 处理 serveHTTP 中的错误
func (m *Mock) handleError(api *doc.API, w http.ResponseWriter, field string, err error) {
	file := string(api.Method) + " " + api.Path.Path

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
}

// 验证单个参数
func validParam(p *doc.Param, val string) error {
	if p == nil {
		return nil
	}

	if val == "" && p.Type != doc.String { // 字符串的默认值可以为 “”
		if p.Optional {
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

func validRequest(p *doc.Request, r *http.Request) error {
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

	contentType := r.Header.Get("Content-Type")
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
