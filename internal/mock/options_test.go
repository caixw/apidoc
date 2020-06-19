// SPDX-License-Identifier: MIT

package mock

import "github.com/caixw/apidoc/v7/internal/ast"

const indent = "    "

var testOptions = &GenOptions{
	Number: func(p *ast.Param) interface{} { return 1024 },
	String: func(p *ast.Param) string {
		switch p.Type.V() {
		case ast.TypeEmail:
			return "user@example.com"
		case ast.TypeURL:
			return "https://example.com"
		}
		return "1024"
	},
	Bool:      func() bool { return true },
	SliceSize: func() int { return 5 },
	Index:     func(max int) int { return 0 },
}
