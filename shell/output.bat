@echo off

echo.
echo ====================================== Proxy Server ======================================

SET FILE_PATH=C:\Users\sd\proxy
if exist %FILE_PATH%\proxy (
    cd %FILE_PATH%
    del /s /Q %FILE_PATH%\proxy
    rd /s /Q %FILE_PATH%\conf
	echo %FILE_PATH%\proxy has been deleted!
	echo %FILE_PATH%\conf has been deleted!
)

::下载所有依赖包
::go mod download
::下载缺少的包，删除不用的包
::go mod tidy

::PROXY_PATH
SET PROXY_PATH=D:\slggame\gamesproxy\proxy\server
cd %PROXY_PATH%

::go build
SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
go build -o %FILE_PATH%\proxy %PROXY_PATH%\main.go

::copy conf
mkdir %FILE_PATH%\conf
xcopy /s/e/y %PROXY_PATH%\conf\* %FILE_PATH%\conf\


echo ====================================== Successfully ======================================
echo.
