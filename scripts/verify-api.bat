@echo off
setlocal enabledelayedexpansion
echo === HIS-Go 集成测试快速验证 ===
echo.

REM 1. Login
for /f "delims=" %%i in ('curl -s -X POST http://localhost:8080/api/auth/login -H "Content-Type: application/json" -d "{\"username\":\"demo-doctor\",\"password\":\"demo123\"}"') do set LOGIN_RESP=%%i
echo [1] Login: !LOGIN_RESP:~0,80!...

REM Extract token (hacky but works for demo)
for /f "tokens=2 delims=:" %%i in ('echo !LOGIN_RESP! ^| findstr "token"') do set TOKEN_PART=%%i
for /f "tokens=1 delims=," %%i in ("!TOKEN_PART!") do set TOKEN=%%i
set TOKEN=!TOKEN:"=!
set TOKEN=!TOKEN: =!
echo Token: !TOKEN:~0,30!...

REM 2. Departments
for /f "delims=" %%i in ('curl -s http://localhost:8080/api/user/departments -H "Authorization: Bearer !TOKEN!"') do set DEPT_RESP=%%i
echo [2] Departments: !DEPT_RESP:~0,80!

REM 3. Schedules
for /f "delims=" %%i in ('curl -s "http://localhost:8080/api/registration/schedules?deptId=dept_001" -H "Authorization: Bearer !TOKEN!"') do set SCHED_RESP=%%i
echo [3] Schedules: !SCHED_RESP:~0,80!

REM 4. Registration list
for /f "delims=" %%i in ('curl -s "http://localhost:8080/api/registration/list" -H "Authorization: Bearer !TOKEN!"') do set REG_RESP=%%i
echo [4] Registrations: !REG_RESP:~0,80!

REM 5. Prescriptions
for /f "delims=" %%i in ('curl -s "http://localhost:8080/api/prescription/list?patientId=patient_001" -H "Authorization: Bearer !TOKEN!"') do set PRESC_RESP=%%i
echo [5] Prescriptions: !PRESC_RESP:~0,80!

REM 6. Bills
for /f "delims=" %%i in ('curl -s "http://localhost:8080/api/billing/list" -H "Authorization: Bearer !TOKEN!"') do set BILL_RESP=%%i
echo [6] Bills: !BILL_RESP:~0,80!

REM 7. Drugs
for /f "delims=" %%i in ('curl -s "http://localhost:8080/api/pharmacy/drugs" -H "Authorization: Bearer !TOKEN!"') do set DRUG_RESP=%%i
echo [7] Drugs: !DRUG_RESP:~0,80!

echo.
echo === 验证完成 ===
