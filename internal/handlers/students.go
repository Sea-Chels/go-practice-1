package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Sea-Chels/go-practice-1/internal/database"
	"github.com/Sea-Chels/go-practice-1/internal/models"
	"github.com/Sea-Chels/go-practice-1/internal/utils"
	"github.com/gorilla/mux"
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

func UpdateStudentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		utils.ErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var student models.Student
	if err := json.NewDecoder(r.Body).Decode(&student); err != nil {
		utils.ErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate input
	if student.ID <= 0 {
		utils.ErrorResponse(w, "Invalid student ID", http.StatusBadRequest)
		return
	}
	if err := utils.ValidateStudentName(student.Name); err != nil {
		utils.ErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := utils.ValidateGrade(student.Grade); err != nil {
		utils.ErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check if student exists
	var exists bool
	err := database.DB.QueryRow(`
		SELECT EXISTS(SELECT 1 FROM students WHERE id = $1 AND deleted_at IS NULL)
	`, student.ID).Scan(&exists)
	
	if err != nil {
		utils.ErrorResponse(w, "Database error", http.StatusInternalServerError)
		return
	}
	
	if !exists {
		utils.ErrorResponse(w, "Student not found", http.StatusNotFound)
		return
	}

	// Update the student
	now := time.Now()
	err = database.DB.QueryRow(`
		UPDATE students 
		SET name = $2, grade = $3, updated_at = $4
		WHERE id = $1 AND deleted_at IS NULL
		RETURNING id, name, grade, created_at, updated_at, deleted_at
	`, student.ID, student.Name, student.Grade, now).Scan(
		&student.ID, &student.Name, &student.Grade, 
		&student.CreatedAt, &student.UpdatedAt, &student.DeletedAt)

	if err != nil {
		utils.ErrorResponse(w, "Failed to update student", http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(w, student, http.StatusOK)
}

func DeleteStudentHandler(w http.ResponseWriter, r *http.Request) {
	
	// Get ID from mux vars or URL path
	vars := mux.Vars(r)
	studentID := vars["id"]
	
	// Fallback: extract from URL path if mux vars are empty
	if studentID == "" {
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) >= 3 && parts[1] == "students" {
			studentID = parts[2]
		}
	}
	
	// Validate input
	studentIDInt, err := strconv.Atoi(studentID)
	if err != nil || studentIDInt <= 0 {
		utils.ErrorResponse(w, "Invalid student ID", http.StatusBadRequest)
		return
	}
	
	// Check if student exists
	var exists bool
	err = database.DB.QueryRow(`
		SELECT EXISTS(SELECT 1 FROM students WHERE id = $1 AND deleted_at IS NULL)
	`, studentIDInt).Scan(&exists)
	
	if err != nil {
		utils.ErrorResponse(w, "Not found", http.StatusNotFound)
		return
	}
	
	if !exists {
		utils.ErrorResponse(w, "Student not found", http.StatusNotFound)
		return
	}

	// Soft delete the student
	now := time.Now()
	_, err = database.DB.Exec(`
		UPDATE students 
		SET deleted_at = $2, updated_at = $3
		WHERE id = $1 AND deleted_at IS NULL
	`, studentIDInt, now, now)

	if err != nil {
		utils.ErrorResponse(w, "Not found", http.StatusNotFound)
		log.Printf("DeleteStudentHandler: error=%v", err)
		return
	}

	utils.SuccessResponse(w, nil, http.StatusNoContent)
}