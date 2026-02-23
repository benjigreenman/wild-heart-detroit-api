@echo off
REM 

go build -o /cmd/local/local_server.exe
if %errorlevel% neq 0 (
    echo Build failed.
    exit /b %errorlevel%
)

REM Run the server
/cmd/local/local_server.exe
