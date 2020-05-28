// SPDX-License-Identifier: MIT

package mock

const indent = "    "

var testOptions = &GenOptions{
	Number:    func() interface{} { return 1024 },
	String:    func() string { return "1024" },
	Bool:      func() bool { return true },
	SliceSize: func() int { return 5 },
	Index:     func(max int) int { return 0 },
}
