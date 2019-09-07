// SPDX-License-Identifier: MIT

// +build !windows

package syslocale

func getLocaleName() (string, error) {
	return getEnvLang(), nil
}
