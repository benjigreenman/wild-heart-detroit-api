@echo off
REM Build and run all Go tests in the project
cd /d %~dp0..

go test ./...
if %errorlevel% neq 0 (
    echo Tests failed.
    exit /b %errorlevel%
)
echo All tests passed.
