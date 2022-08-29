// SPDX-License-Identifier: MIT

package protocol

import (
	"testing"

	"github.com/issue9/assert/v3"

	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/internal/xmlenc"
)

func TestBuildFoldingRange(t *testing.T) {
	a := assert.New(t, false)

	base := xmlenc.Base{}
	a.Equal(BuildFoldingRange(base, false), FoldingRange{Kind: FoldingRangeKindComment})

	base = xmlenc.Base{
		Location: core.Location{
			Range: core.Range{
				Start: core.Position{Line: 1, Character: 11},
				End:   core.Position{Line: 2, Character: 11},
			},
		},
	}
	a.Equal(BuildFoldingRange(base, false), FoldingRange{
		StartLine: 1,
		Kind:      FoldingRangeKindComment,
		EndLine:   2,
	})
	a.Equal(BuildFoldingRange(base, true), FoldingRange{
		StartLine:      1,
		StartCharacter: &base.Location.Range.Start.Character,
		EndLine:        2,
		EndCharacter:   &base.Location.Range.End.Character,
		Kind:           FoldingRangeKindComment,
	})
}
