// SPDX-License-Identifier: MIT

package core

import (
	"testing"

	"github.com/issue9/assert/v3"
	"github.com/issue9/version"
)

// 对一些堂量的基本检测。
func TestVersion(t *testing.T) {
	a := assert.New(t, false)

	a.True(version.SemVerValid(Version()))
	a.True(version.SemVerValid(FullVersion()))
}
