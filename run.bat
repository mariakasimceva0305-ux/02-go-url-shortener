@echo off
setlocal
if not exist go.mod (
  echo Run this script from project root.
  exit /b 1
)
set "GO_BIN=go"
if exist "C:\Program Files\Go\bin\go.exe" (
  set "GO_BIN=C:\Program Files\Go\bin\go.exe"
)
start "" "http://localhost:8080/swagger"
"%GO_BIN%" mod tidy
"%GO_BIN%" run ./cmd/app
