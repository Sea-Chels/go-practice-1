package database

import (
	"fmt"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func SeedDatabase() error {
	// Check if data already exists
	var userCount int
	err := DB.QueryRow("SELECT COUNT(*) FROM users WHERE deleted_at IS NULL").Scan(&userCount)
	if err != nil {
		return fmt.Errorf("failed to check existing users: %w", err)
	}

	if userCount > 0 {
		log.Println("Database already contains data, skipping seed")
		return nil
	}

	// Seed admin user
	adminPassword := "Admin123!"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(adminPassword), 14)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	_, err = DB.Exec(`
		INSERT INTO users (email, password_hash, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
	`, "admin@example.com", string(hashedPassword), time.Now(), time.Now())
	if err != nil {
		return fmt.Errorf("failed to insert admin user: %w", err)
	}

	log.Println("Admin user created (email: admin@example.com, password: Admin123!)")

	// Seed students
	students := []struct {
		name  string
		grade int
	}{
		{"Alice Johnson", 10},
		{"Bob Smith", 11},
		{"Charlie Brown", 9},
		{"Diana Prince", 12},
		{"Ethan Hunt", 10},
	}

	for _, student := range students {
		_, err := DB.Exec(`
			INSERT INTO students (name, grade, created_at, updated_at)
			VALUES ($1, $2, $3, $4)
		`, student.name, student.grade, time.Now(), time.Now())
		if err != nil {
			return fmt.Errorf("failed to insert student %s: %w", student.name, err)
		}
	}

	log.Printf("Seeded %d students successfully", len(students))
	return nil
}