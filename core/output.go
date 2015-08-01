// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package core

import (
	"html/template"
	"os"

	"github.com/caixw/apidoc/core/static"
)

// 将内容以html格式输出到destDir目录下。
func (tree *Tree) OutputHtml(destDir string) error {
	destDir += string(os.PathSeparator)

	t := template.New("core")
	for _, content := range static.Templates {
		template.Must(t.Parse(content))
	}

	if err := tree.outputIndex(t, destDir); err != nil {
		return err
	}

	if err := tree.outputGroup(t, destDir); err != nil {
		return err
	}

	// 输出static
	return static.Output(destDir)
}

// 输出索引页
func (tree *Tree) outputIndex(t *template.Template, destDir string) error {
	index, err := os.Create(destDir + "index.html")
	if err != nil {
		return err
	}
	defer index.Close()

	err = t.ExecuteTemplate(index, "header", nil)
	if err != nil {
		return err
	}
	err = t.ExecuteTemplate(index, "index", tree)
	if err != nil {
		return err
	}
	return t.ExecuteTemplate(index, "footer", tree.Date)
}

// 按分组输出内容页
func (tree *Tree) outputGroup(t *template.Template, destDir string) error {
	for k, v := range tree.Docs {
		group, err := os.Create(destDir + "group_" + k + ".html")
		if err != nil {
			return err
		}
		defer group.Close()

		err = t.ExecuteTemplate(group, "header", nil)
		if err != nil {
			return err
		}
		for _, d := range v {
			err = t.ExecuteTemplate(group, "group", d)
			if err != nil {
				return err
			}
		}
		err = t.ExecuteTemplate(group, "footer", tree.Date)
		if err != nil {
			return err
		}
	}
	return nil
}
