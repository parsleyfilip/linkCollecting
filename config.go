package main

import (
	"fmt"
	"os"
	"strconv"
)

// Config holds all the configuration for the application
type Config struct {
	// Database settings
	DBUser     string
	DBPassword string
	DBHost     string
	DBPort     string
	DBName     string
	
	// Server settings
	ServerPort int
	
	// Session settings
	SessionSecret string
}

// LoadConfig loads configuration from environment variables with fallbacks to default values
func LoadConfig() Config {
	config := Config{
		// Database defaults - CHANGE THESE FOR YOUR STRATO HOSTING
		DBUser:     getEnv("DB_USER", "dbu4004523"),
		DBPassword: getEnv("DB_PASSWORD", "Plesk3317_"),
		DBHost:     getEnv("DB_HOST", "your_db_host"),
		DBPort:     getEnv("DB_PORT", "3306"),
		DBName:     getEnv("DB_NAME", "dbs14018482"),
		
		// Server defaults
		ServerPort: getEnvAsInt("SERVER_PORT", 8080),
		
		// Session defaults
		SessionSecret: getEnv("SESSION_SECRET", "change-this-to-something-secure-in-production"),
	}
	
	return config
}

// GetDSN returns the database connection string
func (c *Config) GetDSN() string {
	// MySQL connection string format: username:password@tcp(host:port)/dbname?parseTime=true
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", 
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName)
}

// Helper function to get environment variable with a default value
func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}

// Helper function to parse int environment variables with a default value
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
} 