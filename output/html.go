// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package output

import (
	"html/template"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/caixw/apidoc/app"
	"github.com/caixw/apidoc/doc"
	"github.com/caixw/apidoc/output/static"
)

// 输出的 html 文件后缀名
const htmlSuffix = ".html"

// 用于页首和页脚的附加信息
type htmlPage struct {
	Content        template.HTML         // 索引文件的其它内容
	Groups         map[string][]*doc.API // 按组名形式组织的文档集合
	GroupName      string                // 当前分组名称
	Group          []*doc.API            // 当前组的文档集合
	Date           string                // 生成日期
	Version        string                // 文档版本
	AppVersion     string                // apidoc 的版本号
	AppName        string                // 程序名称
	AppRepoURL     string                // 仓库地址
	AppOfficialURL string                // 官网地址
	Title          string                // 标题
	Elapsed        time.Duration         // 生成文档所用的时间
}

// 将 docs 的内容以 html 格式输出。
func renderHTML(docs *doc.Doc, opt *Options) error {
	t, err := compileHTMLTemplate(opt)
	if err != nil {
		return err
	}

	p := &htmlPage{
		Content:        template.HTML(strings.Replace(docs.Content, "\n", "<br />", -1)),
		Title:          docs.Title,
		Version:        docs.Version,
		AppVersion:     app.Version,
		AppName:        app.Name,
		AppRepoURL:     app.RepoURL,
		AppOfficialURL: app.OfficialURL,
		Elapsed:        opt.Elapsed,
		Date:           time.Now().Format(time.RFC3339), // TODO 可以自定义时间格式？
		Groups:         make(map[string][]*doc.API, 100),
	}

	for _, api := range docs.Apis { // 按分组名称进行分类
		name := strings.ToLower(api.Group)
		if p.Groups[name] == nil {
			p.Groups[name] = []*doc.API{}
		}
		p.Groups[name] = append(p.Groups[name], api)
	}

	return renderHTMLGroups(p, t, opt.Dir)
}

// 编译模板
func compileHTMLTemplate(opt *Options) (*template.Template, error) {
	t := template.New("html").
		Funcs(template.FuncMap{
			"groupURL": func(name string) string {
				return path.Join(".", name+htmlSuffix)
			},
		})

	if len(opt.Template) > 0 { // 自定义模板
		if _, err := t.ParseGlob(path.Join(opt.Template, "*.html")); err != nil {
			return nil, err
		}
	} else { // 系统模板
		for _, content := range static.Templates {
			template.Must(t.Parse(content))
		}
	}

	return t, nil
}

// 根据模板 t 将页面输出到 destDir 目录
func renderHTMLGroups(p *htmlPage, t *template.Template, destDir string) error {
	for name, group := range p.Groups {
		path := filepath.Join(destDir, name+htmlSuffix)
		file, err := os.Create(path)
		if err != nil {
			return err
		}
		defer file.Close()

		p.GroupName = name
		p.Group = group

		tplName := "group"
		if name == app.DefaultGroupName {
			tplName = "index"
		}
		if err = t.ExecuteTemplate(file, tplName, p); err != nil {
			return err
		}
	}

	// 输出static
	return static.Output(destDir)
}
