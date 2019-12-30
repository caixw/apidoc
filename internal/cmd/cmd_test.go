// SPDX-License-Identifier: MIT

package cmd

import (
	"bytes"
	"testing"

	"github.com/issue9/assert"
	"github.com/issue9/term/colors"

	"github.com/caixw/apidoc/v5/internal/locale"
)

func TestPrintLocale(t *testing.T) {
	a := assert.New(t)

	buf := new(bytes.Buffer)
	printLocale(buf, colors.Default, locale.ErrRequired)
	a.Contains(buf.String(), locale.ErrRequired)
}
