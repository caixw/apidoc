// SPDX-License-Identifier: MIT

package protocol

import (
	"testing"

	"github.com/issue9/assert"
)

func TestCompletionItemKind(t *testing.T) {
	a := assert.New(t)
	a.Equal(25, CompletionItemKindTypeParameter)
}

func TestTextDocumentSyncKind(t *testing.T) {
	a := assert.New(t)
	a.Equal(2, TextDocumentSyncKindIncremental)
}
