@echo off
FOR /F %%i IN ('python.exe "./Fyers-API-Access-Token-Generation-V2-main/New_Fyers_Access_Token.py"') DO set VARIABLE=%%i
echo %VARIABLE%
start "CopierService" cmd /k python pyDatasetCopier.py
go run websocket.go %VARIABLE%
taskkill /FI "WindowTitle eq CopierService*" /T /F