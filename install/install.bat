@echo off

::判断系统
ver|findstr "3\.[0-9]\.[0-9][0-9]*">nul&&(set isXPlevel=1)
ver|findstr "4\.[0-9]\.[0-9][0-9]*">nul&&(set isXPlevel=1)
ver|findstr "5\.[0-9]\.[0-9][0-9]*">nul&&(set isXPlevel=1)
if "%isXPlevel%"=="" (set isXPlevel=0)


:CHOOSE
set /p xz=选择执行方式:(1.安装/2卸载):
echo.
if /i "%xz%"=="1" goto INSTALL
if /i "%xz%"=="2" goto UNINSTALL
echo 选项错误,请重新输入
echo.
echo.

pause 
goto CHOOSE

:INSTALL
:: @echo off

echo 开始安装

if "%isXPlevel%"=="1" (
	set conf=bin\conf.json
) else (
	set conf=%~dp0bin\conf.json
)

for /f "tokens=1* delims=" %%a in (%conf%) do (
set work_path=%%a
goto OUT_FOR)
:OUT_FOR

echo 安装目录:%work_path%

set bin_path=%work_path%\bin
set srv_exe="%bin_path%\IK_Auth_Srv.exe"
set srv_name=IK认证服务器


IF NOT EXIST "%work_path%" MD "%work_path%"
IF NOT EXIST "%bin_path%" MD "%bin_path%"

if "%isXPlevel%"=="1" (
	copy install.bat "%work_path%" > nul
	xcopy bin "%bin_path%" /s /e /i /y > nul
) else (
	copy %~dp0install.bat "%work_path%" > nul
	xcopy %~dp0bin "%bin_path%" /s /e /i /y > nul
)


%srv_exe% test
IF %ERRORLEVEL% neq 0 (
echo 初始化安装失败
goto END
)

REM 安装服务
%srv_exe% install
IF %ERRORLEVEL% neq 0 (
echo 安装服务失败
goto END
)


sc failure %srv_name% reset=0 actions= restart/5000 > nul

netsh firewall add allowedprogram %srv_exe% ENABLE > nul

sc start %srv_name% > nul
IF %ERRORLEVEL% neq 0 (
echo 安装失败
%srv_exe% remove > nul
goto END
)

echo 服务安装成功
pause
exit
:UNINSTALL
echo 开始卸载


if "%isXPlevel%"=="1" (
	set srv_exe=bin\IK_Auth_Srv.exe
) else (
	set srv_exe="%~dp0bin\IK_Auth_Srv.exe"
)
set srv_name=IK认证服务器

sc stop %srv_name%
IF %ERRORLEVEL% neq 0 (
	IF %ERRORLEVEL% neq 1062 (
		echo 服务停止失败
		goto END
	)
)

%srv_exe% remove
IF %ERRORLEVEL% neq 0 (
echo 服务卸载失败
goto END
)

ping -n 2 127.0.0.1>nul
::del "%work_path%" /f /s /q /a > nul
::rd "%work_path%" /s /q > nul
if "%isXPlevel%"=="1" (
	del bin /f /s /q /a > nul
	rd install.bat /s /q > nul
) else (
	del "%~dp0" /f /s /q /a > nul
	rd "%~dp0" /s /q > nul
)
echo 删除根目录
:END
pause
exit