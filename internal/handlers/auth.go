package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/Sea-Chels/go-practice-1/internal/auth"
	"github.com/Sea-Chels/go-practice-1/internal/database"
	"github.com/Sea-Chels/go-practice-1/internal/models"
	"github.com/Sea-Chels/go-practice-1/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.ErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var loginReq models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		utils.ErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate input
	if err := utils.ValidateEmail(loginReq.Email); err != nil {
		utils.ErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get user from database
	var user models.User
	err := database.DB.QueryRow(`
		SELECT id, email, password_hash 
		FROM users 
		WHERE email = $1 AND deleted_at IS NULL
	`, loginReq.Email).Scan(&user.ID, &user.Email, &user.PasswordHash)

	if err == sql.ErrNoRows {
		utils.ErrorResponse(w, "Invalid credentials", http.StatusUnauthorized)
		return
	} else if err != nil {
		utils.ErrorResponse(w, "Database error", http.StatusInternalServerError)
		return
	}

	// Compare password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(loginReq.Password)); err != nil {
		utils.ErrorResponse(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Generate JWT token
	token, expiresAt, err := auth.GenerateToken(user.ID, user.Email)
	if err != nil {
		utils.ErrorResponse(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	response := models.LoginResponse{
		Token:     token,
		ExpiresAt: expiresAt,
	}

	utils.SuccessResponse(w, response, http.StatusOK)
}