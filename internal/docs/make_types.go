// SPDX-License-Identifier: MIT

// +build ignore

package main

import (
	"bytes"
	"io"

	"golang.org/x/text/language"

	"github.com/caixw/apidoc/v7/internal/ast"
	"github.com/caixw/apidoc/v7/internal/cmd"
	"github.com/caixw/apidoc/v7/internal/docs"
	"github.com/caixw/apidoc/v7/internal/docs/localedoc"
	"github.com/caixw/apidoc/v7/internal/docs/makeutil"
	"github.com/caixw/apidoc/v7/internal/locale"
	"github.com/caixw/apidoc/v7/internal/token"
)

func main() {
	for _, tag := range locale.Tags() {
		doc := &localedoc.LocaleDoc{}

		makeutil.PanicError(makeCommands(doc, tag))
		makeutil.PanicError(token.NewTypes(doc, &ast.APIDoc{}, tag))

		target := docs.Dir().Append("localedoc." + tag.String() + ".xml")
		makeutil.PanicError(makeutil.WriteXML(target, doc, "\t"))
	}
}

func makeCommands(doc *localedoc.LocaleDoc, tag language.Tag) error {
	out := new(bytes.Buffer)
	opt := cmd.Init(out, tag)
	names := opt.Commands()

	for _, name := range names {
		out.Reset()
		if err := opt.Exec([]string{"help", name}); err != nil {
			return err
		}

		usage, err := out.ReadString('\n')
		if err != nil && err != io.EOF {
			return err
		}

		if usage[len(usage)-1] == '\n' { // 去掉换行符
			usage = usage[:len(usage)-1]
		}
		doc.Commands = append(doc.Commands, &localedoc.Command{
			Name:  name,
			Usage: usage,
		})
	}

	return nil
}
