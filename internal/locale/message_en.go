// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package locale

import (
	"golang.org/x/text/language"
)

func init() {
	locales[language.MustParse("en")] = map[string]string{}
}
