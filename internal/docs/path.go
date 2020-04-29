// SPDX-License-Identifier: MIT

package docs

import (
	"github.com/issue9/utils"

	"github.com/caixw/apidoc/v7/core"
)

var docsDir core.URI

func init() {
	dir, err := core.FileURI(utils.CurrentPath("../../docs"))
	if err != nil {
		panic(err)
	}

	docsDir = dir
}

// Dir 指向 /docs 的路径
func Dir() core.URI {
	return docsDir
}
