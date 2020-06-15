// SPDX-License-Identifier: MIT

// +build ignore

package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/issue9/pack"

	"github.com/caixw/apidoc/v7/internal/ast"
	"github.com/caixw/apidoc/v7/internal/docs"
)

const (
	pkgName  = "docs"
	varName  = "data"
	distPath = "./static.go"
)

// 允许打包的文件后缀名，以及对应的 mime type 值。
// 不采用 mimetype.TypeByExtension，防止出现空值的可能性。
var allowExts = map[string]string{
	".xml":  "application/xml; charset=utf-8",
	".xsl":  "text/xsl; charset=utf-8",
	".svg":  "image/svg+xml; charset=utf-8",
	".css":  "text/css; charset=utf-8",
	".js":   "application/javascript; charset=utf-8",
	".html": "text/html; charset=utf-8",
	".htm":  "text/html; charset=utf-8",
}

// 允许打包的子目录，仅支持一级目录
var allowDirs = []string{
	ast.MajorVersion + "/",
	"example/",
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
	fis := make([]*docs.FileInfo, 0, 10)

	walk := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		ext := strings.ToLower(filepath.Ext(path))
		if info.IsDir() || allowExts[ext] == "" { // 同时也会过滤各类未知的隐藏文件
			return nil
		}

		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)

		if !isAllowDir(rel) {
			return nil
		}

		content, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		fis = append(fis, &docs.FileInfo{
			Name:        rel,
			ContentType: allowExts[ext],
			Content:     content,
		})

		return nil
	}

	if err := filepath.Walk(root, walk); err != nil {
		return nil, err
	}

	return fis, nil
}

func isAllowDir(rel string) bool {
	if strings.IndexByte(rel, '/') == -1 { // 根目录下的内容始终允许
		return true
	}

	for _, prefix := range allowDirs {
		if strings.HasPrefix(rel, prefix) {
			return true
		}
	}
	return false
}
