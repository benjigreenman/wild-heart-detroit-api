@echo off
REM Build the lambda handler for deployment from the project root
cd /d %~dp0\cmd\lambda

go build -o lambda_handler.exe
if %errorlevel% neq 0 (
    echo Build failed.
    exit /b %errorlevel%
)

tar.exe -a -c -f lambda_handler.zip lambda_handler.exe

echo Lambda handler built as lambda_handler.exe
