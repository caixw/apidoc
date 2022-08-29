// SPDX-License-Identifier: MIT

//go:build ignore
// +build ignore

package main

import (
	"github.com/caixw/apidoc/v7/internal/docs"
	"github.com/caixw/apidoc/v7/internal/docs/site"
)

func main() {
	if err := site.Write(docs.Dir()); err != nil {
		panic(err)
	}
}
