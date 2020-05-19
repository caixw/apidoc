// SPDX-License-Identifier: MIT

// +build ignore

package main

import (
	"github.com/caixw/apidoc/v7/internal/docs"
	"github.com/caixw/apidoc/v7/internal/docs/makeutil"
	"github.com/caixw/apidoc/v7/internal/docs/site"
)

func main() {
	makeutil.PanicError(site.Write(docs.Dir()))
}
