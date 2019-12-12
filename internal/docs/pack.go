// SPDX-License-Identifier: MIT

package docs

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

const header = "当前文件由工具自动生成，请勿手动修改！\n\n"

// FileInfo 被打包文件的信息
type FileInfo struct {
	// 相对于打包根目录的地址，同时也会被作为路由地址
	Name string

	ContentType string
	Content     []byte
}

// Pack 打包 root 目录下的内容到 path 文件
//
// pkgName 输出的包名；
// 如果 path 指定的文件不存在，会尝试创建；
// exclude 用于指定不需要打包的文件名，相对于 root 目录；
// addTo 追加的打包内容；
func Pack(pkgName, path string, exclude []string, addTo ...*FileInfo) error {
	fis, err := getFileInfos(RootDir, exclude)
	if err != nil {
		return err
	}
	fis = append(fis, addTo...)

	buf := bytes.NewBufferString(header)
	buf.WriteString("package ")
	buf.WriteString(pkgName)
	buf.WriteByte('\n')

	for _, info := range fis {
		if err = dump(buf, info, path); err != nil {
			return err
		}
	}

	return nil
}

func getFileInfos(root string, exclude []string) ([]*FileInfo, error) {
	paths := []string{}
	walk := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			paths = append(paths, path)
		}

		return nil
	}

	if err := filepath.Walk(root, walk); err != nil {
		return nil, err
	}

	fis := make([]*FileInfo, 0, len(paths))
LOOP:
	for _, path := range paths {
		for _, e := range exclude {
			if filepath.FromSlash(e) == filepath.FromSlash(path) {
				continue LOOP
			}
		}

		content, err := ioutil.ReadFile(filepath.Join(root, path))
		if err != nil {
			return nil, err
		}
		fis = append(fis, &FileInfo{
			Name:        path,
			Content:     content,
			ContentType: http.DetectContentType(content),
		})
	}

	return fis, nil
}

func dump(buf *bytes.Buffer, file *FileInfo, path string) error {
	content, err := ioutil.ReadFile(path)

	ws := func(c string) {
		if err == nil {
			_, err = buf.WriteString(c)
		}
	}

	ws("{\n")

	ws("Name:\"")
	ws(file.Name)
	ws("\",\n")

	ws("ContentType:\"")
	ws(file.ContentType)
	ws("\",\n")

	ws("Content:[]byte(`")
	ws(string(content))
	ws("`),\n")

	ws("},\n")
	return err
}
