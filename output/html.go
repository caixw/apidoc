// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package output

import (
	"html/template"
	"net/http"
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
	Title          string                // 标题
	Version        string                // 文档版本
	BaseURL        string                // 所有接口地址的前缀
	LicenseName    string                // 文档版权的名称
	LicenseURL     string                // 文档版权的地址
	Content        string                // 索引文件的其它内容
	Groups         map[string][]*doc.API // 按组名形式组织的文档集合
	GroupName      string                // 当前分组名称
	Group          []*doc.API            // 当前组的文档集合
	Date           time.Time             // 生成日期
	AppVersion     string                // apidoc 的版本号
	AppName        string                // 程序名称
	AppRepoURL     string                // 仓库地址
	AppOfficialURL string                // 官网地址
	Elapsed        time.Duration         // 生成文档所用的时间
}

// 将 docs 的内容以 html 格式输出。
func renderHTML(docs *doc.Doc, opt *Options) error {
	t, err := compileHTMLTemplate(opt.Template)
	if err != nil {
		return err
	}

	p := buildHTMLPage(docs, opt)
	return renderHTMLGroups(p, t, opt.Dir)
}

// renderHTMLPlus 是 renderHTML 的调试模式
func renderHTMLPlus(docs *doc.Doc, opt *Options) error {
	p := buildHTMLPage(docs, opt)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// 获取分组名称
		groupName := path.Base(r.URL.Path)
		if path.Ext(groupName) == htmlSuffix {
			groupName = strings.TrimSuffix(groupName, htmlSuffix)
		}

		// 存在该分组，则重新编译模板，并输出内容。
		if group, found := p.Groups[groupName]; found {
			p.GroupName = groupName
			p.Group = group
			handleGroup(w, r, opt, p)
			return
		}

		// 否则当作普通的文件请求
		http.FileServer(http.Dir(opt.Template)).ServeHTTP(w, r)
	})

	return http.ListenAndServe(opt.Port, nil)
}

// 编译模板，并输出 p 的内容。
func handleGroup(w http.ResponseWriter, r *http.Request, opt *Options, p *htmlPage) {
	t, err := compileHTMLTemplate(opt.Template)
	if err != nil {
		if opt.ErrorLog != nil {
			opt.ErrorLog.Println(err)
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	tplName := "group"
	if p.GroupName == app.DefaultGroupName {
		tplName = "index"
	}

	if err = t.ExecuteTemplate(w, tplName, p); err != nil {
		if opt.ErrorLog != nil {
			opt.ErrorLog.Println(err)
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func buildHTMLPage(docs *doc.Doc, opt *Options) *htmlPage {
	p := &htmlPage{
		Title:          docs.Title,
		Version:        docs.Version,
		BaseURL:        docs.BaseURL,
		LicenseName:    docs.LicenseName,
		LicenseURL:     docs.LicenseURL,
		Content:        docs.Content,
		AppVersion:     app.Version,
		AppName:        app.Name,
		AppRepoURL:     app.RepoURL,
		AppOfficialURL: app.OfficialURL,
		Elapsed:        opt.Elapsed,
		Date:           time.Now(),
		Groups:         make(map[string][]*doc.API, 100),
	}

	for _, api := range docs.Apis { // 按分组名称进行分类
		name := strings.ToLower(api.Group)
		if p.Groups[name] == nil {
			p.Groups[name] = []*doc.API{}
		}
		p.Groups[name] = append(p.Groups[name], api)
	}

	return p
}

// 编译模板
//
// tplDir 模板所在的路径，其目录下所有的 .html 文件会被编译，不查找子目录。
// 若 tplDir 为空，则使用系统默认的模板。
func compileHTMLTemplate(tplDir string) (*template.Template, error) {
	t := template.New("html").
		Funcs(template.FuncMap{
			"groupURL": func(name string) string { // 根据分组名称，获取其相应的 URL
				return path.Join(".", name+htmlSuffix)
			},
			"dateFormat": func(t time.Time) string { // 格式化日期
				return t.Format(time.RFC3339)
			},
			"nl2br": func(str string) string { // 将字符串的换行符转成 <br />
				return strings.Replace(str, "\n", "<br />", -1)
			},
			"html": func(str string) interface{} { // 转换成 html
				return template.HTML(str)
			},
			"upper": func(str string) string { // 转大写
				return strings.ToUpper(str)
			},
			"lower": func(str string) string { // 转大写
				return strings.ToLower(str)
			},
		})

	if len(tplDir) > 0 { // 自定义模板
		if _, err := t.ParseGlob(path.Join(tplDir, "*.html")); err != nil {
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
