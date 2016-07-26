// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// +build !windows

package locale

func getLocaleName() (string, error) {
	return getEnvLang(), nil
}
