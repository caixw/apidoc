// SPDX-License-Identifier: MIT

package protocol

import (
	"testing"

	"github.com/issue9/assert"
)

func TestSymbolKind(t *testing.T) {
	a := assert.New(t)
	a.Equal(26, SymbolKindTypeParameter)
}

func TestCompletionItemKind(t *testing.T) {
	a := assert.New(t)
	a.Equal(25, CompletionItemKindTypeParameter)
}
