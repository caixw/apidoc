// SPDX-License-Identifier: MIT

package protocol

import (
	"encoding/json"
	"testing"

	"github.com/issue9/assert/v3"
)

func TestCompletionList_MarshalJSON(t *testing.T) {
	a := assert.New(t, false)

	list := &CompletionList{}
	data, err := json.Marshal(list)
	a.NotError(err).Equal(string(data), "null")

	list.IsIncomplete = true
	list.Items = make([]CompletionItem, 0)
	data, err = json.Marshal(list)
	a.NotError(err).Equal(string(data), "null")

	list.Items = append(list.Items, CompletionItem{})
	data, err = json.Marshal(list)
	a.NotError(err).Equal(string(data), `{"isIncomplete":true,"items":[{"label":""}]}`)
}
