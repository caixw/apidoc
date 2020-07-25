// SPDX-License-Identifier: MIT

package protocol

import (
	"testing"

	"github.com/issue9/assert"

	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/internal/xmlenc"
)

func TestBuildFoldingRange(t *testing.T) {
	a := assert.New(t)

	base := xmlenc.Base{}
	a.Equal(BuildFoldingRange(base, false), FoldingRange{Kind: FoldingRangeKindComment})

	base = xmlenc.Base{Range: core.Range{
		Start: core.Position{Line: 1, Character: 11},
		End:   core.Position{Line: 2, Character: 11},
	}}
	a.Equal(BuildFoldingRange(base, false), FoldingRange{
		StartLine: 1,
		Kind:      FoldingRangeKindComment,
		EndLine:   2,
	})
	a.Equal(BuildFoldingRange(base, true), FoldingRange{
		StartLine:      1,
		StartCharacter: &base.Start.Character,
		EndLine:        2,
		EndCharacter:   &base.End.Character,
		Kind:           FoldingRangeKindComment,
	})
}
