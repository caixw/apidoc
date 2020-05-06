// SPDX-License-Identifier: MIT

package ast

import "github.com/caixw/apidoc/v7/internal/token"

var(
	_ token.Sanitizer = &Param{}
	_ token.Sanitizer = &APIDoc{}
)
