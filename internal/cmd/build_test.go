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

	curr := core.FileURI("./")
	a.NotEmpty(curr)

	uri := getPath(fs)
	a.Equal(curr, uri)

	uri = getPath(nil)
	a.Equal(curr, uri)
}
