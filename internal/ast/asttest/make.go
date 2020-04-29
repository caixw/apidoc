// SPDX-License-Identifier: MIT

// +build ignore

package main

import (
	"io/ioutil"
	"os"

	"github.com/caixw/apidoc/v7/internal/ast/asttest"
	"github.com/caixw/apidoc/v7/internal/token"
)

func main() {
	data, err := token.Encode("\t", asttest.Get())
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(asttest.Filename, data, os.ModePerm)
	if err != nil {
		panic(err)
	}
}
