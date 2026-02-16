@echo off
REM Run all Go tests from the project root
cd /d %~dp0

go test ./...
if %errorlevel% neq 0 (
    echo Tests failed.
    exit /b %errorlevel%
)
echo All tests passed.
