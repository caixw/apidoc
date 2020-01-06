// SPDX-License-Identifier: MIT

package vars

import (
	"path/filepath"
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

func TestDocsDir(t *testing.T) {
	a := assert.New(t)

	p1, err := filepath.Abs("../../docs")
	a.NotError(err).NotEmpty(p1)

	p2, err := filepath.Abs(DocsDir())
	a.NotError(err).NotEmpty(p2)

	a.Equal(p1, p2)
}
