// SPDX-License-Identifier: MIT

package cmd

import (
	"bytes"
	"testing"

	"github.com/issue9/assert/v2"

	"github.com/caixw/apidoc/v7/internal/docs"
)

func TestCmdCheckSyntax(t *testing.T) {
	a := assert.New(t, false)

	buf := new(bytes.Buffer)
	cmd := Init(buf)
	erro, _, succ, _ := resetPrinters()
	err := cmd.Exec([]string{"syntax", "-d", docs.Dir().Append("example").String()})
	a.NotError(err)
	a.Empty(buf.String()).
		Empty(erro.String()).
		NotEmpty(succ.String())
}
