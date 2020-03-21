// SPDX-License-Identifier: MIT

package docs

import (
	"github.com/caixw/apidoc/v6/core"
	"github.com/caixw/apidoc/v6/internal/path"
)

var docsDir core.URI

func init() {
	dir, err := core.FileURI(path.CurrPath("../../docs"))
	if err != nil {
		panic(err)
	}

	docsDir = dir
}

// Dir 指向 /docs 的路径
func Dir() core.URI {
	return docsDir
}
