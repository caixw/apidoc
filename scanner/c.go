// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package scanner

func C(data []byte) ([]rune, int) {
	s := &scanner{
		data:  data,
		pos:   0,
		width: 0,
	}

	block := []rune{}

LOOP:
	for {
		switch {
		case s.match("/*"):
			for {
				if s.match("*/") {
					return block, s.pos
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
			return block, s.pos
		default:
			if s.next() == eof {
				break LOOP
			}
		}
	} // end for

	return nil, -1
}
