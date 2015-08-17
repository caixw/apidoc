// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package scanner

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"
	"sync"
	"unicode"
	"unicode/utf8"

	"github.com/caixw/apidoc/core"
	"github.com/caixw/apidoc/log"
)

const eof = -1

var (
	docs   = []*core.Doc{}
	docsMu sync.Mutex
)

// 扫描scanner中的代码，提取最近的下一个代码块和其开始的行号。
// scanFunc必须是一个无状态的
type scanFunc func(*scanner) ([]rune, int, error)

type scanner struct {
	data  []byte
	pos   int
	width int
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
func scanFile(f scanFunc, path string) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Error(err)
		return
	}

	s := &scanner{
		data: data,
	}

	fileWaiter := sync.WaitGroup{}
	for !s.atEOF() {
		block, lineNum, err := f(s)
		if err != nil {
			log.Error(err)
			return
		}

		fileWaiter.Add(1)
		go func(block []rune, lineNum int, path string) {
			defer fileWaiter.Done()
			doc, err := core.Scan(block, lineNum, path)
			if err != nil {
				log.Error(err)
				return
			}
			if doc == nil {
				return
			}
			docsMu.Lock()
			docs = append(docs, doc)
			docsMu.Unlock()
		}(block, lineNum, path)
	} // end for
	fileWaiter.Wait()
}

// 分析dir目录下的文件。并将其转换为core.Docs类型返回。
// recursive 是否递归查询dir子目录下的内容；
// langName 语言名称，不区分大小写，所有代码都将按该语言的语法进行分析；
// exts 可分析的文件扩展名，扩展名必须以点号开头，若不指定，则使用默认的扩展名。
func Scan(dir string, recursive bool, langName string, exts []string) ([]*core.Doc, error) {
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

	paths, err := recursivePath(dir, recursive, exts...)
	if err != nil {
		return nil, err
	}

	waiter := sync.WaitGroup{}
	for _, path := range paths {
		waiter.Add(1)
		go func(path string) {
			scanFile(l.scan, path)
			waiter.Done()
		}(path)
	}
	waiter.Wait()

	return docs, nil
}
