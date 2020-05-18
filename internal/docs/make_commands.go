// SPDX-License-Identifier: MIT

// +build ignore

package main

import (
	"bytes"
	"io"

	"github.com/caixw/apidoc/v7/internal/cmd"
	"github.com/caixw/apidoc/v7/internal/docs"
	"github.com/caixw/apidoc/v7/internal/docs/makeutil"
	"github.com/caixw/apidoc/v7/internal/locale"
)

type commands struct {
	XMLName struct{} `xml:"commands"`

	Commands []*command `xml:"command"`
}

type command struct {
	Name  string `xml:"name,attr"`
	Usage string `xml:",innerxml"`
}

func main() {
	for _, tag := range locale.Tags() {
		path := "commands." + tag.String() + ".xml"
		out := new(bytes.Buffer)
		opt := cmd.Init(out, tag)
		names := opt.Commands()

		cmds := &commands{
			Commands: make([]*command, 0, len(names)),
		}
		for _, name := range names {
			out.Reset()
			if err := opt.Exec([]string{"help", name}); err != nil {
				panic(err)
			}

			usage, err := out.ReadString('\n')
			if err != nil && err != io.EOF {
				panic(err)
			}

			cmds.Commands = append(cmds.Commands, &command{
				Name:  name,
				Usage: usage,
			})
		}

		makeutil.PanicError(makeutil.WriteXML(docs.Dir().Append(path), cmds, "\t"))
	}
}
