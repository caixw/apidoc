// SPDX-License-Identifier: MIT

package mock

import (
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/caixw/apidoc/v5/doc"
	"github.com/caixw/apidoc/v5/internal/locale"
	"github.com/caixw/apidoc/v5/message"
)

func (m *Mock) build(api *doc.API) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch api.Method {
		case http.MethodGet:
			m.buildGet(api, w, r)
		case http.MethodOptions:
			// TODO
		case http.MethodHead:
			// TODO
		case http.MethodPatch:
			// TODO
		case http.MethodPut:
			// TODO
		case http.MethodPost:
			// TODO
		case http.MethodDelete:
			// TODO
		default:
			// TODO
		}
	})
}

// 处理 serveHTTP 中的错误
func (m *Mock) handleError(w http.ResponseWriter, err error) {
	m.h.Error(message.Erro, err)
	if status, ok := err.(*Error); ok {
		w.WriteHeader(status.Status)
	}
}

func (m *Mock) buildGet(api *doc.API, w http.ResponseWriter, r *http.Request) {
	for _, header := range api.Requests[0].Headers {
		if err := validParam(header, r.Header.Get(header.Name)); err != nil {
			m.handleError(w, err)
			return
		}
	}

	for _, query := range api.Path.Queries {
		if err := validParam(query, r.FormValue(query.Name)); err != nil {
			m.handleError(w, err)
			return
		}
	}

	if err := validRequest(api.Requests[0], r); err != nil {
		m.handleError(w, err)
		return
	}

	data, err := buildResponse(api.Responses[0])
	if err != nil {
		m.handleError(w, err)
		return
	}

	w.WriteHeader(int(api.Responses[0].Status))
	if _, err := w.Write(data); err != nil {
		m.h.Error(message.Erro, err) // 此时状态码已经输出
	}
}

// 验证单个参数
func validParam(p *doc.Param, val string) error {
	if val == "" && p.Type != doc.String { // 字符串的默认值可以为 “”
		if p.Optional {
			return nil
		}

		return newError(http.StatusBadRequest, locale.ErrRequired)
	}

	switch p.Type {
	case doc.Bool:
		if _, err := strconv.ParseBool(val); err != nil {
			return err
		}
	case doc.Number:
		if _, err := strconv.ParseInt(val, 10, 32); err != nil {
			return err
		}
	case doc.String:
	case doc.Object:
	case doc.None:
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
			return newError(http.StatusBadRequest, locale.ErrInvalidValue)
		}
	}

	return nil
}

func validRequest(p *doc.Request, r *http.Request) error {
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
	case "application/xml":
	default:
		// TODO
	}
}

func buildResponse(p *doc.Request) ([]byte, error) {
	// TODO
}
