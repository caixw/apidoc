:: Copyright 2017 by caixw, All rights reserved.
:: Use of this source code is governed by a MIT
:: license that can be found in the LICENSE file.


:: 代码的主目录，变量赋值时，等号两边不能有空格。
set wd=%~dp0\..

:: 程序所在的目录
set mainPath=github.com\caixw\apidoc\cmd

:: 需要修改变量的名名，若为 main，则指接使用 main，而不是全地址
set varsPath=%mainPath%\vars

:: 当前日期，格式为 YYYYMMDD
set builddate=%date:~0,4%%date:~5,2%%date:~8,2% 

:: 获取最后一条前交的 hash
for /f "delims=" %%t in ('git rev-parse HEAD') do set hash=%%t

echo 开始编译
%GOROOT%\bin\go build -o %wd%\apidoc.exe -ldflags "-X %varsPath%.buildDate=%builddate% -X %varsPath%.commitHash=%hash%" -v %mainPath%
