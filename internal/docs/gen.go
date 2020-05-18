// SPDX-License-Identifier: MIT

//go:generate go run ./make_config.go
//go:generate go run ./make_types.go
//go:generate go run ../../cmd/apidoc/main.go build ../../docs/example

// make_static 需要将以上的内容打包到 static.go，所以要放在最后调用。
//go:generate go run ./make_static.go

package docs
