// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package core

import (
	"html/template"
	"os"
	"time"

	"github.com/caixw/apidoc/core/static"
)

// 用于页首和页脚的附加信息
type info struct {
	Groups    map[string]string // 分组名称与文件的对照表
	CurrGroup string            // 当前所在的分组页，若为空，表示在列表页
	Date      string            // 编译日期
	Version   string            // 程序版本
}

// 将内容以html格式输出到destDir目录下。
func (tree *Tree) OutputHtml(destDir string) error {
	destDir += string(os.PathSeparator)

	t := template.New("core")
	for _, content := range static.Templates {
		template.Must(t.Parse(content))
	}

	i := &info{
		Version: tree.Version,
		Date:    time.Now().Format(time.RFC3339),
		Groups:  make(map[string]string, len(tree.Docs)),
	}
	for k, _ := range tree.Docs {
		i.Groups[k] = "./group_" + k + ".html"
	}

	if err := tree.outputIndex(t, i, destDir); err != nil {
		return err
	}

	if err := tree.outputGroup(t, i, destDir); err != nil {
		return err
	}

	// 输出static
	return static.Output(destDir)
}

// 输出索引页
func (tree *Tree) outputIndex(t *template.Template, i *info, destDir string) error {
	index, err := os.Create(destDir + "index.html")
	if err != nil {
		return err
	}
	defer index.Close()

	err = t.ExecuteTemplate(index, "header", i)
	if err != nil {
		return err
	}

	err = t.ExecuteTemplate(index, "index", i)
	if err != nil {
		return err
	}
	return t.ExecuteTemplate(index, "footer", i)
}

// 按分组输出内容页
func (tree *Tree) outputGroup(t *template.Template, i *info, destDir string) error {
	for k, v := range tree.Docs {
		group, err := os.Create(destDir + "group_" + k + ".html")
		if err != nil {
			return err
		}
		defer group.Close()

		i.CurrGroup = k
		err = t.ExecuteTemplate(group, "header", i)
		if err != nil {
			return err
		}
		for _, d := range v {
			err = t.ExecuteTemplate(group, "group", d)
			if err != nil {
				return err
			}
		}
		err = t.ExecuteTemplate(group, "footer", i)
		if err != nil {
			return err
		}
	}
	return nil
}
