-- Migration: add_teachers_table
-- Created at: Thu Aug  7 21:12:00 CDT 2025

CREATE TABLE IF NOT EXISTS teachers (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    subject VARCHAR(100),
    hire_date DATE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Create an index on email for faster lookups
CREATE INDEX IF NOT EXISTS idx_teachers_email ON teachers(email);

-- Create an index on deleted_at for soft delete queries
CREATE INDEX IF NOT EXISTS idx_teachers_deleted_at ON teachers(deleted_at);