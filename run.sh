#!/bin/bash

# Simple script to run the LinkCollector application

echo "Starting LinkCollector..."

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "Go is not installed. Please install Go to run this application."
    exit 1
fi

# Download dependencies if needed
echo "Checking dependencies..."
go mod download

# Set MySQL connection variables if needed
# Uncomment and modify these lines with your Strato MySQL credentials
# export DB_USER=your_strato_db_user
# export DB_PASSWORD=your_strato_db_password
# export DB_HOST=your_strato_db_host
# export DB_PORT=3306
# export DB_NAME=your_strato_db_name
# export SESSION_SECRET=your-secure-session-key

# Run the application
echo "Running LinkCollector on http://localhost:8080"
go run main.go config.go 