// SPDX-License-Identifier: MIT

// +build ignore

package main

import (
	"encoding/xml"
	"io/ioutil"
	"os"

	"github.com/caixw/apidoc/v7/spec/spectest"
)

func main() {
	data, err := xml.MarshalIndent(spectest.Get(), "", "\t")
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(spectest.Filename, data, os.ModePerm)
	if err != nil {
		panic(err)
	}
}
