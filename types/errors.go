// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package types

import (
	"github.com/caixw/apidoc/locale"
	"github.com/caixw/apidoc/vars"
)

// OptionsError 提供对配置项错误的描述
type OptionsError struct {
	Field   string
	Message string
}

func (err *OptionsError) Error() string {
	return locale.Sprintf(locale.OptionsError, vars.ConfigFilename, err.Field, err.Message)
}
