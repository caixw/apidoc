// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package locale

import (
	"syscall"
	"unsafe"
)

// GetUserDefaultLocaleName 第二个参数
const maxlen = 85

func getLocaleName() (string, error) {
	if name := getEnvLang(); len(name) > 0 {
		return name, nil
	}

	k32, err := syscall.LoadDLL("kernel32.dll")
	if err != nil {
		return "", err
	}
	defer k32.Release()

	f, err := k32.FindProc("GetUserDefaultLocaleName")
	if err != nil {
		return "", err
	}

	buf := make([]uint16, maxlen)
	r1, _, err := f.Call(uintptr(unsafe.Pointer(&buf[0])), uintptr(maxlen))
	if uint32(r1) == 0 {
		return "", err
	}

	return syscall.UTF16ToString(buf), nil
}
