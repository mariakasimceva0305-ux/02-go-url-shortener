@echo off
setlocal
echo Starting demo script...
powershell -NoProfile -ExecutionPolicy Bypass -File "%~dp0record_demo_project_02.ps1"
echo Demo script finished.
pause
