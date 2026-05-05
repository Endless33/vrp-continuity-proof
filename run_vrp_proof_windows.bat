@echo off
setlocal

title VRP Continuity Proof Runner

echo ============================================================
echo VRP / Jumping VPN - Windows Proof Runner
echo ============================================================
echo.
echo This runner verifies reproducible VRP execution invariants.
echo.
echo Expected:
echo CLEAN: CONSISTENT
echo CHAOS: CONSISTENT
echo ATTACK: CONSISTENT
echo CONSENSUS: CONSISTENT
echo OVERALL VERDICT: CONTINUITY PRESERVED
echo.

where go >nul 2>nul
if errorlevel 1 (
    echo [ERROR] Go is not installed or not available in PATH.
    echo.
    echo Install Go from:
    echo https://go.dev/dl/
    echo.
    pause
    exit /b 1
)

echo [OK] Go detected:
go version
echo.

echo [RUN] Oracle unified proof runner
echo ------------------------------------------------------------
go run ./cmd/oracle_unified_proof_runner
echo ------------------------------------------------------------
echo.

echo [INFO] UDP continuity proof requires two terminals.
echo.
echo Terminal 1:
echo   go run ./cmd/udp_continuity_server
echo.
echo Terminal 2:
echo   go run ./cmd/udp_continuity_client
echo.
echo Expected UDP server output:
echo   [PATH CHANGE DETECTED]
echo   session_identity_preserved=true
echo   VERDICT: UDP CONTINUITY PRESERVED
echo.

echo ============================================================
echo VRP proof runner finished.
echo ============================================================
echo.
pause