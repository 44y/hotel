@echo off

::�ж�ϵͳ
ver|findstr "3\.[0-9]\.[0-9][0-9]*">nul&&(set isXPlevel=1)
ver|findstr "4\.[0-9]\.[0-9][0-9]*">nul&&(set isXPlevel=1)
ver|findstr "5\.[0-9]\.[0-9][0-9]*">nul&&(set isXPlevel=1)
if "%isXPlevel%"=="" (set isXPlevel=0)


:CHOOSE
set /p xz=ѡ��ִ�з�ʽ:(1.��װ/2ж��):
echo.
if /i "%xz%"=="1" goto INSTALL
if /i "%xz%"=="2" goto UNINSTALL
echo ѡ�����,����������
echo.
echo.

pause 
goto CHOOSE

:INSTALL
:: @echo off

echo ��ʼ��װ

if "%isXPlevel%"=="1" (
	set conf=bin\conf.json
) else (
	set conf=%~dp0bin\conf.json
)

for /f "tokens=1* delims=" %%a in (%conf%) do (
set work_path=%%a
goto OUT_FOR)
:OUT_FOR

echo ��װĿ¼:%work_path%

set bin_path=%work_path%\bin
set srv_exe="%bin_path%\IK_Auth_Srv.exe"
set srv_name=IK��֤������


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
echo ��ʼ����װʧ��
goto END
)

REM ��װ����
%srv_exe% install
IF %ERRORLEVEL% neq 0 (
echo ��װ����ʧ��
goto END
)


sc failure %srv_name% reset=0 actions= restart/5000 > nul

netsh firewall add allowedprogram %srv_exe% ENABLE > nul

sc start %srv_name% > nul
IF %ERRORLEVEL% neq 0 (
echo ��װʧ��
%srv_exe% remove > nul
goto END
)

echo ����װ�ɹ�
pause
exit
:UNINSTALL
echo ��ʼж��


if "%isXPlevel%"=="1" (
	set srv_exe=bin\IK_Auth_Srv.exe
) else (
	set srv_exe="%~dp0bin\IK_Auth_Srv.exe"
)
set srv_name=IK��֤������

sc stop %srv_name%
IF %ERRORLEVEL% neq 0 (
	IF %ERRORLEVEL% neq 1062 (
		echo ����ֹͣʧ��
		goto END
	)
)

%srv_exe% remove
IF %ERRORLEVEL% neq 0 (
echo ����ж��ʧ��
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
echo ɾ����Ŀ¼
:END
pause
exit