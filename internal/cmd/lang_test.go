// SPDX-License-Identifier: MIT

package cmd

import (
	"bufio"
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/issue9/assert"

	"github.com/caixw/apidoc/v5/internal/lang"
)

func TestLanguage(t *testing.T) {
	a := assert.New(t)
	w := new(bytes.Buffer)

	lines := func(w *bytes.Buffer) []string {
		b := bufio.NewReader(w)
		lines := make([]string, 0, 100)
		for line, err := b.ReadString('\n'); err != io.EOF; line, err = b.ReadString('\n') {
			lines = append(lines, line)
		}

		return lines
	}

	a.NotError(language(w))
	ls := lines(w)
	a.Equal(len(ls), len(lang.Langs())+1)
	for _, l := range ls {
		cnt := strings.Count(l, strings.Repeat(" ", 3))
		a.True(cnt >= 2)
	}
}
