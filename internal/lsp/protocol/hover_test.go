// SPDX-License-Identifier: MIT

package protocol

import (
	"encoding/json"
	"testing"

	"github.com/issue9/assert/v3"

	"github.com/caixw/apidoc/v7/core"
)

func TestHover_MarshalJSON(t *testing.T) {
	a := assert.New(t, false)

	h := &Hover{}
	data, err := json.Marshal(h)
	a.NotError(err).Equal(string(data), `null`)

	h.Range = core.Range{
		Start: core.Position{Line: 0, Character: 10},
		End:   core.Position{Line: 1, Character: 10},
	}
	data, err = json.Marshal(h)
	a.NotError(err).Equal(string(data), `null`)

	h.Contents = MarkupContent{
		Kind: MarkupKindPlainText,
	}
	data, err = json.MarshalIndent(h, "", "\t")
	a.NotError(err).Equal(string(data), `{
	"contents": {
		"kind": "plaintext",
		"value": ""
	},
	"range": {
		"start": {
			"line": 0,
			"character": 10
		},
		"end": {
			"line": 1,
			"character": 10
		}
	}
}`)
}
