// SPDX-License-Identifier: MIT

package cmd

import (
	"bytes"
	"testing"

	"github.com/caixw/apidoc/v7"
	"github.com/issue9/assert"
)

func TestCmdVersion(t *testing.T) {
	a := assert.New(t)

	buf := new(bytes.Buffer)
	cmd := Init(buf)
	resetPrinters()
	a.NotError(cmd.Exec([]string{"version"}))
	a.Contains(buf.String(), apidoc.LSPVersion()).
		Contains(buf.String(), apidoc.DocVersion()).
		Contains(buf.String(), apidoc.Version(true))
}
