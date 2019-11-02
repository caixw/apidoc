// SPDX-License-Identifier: MIT

package apidoc

import (
	"testing"

	"github.com/issue9/assert"
	"github.com/issue9/version"

	"github.com/caixw/apidoc/v5/internal/vars"
)

func TestVersion(t *testing.T) {
	a := assert.New(t)

	a.Equal(Version(), vars.Version())
	a.True(version.SemVerValid(Version()))
}
