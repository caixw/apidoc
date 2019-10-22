// SPDX-License-Identifier: MIT

package mock

import (
	"net/http"

	"github.com/caixw/apidoc/v5/doc"
)

func build(api *doc.API) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch api.Method {
		case http.MethodGet:
			buildGet(api)
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

func buildGet(api *doc.API) {
	// TODO
}
