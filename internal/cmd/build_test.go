// SPDX-License-Identifier: MIT

package cmd

import (
	"flag"
	"testing"

	"github.com/caixw/apidoc/v7/core"
	"github.com/issue9/assert"
)

func TestGetPath(t *testing.T) {
	a := assert.New(t)
	fs := flag.NewFlagSet("test", flag.ContinueOnError)

	curr, err := core.FileURI("./")
	a.NotError(err).NotEmpty(curr)

	uri, err := getPath(fs)
	a.NotError(err).Equal(curr, uri)

	uri, err = getPath(nil)
	a.NotError(err).Equal(curr, uri)
}
