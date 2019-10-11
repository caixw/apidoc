// SPDX-License-Identifier: MIT

package vars

import (
	"testing"

	"github.com/issue9/assert"
	v "github.com/issue9/version"
)

// 对一些堂量的基本检测。
func TestVersion(t *testing.T) {
	a := assert.New(t)

	a.True(v.SemVerValid(mainVersion))
	a.True(v.SemVerValid(version))
	a.True(v.SemVerValid(Version()))
}

func TestDocVersion(t *testing.T) {
	a := assert.New(t)
	v := DocVersion()
	a.True(len(v) >= 2).
		Equal(v[0], 'v').
		NotContains(v, '.')
}
