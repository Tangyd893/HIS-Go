@echo off
REM 加载 deploy/config/demo.env 后执行命令（用于本地 go run 微服务）
setlocal enabledelayedexpansion
set "ENV=%~dp0..\deploy\config\demo.env"
if exist "%ENV%" (
  for /f "usebackq eol=# tokens=1,* delims==" %%a in ("%ENV%") do (
    if not "%%a"=="" set "%%a=%%b"
  )
) else (
  echo [WARN] %ENV% 不存在，就诊助手将仅使用关键词模式
)
cd /d "%~dp0..\backend"
%*
