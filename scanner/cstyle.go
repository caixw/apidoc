// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package scanner

func cstyle(s *scanner) ([]byte, error) {
	block := []rune{}

LOOP:
	for {
		switch {
		case s.match("/*"):
			for {
				if s.match("*/") {
					// TODO 两层转换，是否可以去掉
					return []byte(string(block)), nil
				}
				block = append(block, s.next())
			} // end for
		case s.match("//"):
		LOOP2:
			for {
				r := s.next()
				block = append(block, r)
				if r == '\n' {
					if s.match("//") {
						continue
					}
					break LOOP2
				}
			} // end for
			return []byte(string(block)), nil
		default:
			if s.next() == eof {
				break LOOP
			}
		}
	} // end for

	return nil, nil
}
