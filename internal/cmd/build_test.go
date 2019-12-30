// SPDX-License-Identifier: MIT

package cmd

import (
	"flag"
	"testing"

	"github.com/issue9/assert"
)

func TestGetPath(t *testing.T) {
	a := assert.New(t)
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	a.Equal("./", getPath(fs))
	a.Equal("./", getPath(nil))
}
