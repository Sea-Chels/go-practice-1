package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

func main() {
	// Connect to database
	db, err := sql.Open("postgres", "host=localhost port=5432 user=devuser password=devpass123 dbname=school_db sslmode=disable")
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Query all students
	fmt.Println("=== ALL STUDENTS IN DATABASE ===")
	rows, err := db.Query("SELECT id, name, grade, created_at, deleted_at FROM students ORDER BY id")
	if err != nil {
		log.Fatal("Failed to query students:", err)
	}
	defer rows.Close()

	count := 0
	deletedCount := 0
	for rows.Next() {
		var id int
		var name string
		var grade int
		var createdAt time.Time
		var deletedAt sql.NullTime

		err := rows.Scan(&id, &name, &grade, &createdAt, &deletedAt)
		if err != nil {
			log.Fatal("Failed to scan row:", err)
		}

		count++
		status := "ACTIVE"
		if deletedAt.Valid {
			status = fmt.Sprintf("DELETED at %s", deletedAt.Time.Format("2006-01-02 15:04:05"))
			deletedCount++
		}

		fmt.Printf("ID: %d, Name: %s, Grade: %d, Status: %s\n", id, name, grade, status)
	}

	fmt.Printf("\nTotal students: %d (Active: %d, Deleted: %d)\n", count, count-deletedCount, deletedCount)
}