// SPDX-License-Identifier: MIT

package mock

import (
	"golang.org/x/text/message"

	"github.com/caixw/apidoc/v5/internal/locale"
)

// Error 带状态码的错误信息
type Error struct {
	Status  int
	Message string
}

func (err *Error) Error() string {
	return err.Message
}

func newError(status int, key message.Reference, v ...interface{}) error {
	return &Error{
		Status:  status,
		Message: locale.Sprintf(key, v...),
	}
}
