@echo off
chcp 65001 >nul
title HIS-Go 患者端演示 - 一键启动

echo ============================================
echo   HIS-Go 患者端演示 (小程序 + H5)
echo ============================================
echo.

set "ROOT=%~dp0.."
set "ENV_FILE=%ROOT%\deploy\config\demo.env"

REM ========== 1. 基础设施 (无需 RabbitMQ) ==========
echo [1/4] 启动 Docker 基础设施...
docker compose -f "%ROOT%\deploy\compose\demo-patient.yml" --env-file "%ENV_FILE%" up -d postgresql redis 2>nul
if %errorlevel% neq 0 (
    echo [错误] Docker 未运行，请先启动 Docker Desktop
    pause
    exit /b 1
)
echo   等待 PostgreSQL 就绪...
timeout /t 10 /nobreak >nul

REM ========== 2. 初始化数据库 ==========
echo [2/4] 初始化数据库...
docker exec his-postgres-demo psql -U his_admin -f /sql-init/init_all.sql >nul 2>&1
docker exec his-postgres-demo psql -U his_admin -f /sql-init/seed_data.sql >nul 2>&1
echo   数据库初始化完成

REM ========== 3. 启动后端微服务 (患者端 Profile) ==========
echo [3/4] 启动后端服务...
set "DB_PASSWORD=change_me_123"
set "REDIS_PASSWORD=change_me_456"

cd /d "%ROOT%\backend"
start "HIS-Gateway"       go run .\cmd\gateway
start "HIS-Auth"          go run .\cmd\auth
start "HIS-User"          go run .\cmd\user
start "HIS-Reg"           go run .\cmd\registration
start "HIS-Schedule"      go run .\cmd\schedule
start "HIS-Prescription"  go run .\cmd\prescription
start "HIS-Examination"   go run .\cmd\examination
start "HIS-Followup"      go run .\cmd\followup
start "HIS-HealthRecord"  go run .\cmd\health_record

echo   等待 Gateway 就绪...
for /l %%i in (1,1,15) do (
    curl -s http://localhost:8080/health >nul 2>&1 && goto :gateway_ready
    timeout /t 1 /nobreak >nul
)
:gateway_ready

REM ========== 4. 启动患者 H5 ==========
echo [4/4] 启动患者端 H5 (Vite)...
cd /d "%ROOT%\frontend\patient"
start "HIS-Patient-H5" cmd /c "npm run dev -- --host && pause"

REM ========== 完成 ==========
timeout /t 3 /nobreak >nul
echo.
echo ============================================
echo   患者端演示已启动！
echo.
echo   H5 地址: http://localhost:5174
echo   API:     http://localhost:8080
echo   账号:    demo-patient / demo123
echo.
echo   小程序: 微信开发者工具 → 导入
echo           %ROOT%\frontend\mp-webview
echo ============================================
pause >nul
start http://localhost:5174
