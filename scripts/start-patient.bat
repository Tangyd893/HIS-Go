@echo off
title HIS-Go Patient Demo
setlocal enabledelayedexpansion

set "ROOT=%~dp0.."
set "BACKEND=%ROOT%\backend"
set "FRONTEND=%ROOT%\frontend\patient"
set "COMPOSE=%ROOT%\deploy\compose\demo-patient.yml"
set "ENV=%ROOT%\deploy\config\demo.env"

echo ============================================
echo   HIS-Go Patient Demo Launcher
echo ============================================
echo.

REM ===== Step 1: Docker Infrastructure (no RabbitMQ) =====
echo [1/4] Starting Docker infrastructure...
docker compose -f "%COMPOSE%" --env-file "%ENV%" up -d postgresql redis >nul 2>&1
if %errorlevel% neq 0 (
    echo [ERROR] Docker not running. Start Docker Desktop first.
    pause
    exit /b 1
)
echo   Waiting for PostgreSQL...
timeout /t 10 /nobreak >nul
echo   OK

REM ===== Step 2: Init Database =====
echo [2/4] Initializing database...
docker exec his-postgres-demo psql -U his_admin -f /sql-init/init_all.sql >nul 2>&1
docker exec his-postgres-demo psql -U his_admin -f /sql-init/seed_data.sql >nul 2>&1
echo   OK

REM ===== Step 3: Start Go Services (patient profile) =====
echo [3/4] Starting backend services...
set DB_PASSWORD=change_me_123
set REDIS_PASSWORD=change_me_456

start "Gateway"       /D "%BACKEND%" cmd /c "go run ./cmd/gateway"
start "Auth"          /D "%BACKEND%" cmd /c "go run ./cmd/auth"
start "User"          /D "%BACKEND%" cmd /c "go run ./cmd/user"
start "Registration"  /D "%BACKEND%" cmd /c "go run ./cmd/registration"
start "Schedule"      /D "%BACKEND%" cmd /c "go run ./cmd/schedule"
start "Prescription"  /D "%BACKEND%" cmd /c "go run ./cmd/prescription"
start "Examination"   /D "%BACKEND%" cmd /c "go run ./cmd/examination"
start "Followup"      /D "%BACKEND%" cmd /c "go run ./cmd/followup"
start "HealthRecord"  /D "%BACKEND%" cmd /c "go run ./cmd/health_record"

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

REM ===== Step 4: Start Patient H5 =====
echo [4/4] Starting patient H5 (Vite)...
start "Patient-H5" /D "%FRONTEND%" cmd /c "npm run dev -- --host"

echo.
echo ============================================
echo   Patient Demo Ready!
echo.
echo   H5       : http://localhost:5174
echo   API      : http://localhost:8080
echo   Login    : demo-patient / demo123
echo.
echo   MiniProgram : WeChat DevTools -^> Import
echo                %ROOT%\frontend\mp-webview
echo.
echo   Press any key to open H5 in browser...
echo ============================================
pause >nul
start http://localhost:5174
