// SPDX-License-Identifier: MIT

//go:generate go run ./make_site.go
//go:generate go run ../../cmd/apidoc/ build -d=../../docs/example

// make_static 需要将以上的内容打包到 static.go，所以要放在最后调用。
//go:generate go run ./make_static.go

package docs
