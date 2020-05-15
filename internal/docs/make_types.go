// SPDX-License-Identifier: MIT

// +build ignore

package main

import (
	"github.com/caixw/apidoc/v7/internal/ast"
	"github.com/caixw/apidoc/v7/internal/docs/makeutil"
	"github.com/caixw/apidoc/v7/internal/locale"
	"github.com/caixw/apidoc/v7/internal/token"
)

func main() {
	for _, tag := range locale.Tags() {
		types, err := token.NewTypes(&ast.APIDoc{}, tag)
		makeutil.PanicError(err)

		target := docs.Dir().Append("types." + tag.String() + ".xml")
		makeutil.PanicError(makeutil.WriteXML(target, types, "\t"))
	}
}
