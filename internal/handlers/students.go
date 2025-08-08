package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Sea-Chels/go-practice-1/internal/database"
	"github.com/Sea-Chels/go-practice-1/internal/models"
	"github.com/Sea-Chels/go-practice-1/internal/utils"
)

func GetStudentsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.ErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check for include_deleted query parameter
	includeDeleted := r.URL.Query().Get("include_deleted") == "true"
	log.Printf("GetStudentsHandler: include_deleted=%v, URL=%s", includeDeleted, r.URL.String())
	
	var query string
	if includeDeleted {
		query = `
			SELECT id, name, grade, created_at, updated_at, deleted_at 
			FROM students 
			ORDER BY id
		`
		log.Println("Including deleted students in query")
	} else {
		query = `
			SELECT id, name, grade, created_at, updated_at, deleted_at 
			FROM students 
			WHERE deleted_at IS NULL
			ORDER BY id
		`
		log.Println("Excluding deleted students from query")
	}

	rows, err := database.DB.Query(query)
	if err != nil {
		log.Printf("Database query error: %v", err)
		utils.ErrorResponse(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	
	log.Println("Query executed successfully")

	var students []models.Student
	for rows.Next() {
		var student models.Student
		err := rows.Scan(&student.ID, &student.Name, &student.Grade, 
			&student.CreatedAt, &student.UpdatedAt, &student.DeletedAt)
		if err != nil {
			utils.ErrorResponse(w, "Error scanning results", http.StatusInternalServerError)
			return
		}
		log.Printf("Student: ID=%d, Name=%s, DeletedAt=%v", student.ID, student.Name, student.DeletedAt)
		students = append(students, student)
	}

	if err = rows.Err(); err != nil {
		utils.ErrorResponse(w, "Database error", http.StatusInternalServerError)
		return
	}

	log.Printf("Total students found: %d", len(students))
	
	response := models.StudentsResponse{
		Students: students,
		Count:    len(students),
	}

	utils.SuccessResponse(w, response, http.StatusOK)
}

func CreateStudentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.ErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var createReq models.CreateStudentRequest
	if err := json.NewDecoder(r.Body).Decode(&createReq); err != nil {
		utils.ErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate input
	if err := utils.ValidateStudentName(createReq.Name); err != nil {
		utils.ErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := utils.ValidateGrade(createReq.Grade); err != nil {
		utils.ErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	var student models.Student
	now := time.Now()
	err := database.DB.QueryRow(`
		INSERT INTO students (name, grade, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id, name, grade, created_at, updated_at
	`, createReq.Name, createReq.Grade, now, now).Scan(
		&student.ID, &student.Name, &student.Grade, 
		&student.CreatedAt, &student.UpdatedAt)

	if err != nil {
		utils.ErrorResponse(w, "Failed to create student", http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(w, student, http.StatusCreated)
}