:: SPDX-License-Identifier: MIT

set wd=%~dp0\..\cmd\apidoc

set mainPath=github.com\caixw\apidoc

set varsPath=%mainPath%\v6\internal\vars

set builddate=%date:~0,4%%date:~5,2%%date:~8,2% 

for /f "delims=" %%t in ('git rev-parse HEAD') do set hash=%%t

echo compile
go build -o %wd%\apidoc.exe -ldflags "-X %varsPath%.buildDate=%builddate% -X %varsPath%.commitHash=%hash%" -v %wd%
