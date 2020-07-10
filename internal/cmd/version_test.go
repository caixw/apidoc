// SPDX-License-Identifier: MIT

package cmd

import (
	"bytes"
	"testing"

	"github.com/issue9/assert"

	"github.com/caixw/apidoc/v7"
	"github.com/caixw/apidoc/v7/internal/openapi"
)

func TestCmdVersion(t *testing.T) {
	a := assert.New(t)

	buf := new(bytes.Buffer)
	cmd := Init(buf)
	resetPrinters()
	a.NotError(cmd.Exec([]string{"version"}))
	a.Contains(buf.String(), apidoc.LSPVersion).
		Contains(buf.String(), apidoc.DocVersion).
		Contains(buf.String(), apidoc.Version(true))

	buf.Reset()
	cmd = Init(buf)
	resetPrinters()
	a.NotError(cmd.Exec([]string{"version", "-kind", "apidoc"}))
	a.Equal(buf.String(), apidoc.Version(true)+"\n")

	buf.Reset()
	cmd = Init(buf)
	resetPrinters()
	a.NotError(cmd.Exec([]string{"version", "-kind", "lsp"}))
	a.Equal(buf.String(), apidoc.LSPVersion+"\n")

	buf.Reset()
	cmd = Init(buf)
	resetPrinters()
	a.NotError(cmd.Exec([]string{"version", "-kind", "doc"}))
	a.Equal(buf.String(), apidoc.DocVersion+"\n")

	buf.Reset()
	cmd = Init(buf)
	resetPrinters()
	a.NotError(cmd.Exec([]string{"version", "-kind", "openapi"}))
	a.Equal(buf.String(), openapi.LatestVersion+"\n")
}
