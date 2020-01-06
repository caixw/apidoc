// SPDX-License-Identifier: MIT

// +build ignore

package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/issue9/utils"

	"github.com/caixw/apidoc/v5/internal/static"
)

const (
	header    = "// 当前文件由工具自动生成，请勿手动修改！\n\n"
	pkgName   = "static"
	varName   = "data"
	distPath  = "./data.go"
	sourceDir = "../../docs"
)

// 允许打包的文件后缀名，以及对应的 mime type 值。
// 不采用 mimetype.TypeByExtension，防止出现空值的可能性。
var allowFiles = map[string]string{
	".xml":  "application/xml",
	".xsl":  "text/xsl",
	".svg":  "image/svg+xml",
	".css":  "text/css",
	".js":   "application/javascript",
	".html": "text/html",
	".htm":  "text/html",
}

func main() {
	if err := pack(); err != nil {
		panic(err)
	}
}

// NOTE: 隐藏文件不会被打包
func pack() error {
	fis, err := getFileInfos(sourceDir)
	if err != nil {
		return err
	}

	buf := bytes.NewBufferString(header)

	ws := func(str ...string) {
		for _, s := range str {
			if err == nil {
				_, err = buf.WriteString(s)
			}
		}
	}

	ws("package ", pkgName, "\n\n")

	ws("var ", varName, "= []*FileInfo{")
	for _, info := range fis {
		if err = dump(buf, info); err != nil {
			return err
		}
	}

	// end var pack.FileInfo
	ws("}\n")

	if err != nil {
		return err
	}

	return utils.DumpGoFile(distPath, buf.String())
}

func getFileInfos(root string) ([]*static.FileInfo, error) {
	var paths []string

	walk := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 过滤各类未知的隐藏文件
		if info.IsDir() || allowFiles[filepath.Ext(info.Name())] == "" {
			return nil
		}

		relPath, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}

		paths = append(paths, relPath)

		return nil
	}

	if err := filepath.Walk(root, walk); err != nil {
		return nil, err
	}

	fis := make([]*static.FileInfo, 0, len(paths))
	for _, path := range paths {
		content, err := ioutil.ReadFile(filepath.Join(root, path))
		if err != nil {
			return nil, err
		}
		fis = append(fis, &static.FileInfo{
			Name:        filepath.ToSlash(path),
			Content:     content,
			ContentType: allowFiles[filepath.Ext(path)],
		})
	}

	return fis, nil
}

func dump(buf *bytes.Buffer, file *static.FileInfo) (err error) {
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
