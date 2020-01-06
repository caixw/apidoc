// SPDX-License-Identifier: MIT

// 以下 go:generate 有依赖关系，顺序不能乱！

//go:generate go run ./make_config.go
//go:generate go run ../../cmd/apidoc/main.go build ../../docs/example
//go:generate go run ./make_static.go

package docs
