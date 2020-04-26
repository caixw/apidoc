// SPDX-License-Identifier: MIT

package ast

import (
	"strconv"
	"testing"

	"github.com/issue9/assert"
	"github.com/issue9/version"
)

func TestVersion(t *testing.T) {
	a := assert.New(t)
	a.True(version.SemVerValid(Version))

	v := &version.SemVersion{}
	a.NotError(version.Parse(v, Version))
	major, err := strconv.Atoi(MajorVersion[1:])
	a.NotError(err)
	a.Equal(major, v.Major)
}
