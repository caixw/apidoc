// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package core

import (
	"bytes"
	"io/ioutil"
	"sync"

	"github.com/caixw/apidoc/log"
)

func scanFile(docs *docs, f ScanFunc, path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	var block []rune
	pos := 0
	ln := 0
	wg := sync.WaitGroup{}
	for {
		block, pos = f(data)
		if pos >= len(data) || pos < 0 {
			break
		}
		ln += bytes.Count(data[:pos], []byte("\n"))
		wg.Add(1)
		go func(block []rune, lineNum int, path string) {
			defer wg.Done()
			doc, err := scan(block, lineNum, path)
			if err != nil {
				log.Error(err)
				return
			}
			if doc == nil {
				return
			}
			docs.mux.Lock()
			docs.items = append(docs.items, doc)
			docs.mux.Unlock()
		}(block, ln, path)

		data = data[pos:]
	} // end for
	wg.Wait()

	return nil
}

// 扫描所有的paths文件内容。
func ScanFiles(paths []string, f ScanFunc) ([]*Doc, error) {
	docs := &docs{items: make([]*Doc, 0, 100)}
	wg := sync.WaitGroup{}
	for _, path := range paths {
		wg.Add(1)
		go func(p string) {
			scanFile(docs, f, p)
			wg.Done()
		}(path)
	}
	wg.Wait()

	return docs.items, nil
}
