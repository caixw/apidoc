// SPDX-License-Identifier: MIT

package main

import (
	"net/http"

	"github.com/caixw/apidoc/v5"
)

func main() {
	http.Handle("/apidoc/", apidoc.Handle("../../docs/apidoc.xml", "", nil))
	err := http.ListenAndServe(":8080", nil)
	panic(err)
}
