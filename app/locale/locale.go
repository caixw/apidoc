// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package locale

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// defaultTag 默认的语言
func Init(defaultTag string) {
	list := map[string]map[string]string{
		"zh-cmn-Hans": zh_cmn_Hans,
		"zh-cmn-Hant": zh_cmn_Hant,
	}

	tag, found := list[defaultTag]
	if !found {
		panic("参数 defaultTag 所指的语言不存在")
	}
	list["und"] = tag

	for id, messages := range list {
		tag := language.MustParse(id)
		for key, val := range messages {
			message.SetString(tag, key, val)
		}
	}
}
