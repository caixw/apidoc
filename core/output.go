// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package core

import (
	"html/template"
	"os"

	"github.com/caixw/apidoc/core/static"
)

func (tree *Tree) OutputHtml(destDir string) error {
	destDir += string(os.PathSeparator)
	path := destDir + "index.html"

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	t := template.New("core")
	for _, content := range static.Templates {
		template.Must(t.Parse(content))
	}

	err = t.ExecuteTemplate(file, "main", tree)
	if err != nil {
		return err
	}

	// 输出static
	return static.Output(destDir)
}
