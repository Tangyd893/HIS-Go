@echo off
title HIS-Go Admin Demo
setlocal enabledelayedexpansion

set "ROOT=%~dp0.."
set "BACKEND=%ROOT%\backend"
set "FRONTEND=%ROOT%\frontend\admin"
set "COMPOSE=%ROOT%\deploy\compose\demo-admin.yml"
set "ENV=%ROOT%\deploy\config\demo.env"

echo ============================================
echo   HIS-Go Admin Demo Launcher
echo ============================================
echo.

REM ===== Step 1: Docker Infrastructure =====
echo [1/4] Starting Docker infrastructure...
docker compose -f "%COMPOSE%" --env-file "%ENV%" up -d postgresql redis rabbitmq >nul 2>&1
if %errorlevel% neq 0 (
    echo [ERROR] Docker not running. Start Docker Desktop first.
    pause
    exit /b 1
)
echo   Waiting for PostgreSQL...
docker exec his-postgres-demo pg_isready -U his_admin >nul 2>&1
if %errorlevel% neq 0 (
    timeout /t 20 /nobreak >nul
)
echo   OK

REM ===== Step 2: Init Database =====
echo [2/4] Initializing database...
docker exec his-postgres-demo psql -U his_admin -f /sql-init/init_all.sql >nul 2>&1
docker exec his-postgres-demo psql -U his_admin -f /sql-init/seed_data.sql >nul 2>&1
echo   OK

REM ===== Step 3: Start Go Services =====
echo [3/4] Starting backend services...
set DB_PASSWORD=change_me_123
set REDIS_PASSWORD=change_me_456
set RABBITMQ_PASSWORD=change_me_789

start "Gateway"       /D "%BACKEND%" cmd /c "go run ./cmd/gateway"
start "Auth"          /D "%BACKEND%" cmd /c "go run ./cmd/auth"
start "User"          /D "%BACKEND%" cmd /c "go run ./cmd/user"
start "Registration"  /D "%BACKEND%" cmd /c "go run ./cmd/registration"
start "Schedule"      /D "%BACKEND%" cmd /c "go run ./cmd/schedule"
start "Clinic"        /D "%BACKEND%" cmd /c "go run ./cmd/clinic"
start "Prescription"  /D "%BACKEND%" cmd /c "go run ./cmd/prescription"
start "Billing"       /D "%BACKEND%" cmd /c "go run ./cmd/billing"
start "Pharmacy"      /D "%BACKEND%" cmd /c "go run ./cmd/pharmacy"
start "System"        /D "%BACKEND%" cmd /c "go run ./cmd/system"

echo   Waiting for Gateway (port 8080)...
set /a TRIES=0
:wait_gw
timeout /t 2 /nobreak >nul
set /a TRIES+=1
curl -s http://localhost:8080/health >nul 2>&1
if %errorlevel% equ 0 goto gw_ok
if %TRIES% lss 15 goto wait_gw
echo   [WARN] Gateway may still be starting...
:gw_ok
echo   OK

REM ===== Step 4: Start Frontend =====
echo [4/4] Starting frontend (Vite)...
start "Admin-Web" /D "%FRONTEND%" cmd /c "npm run dev"

echo.
echo ============================================
echo   Admin Demo Ready!
echo.
echo   Frontend : http://localhost:5173
echo   API      : http://localhost:8080
echo   Login    : demo-doctor / demo123
echo.
echo   Press any key to open browser...
echo ============================================
pause >nul
start http://localhost:5173
