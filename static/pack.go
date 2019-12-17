// SPDX-License-Identifier: MIT

package static

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/issue9/utils"
)

const modulePath = "github.com/caixw/apidoc/v5/static"

const header = "// 当前文件由工具自动生成，请勿手动修改！\n\n"

var allowPackExts = []string{
	".xml", ".xsl", ".svg",
	".css", ".js", ".html", ".htm",
}

// FileInfo 被打包文件的信息
type FileInfo struct {
	// 相对于打包根目录的地址，同时也会被作为路由地址
	Name string

	ContentType string
	Content     []byte
}

// Pack 打包 root 目录下的内容到 path 文件
//
// root 需要打包的目录；
// pkgName 输出的包名；
// varName 输出的变量名；
// path 内容保存的文件名；
// t 打包的文件类型，如果为 TypeNone，则只打包 addTo 的内容；
// addTo 追加的打包内容；
//
// NOTE: 隐藏文件不会被打包
func Pack(root, pkgName, varName, path string, t Type, addTo ...*FileInfo) error {
	fis, err := getFileInfos(root, t)
	if err != nil {
		return err
	}
	fis = append(fis, addTo...)

	buf := bytes.NewBufferString(header)

	ws := func(str ...string) {
		for _, s := range str {
			if err == nil {
				_, err = buf.WriteString(s)
			}
		}
	}

	ws("package ", pkgName, "\n\n")

	ws("import \"", modulePath, "\"\n\n")

	ws("var ", varName, "= []*static.FileInfo{")
	for _, info := range fis {
		if err = dump(buf, info); err != nil {
			return err
		}
	}

	// end var pack.FileInfo
	buf.WriteString("}\n")

	return utils.DumpGoFile(path, buf.String())
}

func getFileInfos(root string, t Type) ([]*FileInfo, error) {
	if t == TypeNone {
		return nil, nil
	}

	var paths []string

	walk := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 过滤各类未知的隐藏文件
		if info.IsDir() || !isAllowPackFile(filepath.Ext(info.Name())) {
			return nil
		}

		relPath, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}

		if t != TypeStylesheet || isStylesheetFile(relPath) {
			paths = append(paths, relPath)
		}

		return nil
	}

	if err := filepath.Walk(root, walk); err != nil {
		return nil, err
	}

	fis := make([]*FileInfo, 0, len(paths))
	for _, path := range paths {
		content, err := ioutil.ReadFile(filepath.Join(root, path))
		if err != nil {
			return nil, err
		}
		fis = append(fis, &FileInfo{
			Name:        filepath.ToSlash(path),
			Content:     content,
			ContentType: http.DetectContentType(content),
		})
	}

	return fis, nil
}

func dump(buf *bytes.Buffer, file *FileInfo) (err error) {
	ws := func(str ...string) {
		for _, s := range str {
			if err == nil {
				_, err = buf.WriteString(s)
			}
		}
	}

	ws("{\n")

	ws("Name:\"", file.Name, "\",\n")
	ws("ContentType:\"", file.ContentType, "\",\n")
	ws("Content:[]byte(`", string(file.Content), "`),\n")

	ws("},\n")
	return err
}

func isAllowPackFile(ext string) bool {
	for _, e := range allowPackExts {
		if e == ext {
			return true
		}
	}
	return false
}
