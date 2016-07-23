// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package locale

func init() {
	locales["zh-cmn-Hans"] = map[string]string{
		SyntaxError:  "在[%v:%v]出现语法错误[%v]",
		OptionsError: "配置文件[%v]中配置项[%v]错误:[%v]",
	}
}
