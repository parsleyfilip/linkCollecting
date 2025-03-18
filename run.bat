@echo off
ECHO Starting LinkCollector...

REM Check if Go is installed
WHERE go >nul 2>nul
IF %ERRORLEVEL% NEQ 0 (
    ECHO Go is not installed. Please install Go to run this application.
    EXIT /B
)

ECHO Checking dependencies...
go mod download

REM Set MySQL connection variables if needed
REM Uncomment and modify these lines with your Strato MySQL credentials
REM SET DB_USER=your_strato_db_user
REM SET DB_PASSWORD=your_strato_db_password
REM SET DB_HOST=your_strato_db_host
REM SET DB_PORT=3306
REM SET DB_NAME=your_strato_db_name
REM SET SESSION_SECRET=your-secure-session-key

ECHO Running LinkCollector on http://localhost:8080
go run main.go config.go

PAUSE 