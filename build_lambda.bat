@echo off
REM Build the lambda handler for deployment from the project root
cd /d %~dp0\cmd\lambda

$env:GOOS="linux"
$env:GOARCH="amd64"
go build -o bootstrap main.go

if %errorlevel% neq 0 (
    echo Build failed.
    exit /b %errorlevel%
)

tar.exe -a -c -f function.zip bootstrap

echo Lambda handler built as bootstrap
