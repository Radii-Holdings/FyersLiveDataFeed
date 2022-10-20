@echo off
FOR /F %%i IN ('python.exe "./Fyers-API-Access-Token-Generation-V2-main/New_Fyers_Access_Token.py"') DO set VARIABLE=%%i
echo %VARIABLE%
go run websocket.go %VARIABLE%