// SPDX-License-Identifier: MIT

// +build ignore

package main

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/issue9/pack"

	"github.com/caixw/apidoc/v7/internal/docs"
)

const (
	pkgName  = "docs"
	varName  = "data"
	distPath = "./static.go"
)

// 允许打包的文件后缀名，以及对应的 mime type 值。
// 不采用 mimetype.TypeByExtension，防止出现空值的可能性。
var allowFiles = map[string]string{
	".xml":  "application/xml; charset=utf-8",
	".xsl":  "text/xsl; charset=utf-8",
	".svg":  "image/svg+xml; charset=utf-8",
	".css":  "text/css; charset=utf-8",
	".js":   "application/javascript; charset=utf-8",
	".html": "text/html; charset=utf-8",
	".htm":  "text/html; charset=utf-8",
}

func panicError(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	dir, err := docs.Dir().File()
	panicError(err)

	fis, err := getFileInfos(dir)
	panicError(err)

	panicError(pack.File(fis, pkgName, varName, docs.FileHeader, "", distPath))
}

func getFileInfos(root string) ([]*docs.FileInfo, error) {
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

	fis := make([]*docs.FileInfo, 0, len(paths))
	for _, path := range paths {
		content, err := ioutil.ReadFile(filepath.Join(root, path))
		if err != nil {
			return nil, err
		}
		fis = append(fis, &docs.FileInfo{
			Name:        filepath.ToSlash(path),
			ContentType: allowFiles[filepath.Ext(path)],
			Content:     content,
		})
	}

	return fis, nil
}
