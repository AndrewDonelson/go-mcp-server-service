@echo off
setlocal EnableDelayedExpansion

rem Build configuration
set "BINARY_NAME=notes-server"
set "SERVICE_NAME=notes-service"
set "BUILD_DIR=bin"
set "VERSION=0.1.0"

rem Command processing
if "%~1"=="" goto help
if "%~1"=="clean" goto clean
if "%~1"=="dev" goto dev
if "%~1"=="release" goto release
if "%~1"=="release-windows" goto release-windows
if "%~1"=="help" goto help
echo Unknown command: %~1
goto help

:clean
echo Cleaning build directory...
if exist "%BUILD_DIR%" (
    rd /s /q "%BUILD_DIR%"
    if errorlevel 1 goto error
    echo Cleaned %BUILD_DIR% directory
) else (
    echo No %BUILD_DIR% directory found
)
goto :eof

:dev
echo Building development version...
if not exist "%BUILD_DIR%\dev\windows" mkdir "%BUILD_DIR%\dev\windows"
if errorlevel 1 goto error

echo Building command line app...
go build -o "%BUILD_DIR%\dev\windows\%BINARY_NAME%.exe" .\cmd
if errorlevel 1 goto error

echo Building service...
go build -o "%BUILD_DIR%\dev\windows\%SERVICE_NAME%.exe" .\service
if errorlevel 1 goto error
goto :eof

:release
call :release-windows
goto :eof

:release-windows
echo Building Windows release version...
if not exist "%BUILD_DIR%\release\windows" mkdir "%BUILD_DIR%\release\windows"
if errorlevel 1 goto error

echo Building command line app...
go build -o "%BUILD_DIR%\release\windows\%BINARY_NAME%.exe" .\cmd
if errorlevel 1 goto error

echo Building service...
go build -o "%BUILD_DIR%\release\windows\%SERVICE_NAME%.exe" .\service
if errorlevel 1 goto error
goto :eof

:help
echo.
echo Build script for the notes server project
echo.
echo Usage:
echo   build.bat [command]
echo.
echo Available commands:
echo   clean           - Remove build artifacts
echo   dev             - Build development version
echo   release         - Build release version
echo   release-windows - Build Windows release version
echo   help           - Show this help
goto :eof

:error
echo.
echo Build failed with error #%errorlevel%
exit /b %errorlevel%

:eof
exit /b 0