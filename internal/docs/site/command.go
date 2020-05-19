// SPDX-License-Identifier: MIT

package site

import (
	"bytes"
	"io"

	"github.com/caixw/apidoc/v7/internal/cmd"
)

func (d *doc) newCommands() error {
	out := new(bytes.Buffer)
	opt := cmd.Init(out)
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
		d.Commands = append(d.Commands, &command{
			Name:  name,
			Usage: usage,
		})
	}

	return nil
}
