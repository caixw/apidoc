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

// 扫描scanner中的代码，提取最近的下一个代码块和其开始的行号。
type scanFunc func(*scanner) ([]rune, int, error)

type scanner struct {
	f     scanFunc
	docs  core.Docs
	data  []byte
	pos   int
	width int
}

func newScanner(f scanFunc) (*scanner, error) {
	return &scanner{
		f:    f,
		docs: core.NewDocs(),
	}, nil
}

// 是否已经在文件末尾。
func (s *scanner) atEOF() bool {
	return s.pos >= len(s.data)
}

// 获取当前的字符，并将指针指向下一个字符。
func (s *scanner) next() rune {
	if s.atEOF() {
		return eof
	}

	r, w := utf8.DecodeRune(s.data[s.pos:])
	s.pos += w
	s.width = w
	return r
}

// 是否匹配指定的字符串，若匹配，则将指定移向该字符串这后，否则不作任何操作。
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

// 撤消s.next()/s.match()的最后一次操作。
func (s *scanner) backup() {
	s.pos -= s.width
	s.width = 0
}

// 跳过之后的所有空白字符。
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
	// 当前行应该是\n数量加1
	return bytes.Count(s.data[:s.pos], []byte("\n")) + 1
}

// 扫描指定的文件到docs
func (s *scanner) scan(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	s.data = data
	s.width = 0
	s.pos = 0

	for !s.atEOF() {
		block, lineNum, err := s.f(s)
		if err != nil {
			return err
		}

		flag := false
		for _, r := range block {
			if r == '@' {
				flag = true
			}
		}
		if !flag {
			continue
		}

		err = s.docs.Scan(block, lineNum, path)
		if err != nil {
			return err
		}
	} // end for

	return nil
}

// 分析dir目录下的文件。并将其转换为core.Docs类型返回。
// recursive 是否递归查询dir子目录下的内容；
// langName 语言名称，不区分大小写，所有代码都将按该语言的语法进行分析；
// exts 可分析的文件扩展名，扩展名必须以点号开头，若不指定，则使用默认的扩展名。
func Scan(dir string, recursive bool, langName string, exts []string) (core.Docs, error) {
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

	return s.docs, nil
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
