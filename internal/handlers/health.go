package handlers

import (
	"net/http"

	"github.com/Sea-Chels/go-practice-1/internal/database"
	"github.com/Sea-Chels/go-practice-1/internal/utils"
)

type HealthResponse struct {
	Status   string `json:"status"`
	Database string `json:"database"`
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.ErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	dbStatus := "connected"
	if err := database.DB.Ping(); err != nil {
		dbStatus = "disconnected"
	}

	response := HealthResponse{
		Status:   "ok",
		Database: dbStatus,
	}

	statusCode := http.StatusOK
	if dbStatus != "connected" {
		statusCode = http.StatusServiceUnavailable
	}

	utils.SuccessResponse(w, response, statusCode)
}