// SPDX-License-Identifier: MIT

package vars

import (
	"testing"

	"github.com/issue9/assert"
	"github.com/issue9/is"
)

// 对一些堂量的基本检测。
func TestConst(t *testing.T) {
	a := assert.New(t)

	a.True(len(Name) > 0)
	a.True(is.URL(RepoURL))
	a.True(is.URL(OfficialURL))
}

func TestAllowConfigFilenames(t *testing.T) {
	a := assert.New(t)

	a.True(len(AllowConfigFilenames) > 0)
	for _, name := range AllowConfigFilenames {
		a.True(len(name) > 1)
	}
}
