#!/usr/bin/env bash
set -e
xdg-open "http://localhost:8080/swagger" >/dev/null 2>&1 || true
go mod tidy
go run ./cmd/app
