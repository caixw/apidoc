// SPDX-License-Identifier: MIT

package mock

import "net/http"

var _ http.Handler = &Mock{}
