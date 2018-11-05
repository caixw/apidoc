// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package openapi

import "github.com/caixw/apidoc/errors"

// Sanitizer 数据验证接口
type Sanitizer interface {
	Sanitize() *errors.Error
}
