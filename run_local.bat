@echo off
REM Build and run the local server from the project root
cd /d %~dp0\cmd\local

go build -o local_server.exe main.go
if %errorlevel% neq 0 (
    echo Build failed.
    exit /b %errorlevel%
)

REM Run the server
local_server.exe
