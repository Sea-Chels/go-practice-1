package models

import (
	"time"
)

type Student struct {
	ID        int        `json:"id"`
	Name      string     `json:"name"`
	Grade     int        `json:"grade"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

type CreateStudentRequest struct {
	Name  string `json:"name"`
	Grade int    `json:"grade"`
}

type StudentsResponse struct {
	Students []Student `json:"students"`
	Count    int       `json:"count"`
}