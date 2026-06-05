@echo off
chcp 65001 >nul
title HIS-Go 管理端演示 - 一键启动

echo ============================================
echo   HIS-Go 管理端本地演示
echo ============================================
echo.

set "ROOT=%~dp0.."
set "ENV_FILE=%ROOT%\deploy\config\demo.env"

REM ========== 1. 基础设施 ==========
echo [1/4] 启动 Docker 基础设施...
docker compose -f "%ROOT%\deploy\compose\demo-admin.yml" --env-file "%ENV_FILE%" up -d postgresql redis rabbitmq 2>nul
if %errorlevel% neq 0 (
    echo [错误] Docker 未运行，请先启动 Docker Desktop
    pause
    exit /b 1
)
echo   等待 PostgreSQL 就绪...
docker exec his-postgres-demo pg_isready -U his_admin -d his_auth >nul 2>&1
if %errorlevel% neq 0 (
    echo   等待中... (最多 30 秒)
    timeout /t 30 /nobreak >nul
)

REM ========== 2. 初始化数据库 ==========
echo [2/4] 初始化数据库...
docker exec his-postgres-demo psql -U his_admin -f /sql-init/init_all.sql >nul 2>&1
docker exec his-postgres-demo psql -U his_admin -f /sql-init/seed_data.sql >nul 2>&1
echo   数据库初始化完成

REM ========== 3. 启动后端微服务 ==========
echo [3/4] 启动后端服务...
set "DB_PASSWORD=change_me_123"
set "REDIS_PASSWORD=change_me_456"
set "RABBITMQ_PASSWORD=change_me_789"

cd /d "%ROOT%\backend"
start "HIS-Gateway"    go run .\cmd\gateway
start "HIS-Auth"       go run .\cmd\auth
start "HIS-User"       go run .\cmd\user
start "HIS-Reg"        go run .\cmd\registration
start "HIS-Schedule"   go run .\cmd\schedule
start "HIS-Clinic"     go run .\cmd\clinic
start "HIS-Prescription" go run .\cmd\prescription
start "HIS-Billing"    go run .\cmd\billing
start "HIS-Pharmacy"   go run .\cmd\pharmacy
start "HIS-System"     go run .\cmd\system

echo   等待 Gateway 就绪 (最多 15 秒)...
for /l %%i in (1,1,15) do (
    curl -s http://localhost:8080/health >nul 2>&1 && goto :gateway_ready
    timeout /t 1 /nobreak >nul
)
:gateway_ready

REM ========== 4. 启动前端 ==========
echo [4/4] 启动管理端前端 (Vite)...
cd /d "%ROOT%\frontend\admin"
start "HIS-Admin-Web" cmd /c "npm run dev && pause"

REM ========== 完成 ==========
timeout /t 3 /nobreak >nul
echo.
echo ============================================
echo   管理端演示已启动！
echo.
echo   前端: http://localhost:5173
echo   API:  http://localhost:8080
echo   账号: demo-doctor / demo123
echo.
echo   按任意键打开浏览器...
echo ============================================
pause >nul
start http://localhost:5173
