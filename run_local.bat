@echo off
REM Build and run the local server from the project root


go build -o \cmd\locallocal_server.exe
if %errorlevel% neq 0 (
    echo Build failed.
    exit /b %errorlevel%
)

REM Run the server
\cmd\locallocal_server.exe
