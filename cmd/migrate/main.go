package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/Sea-Chels/go-practice-1/internal/database"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Initialize database connection
	if err := database.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.CloseDB()

	// Get migration directory
	migrationsDir := "migrations"
	if len(os.Args) > 1 {
		migrationsDir = os.Args[1]
	}

	// Check if migrations directory exists
	if _, err := os.Stat(migrationsDir); os.IsNotExist(err) {
		log.Fatalf("Migrations directory '%s' does not exist", migrationsDir)
	}

	// Get all migration files
	files, err := filepath.Glob(filepath.Join(migrationsDir, "*.sql"))
	if err != nil {
		log.Fatalf("Failed to read migration files: %v", err)
	}

	if len(files) == 0 {
		fmt.Println("No migration files found")
		return
	}

	// Create migrations table if it doesn't exist
	createMigrationsTable := `
		CREATE TABLE IF NOT EXISTS migrations (
			id SERIAL PRIMARY KEY,
			filename VARCHAR(255) NOT NULL UNIQUE,
			executed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`
	if _, err := database.DB.Exec(createMigrationsTable); err != nil {
		log.Fatalf("Failed to create migrations table: %v", err)
	}

	// Track migrations executed
	migrationsRun := 0
	migrationsSkipped := 0

	// Run each migration
	for _, file := range files {
		filename := filepath.Base(file)

		// Check if migration has already been executed
		var count int
		err := database.DB.QueryRow("SELECT COUNT(*) FROM migrations WHERE filename = $1", filename).Scan(&count)
		if err != nil {
			log.Printf("Failed to check migration status for %s: %v", filename, err)
			continue
		}

		if count > 0 {
			log.Printf("Skipping already executed migration: %s", filename)
			migrationsSkipped++
			continue
		}

		// Read migration file
		content, err := os.ReadFile(file)
		if err != nil {
			log.Printf("Failed to read migration file %s: %v", filename, err)
			continue
		}

		// Execute migration
		if _, err := database.DB.Exec(string(content)); err != nil {
			log.Printf("Failed to execute migration %s: %v", filename, err)
			log.Fatalf("Migration failed. Stopping execution.")
		}

		// Record successful migration
		if _, err := database.DB.Exec("INSERT INTO migrations (filename) VALUES ($1)", filename); err != nil {
			log.Printf("Failed to record migration %s: %v", filename, err)
		}

		log.Printf("Successfully executed migration: %s", filename)
		migrationsRun++
	}

	// Summary
	fmt.Printf("\nMigration Summary:\n")
	fmt.Printf("  Migrations run: %d\n", migrationsRun)
	fmt.Printf("  Migrations skipped: %d\n", migrationsSkipped)
	fmt.Printf("  Total migrations: %d\n", len(files))

	if migrationsRun > 0 {
		fmt.Println("\nAll pending migrations have been executed successfully!")
	} else {
		fmt.Println("\nNo new migrations to run.")
	}
}