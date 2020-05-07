// SPDX-License-Identifier: MIT

package makeutil

// PanicError 如果 err 不为 nil，则 panic
func PanicError(err error) {
	if err != nil {
		panic(err)
	}
}
