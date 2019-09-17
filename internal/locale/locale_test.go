// SPDX-License-Identifier: MIT

package locale

import (
	"testing"

	"github.com/issue9/assert"
)

func TestDisplayNames(t *testing.T) {
	a := assert.New(t)

	a.Equal(len(DisplayNames()), len(locales))
	for tag := range DisplayNames() {
		_, found := locales[tag]
		a.True(found)
	}
}
