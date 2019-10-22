// SPDX-License-Identifier: MIT

package mock

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/caixw/apidoc/v5/doc"
)

func build(api *doc.API) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch api.Method {
		case http.MethodGet:
			buildGet(api, w, r)
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

func buildGet(api *doc.API, w http.ResponseWriter, r *http.Request) {
	for _, header := range api.Requests[0].Headers {
		if err := validSimpleParam(header, r.Header.Get(header.Name)); err != nil {
			// TODO
		}
	}

	for _, query := range api.Path.Queries {
		if err := validSimpleParam(query, r.FormValue(query.Name)); err != nil {
			// TODO
		}
	}

	if err := validRequest(api.Requests[0], r); err != nil {
		// TODO
	}

	data, err := buildResponse(api.Responses[0])
	if err != nil {
		// TODO
	}

	w.WriteHeader(int(api.Responses[0].Status))
	if _, err := w.Write(data); err != nil {
		// TODO
	}
}

func validSimpleParam(p *doc.SimpleParam, val string) error {
	if val == "" {
		if p.Optional {
			return nil
		}

		// TODO 如何解决参数错误返回的状态类型？
		return errors.New("400")
	}

	var err error
	switch p.Type {
	case doc.Bool:
		_, err = strconv.ParseBool(val)
	case doc.String:
	case doc.Number:
		_, err = strconv.ParseInt(val, 10, 32)
	case doc.Object:
		err = errors.New("无效的值类型")
	case doc.None:
		err = errors.New("无效的值类型")
	}
	if err != nil {
		return err
	}

	if p.IsEnum() {
		// TODO
	}

	return nil
}

func validRequest(p *doc.Request, r *http.Request) error {
	// TODO
}

func buildResponse(p *doc.Request) ([]byte, error) {
	//
}
