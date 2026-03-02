@echo off
chcp 65001 >nul
setlocal
cd /d %~dp0
echo 1. 正在安装依赖...
go mod tidy
echo 2. 正在编译项目...
go build -o main.exe ./cmd/main.go
echo 3. 正在运行项目...
main.exe
endlocal