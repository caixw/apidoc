// SPDX-License-Identifier: MIT

package mock

import (
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/issue9/is"
	"github.com/issue9/qheader"

	"github.com/caixw/apidoc/v5/doc"
	"github.com/caixw/apidoc/v5/internal/locale"
	"github.com/caixw/apidoc/v5/internal/vars"
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

		ct := r.Header.Get("Content-Type")
		if ct == "" || ct == "*/*" {
			m.handleError(api, w, "headers[content-type]", locale.Errorf(locale.ErrRequired))
			return
		}
		req, ct := m.findRequestByContentType(api.Requests, []*qheader.Header{{Value: ct}})
		if req == nil {
			m.handleError(api, w, "headers[content-type]", locale.Errorf(locale.ErrInvalidValue))
			return
		}

		if err := validRequest(req, r, ct); err != nil {
			m.handleError(api, w, "request.body.", err)
			return
		}

		accepts, err := qheader.Accept(r)
		if err != nil {
			m.handleError(api, w, "request.headers[Accept]", err)
			return
		}
		resp, accept := m.findRequestByContentType(api.Responses, accepts)
		if resp == nil {
			m.handleError(api, w, "headers[Accept]", locale.Errorf(locale.ErrInvalidValue))
			return
		}

		data, err := buildResponse(resp, r)
		if err != nil {
			m.handleError(api, w, "response.body.", err)
			return
		}

		w.Header().Set("Content-Type", accept) // 需要在输出状态码之前
		w.Header().Set("Server", vars.Name)
		w.WriteHeader(int(resp.Status))
		for _, item := range resp.Headers {
			switch item.Type {
			case doc.Bool:
				w.Header().Set(item.Name, strconv.FormatBool(generateBool()))
			case doc.Number:
				w.Header().Set(item.Name, strconv.FormatInt(generateNumber(item), 10))
			case doc.String:
				w.Header().Set(item.Name, generateString(item))
			default:
				m.handleError(api, w, "response.headers", locale.Errorf(locale.ErrInvalidFormat))
				return
			}
		}
		if _, err := w.Write(data); err != nil {
			m.h.Error(message.Erro, err) // 此时状态码已经输出
		}
	})
}

// 查找第一个符合条件的 Request 实例，如果用户定义了多个相同 mimetype 的实例，也只返回第一符合要求的
func (m *Mock) findRequestByContentType(request []*doc.Request, accepts []*qheader.Header) (*doc.Request, string) {
	var none *doc.Request // 未指定任何 mimetype 值的 doc.Request
	var canMatchAny bool

	for _, req := range request {
		for _, accept := range accepts {
			if accept.Value != "*/*" && accept.Value == req.Mimetype {
				return req, accept.Value
			} else if accept.Value == "*" {
				canMatchAny = true
			}
		}

		if req.Mimetype == "" {
			none = req
		}
	}

	if none != nil {
		for _, mt := range m.doc.Mimetypes {
			for _, accept := range accepts {
				if accept.Value != "*/*" && accept.Value == mt {
					return none, accept.Value
				}
			}
		}
	}

	if canMatchAny {
		mimetype := request[0].Mimetype
		if mimetype == "" {
			mimetype = m.doc.Mimetypes[0]
		}
		return request[0], mimetype
	}

	return nil, ""
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
	w.WriteHeader(http.StatusBadRequest)
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
