// SPDX-License-Identifier: MIT

// +build ignore

package main

import (
	"encoding/xml"
	"io/ioutil"
	"os"

	"github.com/caixw/apidoc/v5/doc/doctest"
)

func main() {
	data, err := xml.MarshalIndent(doctest.Get(), "", "\t")
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(doctest.Filename, data, os.ModePerm)
	if err != nil {
		panic(err)
	}
}
