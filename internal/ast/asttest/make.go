// SPDX-License-Identifier: MIT

//go:build ignore
// +build ignore

package main

import (
	"io/ioutil"
	"os"

	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/internal/ast/asttest"
	"github.com/caixw/apidoc/v7/internal/xmlenc"
)

func main() {
	data, err := xmlenc.Encode("\t", asttest.Get(), core.XMLNamespace, "aa")
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(asttest.Filename, data, os.ModePerm)
	if err != nil {
		panic(err)
	}
}
