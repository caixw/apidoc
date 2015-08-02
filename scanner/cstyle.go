// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package scanner

func cstyle(s *scanner) ([]rune, int, error) {
	block := []rune{}

LOOP:
	for {
		switch {
		case s.match("/*"):
			lineNum := s.lineNumber()
			for {
				if s.match("*/") {
					return block, lineNum, nil
				}
				block = append(block, s.next())
			} // end for
		case s.match("//"):
			lineNum := s.lineNumber()
		LOOP2:
			for {
				r := s.next()
				block = append(block, r)
				if r == '\n' {
					s.skipSpace()
					if s.match("//") {
						continue
					}
					break LOOP2
				}
			} // end for
			return block, lineNum, nil
		default:
			if s.next() == eof {
				break LOOP
			}
		}
	} // end for

	return nil, 0, nil
}
