// SPDX-License-Identifier: MIT

package apidoc

import (
	"io/ioutil"
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

func TestValid(t *testing.T) {
	a := assert.New(t)

	data, err := ioutil.ReadFile("./docs/example/index.xml")
	a.NotError(err).NotNil(data)
	a.NotError(Valid(data))
}
