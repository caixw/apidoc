// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package scanner

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/caixw/apidoc/core"
)

const eof = -1

type scanFunc func(*scanner) ([]byte, error)

type scanner struct {
	f     scanFunc
	tree  *core.Tree
	data  []byte
	pos   int
	width int
}

func newScanner(f scanFunc) (*scanner, error) {
	return &scanner{
		f:    f,
		tree: core.NewTree(),
	}, nil
}

func (s *scanner) atEOF() bool {
	return s.pos >= len(s.data)
}

func (s *scanner) next() rune {
	if s.atEOF() {
		return eof
	}

	r, w := utf8.DecodeRune(s.data[s.pos:])
	s.pos += w
	s.width = w
	return r
}

func (s *scanner) match(str string) bool {
	rs := []rune(str)
	if s.atEOF() {
		return false
	}

	width := 0
	for _, r := range rs {
		rr, w := utf8.DecodeRune(s.data[s.pos:])
		if rr != r {
			s.pos -= width
			return false
		}

		s.pos += w
		width += w
	}

	s.width = width
	return true
}

func (s *scanner) backup() {
	s.pos -= s.width
	s.width = 0
}

func (s *scanner) skipSpace() {
	if s.atEOF() {
		return
	}

	for {
		if !unicode.IsSpace(s.next()) {
			s.backup()
			return
		}
	}
}

// 当前所在的行号
func (s *scanner) lineNumber() int {
	return bytes.Count(s.data[:s.pos], []byte("\n"))
}

// 扫描指定的文件到tree
func (s *scanner) scan(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	s.data = data
	s.width = 0
	s.pos = 0

	for !s.atEOF() {
		block, err := s.f(s)
		if err != nil {
			return err
		}

		err = s.tree.Scan(block, s.lineNumber(), path)
		if err != nil {
			return err
		}
	} // end for

	return nil
}

// 分析dir目录下的文件。并将其转换为*core.Tree类型返回。
// recursive 是否递归查询dir子目录下的内容；
// langName 语言名称，不区分大小写，所有代码都将按该语言的语法进行分析；
// exts 可分析的文件扩展名，扩展名必须以点号开头，若不指定，则使用默认的扩展名。
func Scan(dir string, recursive bool, langName string, exts []string) (*core.Tree, error) {
	if len(langName) == 0 {
		var err error
		if len(exts) == 0 {
			langName, err = detectDirLangType(dir)
		} else {
			langName, err = detectLangType(exts)
		}

		if err != nil {
			return nil, err
		}
	}

	l, found := langs[strings.ToLower(langName)]
	if !found {
		return nil, fmt.Errorf("不支持的语言:%v", langName)
	}
	if len(exts) == 0 {
		exts = l.exts
	}

	fmt.Println("scanner:", langName)
	fmt.Println("exts:", exts)

	s, err := newScanner(l.scan)
	if err != nil {
		return nil, err
	}

	paths, err := recursivePath(dir, recursive, exts...)
	if err != nil {
		return nil, err
	}
	for _, path := range paths {
		err = s.scan(path)
		if err != nil {
			return nil, err
		}
	}

	return s.tree, nil
}

// 从扩展名检测其所属的语言名称。
// 以第一个匹配extsIndex的文件扩展名为准。
func detectLangType(exts []string) (string, error) {
	for _, ext := range exts {
		if lang, found := extsIndex[ext]; found {
			return lang, nil
		}
	}
	return "", fmt.Errorf("无法找到与这些扩展名[%v]相匹配的代码扫描函数", exts)
}

// 检测目录下的文件类型。
// 以第一个匹配extsIndex的文件扩展名为准。
func detectDirLangType(dir string) (string, error) {
	var lang string

	walk := func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if fi.IsDir() || len(lang) > 0 {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		lang, _ = extsIndex[ext]
		return nil
	}

	if err := filepath.Walk(dir, walk); err != nil {
		return "", err
	}

	if len(lang) == 0 {
		return lang, fmt.Errorf("无法检测到[%v]目录下的文件类型", dir)
	}

	return lang, nil
}

// 根据recursive值确定是否递归查找paths每个目录下的子目录。
func recursivePath(dir string, recursive bool, exts ...string) ([]string, error) {
	paths := []string{}
	dir += string(os.PathSeparator)

	extIsEnabled := func(ext string) bool {
		for _, v := range exts {
			if ext == v {
				return true
			}
		}
		return false
	}

	walk := func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if fi.IsDir() && !recursive && path != dir {
			return filepath.SkipDir
		} else if extIsEnabled(filepath.Ext(path)) {
			paths = append(paths, path)
		}
		return nil
	}

	if err := filepath.Walk(dir, walk); err != nil {
		return nil, err
	}

	return paths, nil
}
