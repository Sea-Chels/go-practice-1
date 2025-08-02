-- This file is executed when the PostgreSQL container is first initialized
-- Add any initial setup SQL here if needed

-- Ensure the database exists (this is handled by POSTGRES_DB env var, but included for clarity)
-- CREATE DATABASE school_db;

-- Grant all privileges to the user
GRANT ALL PRIVILEGES ON DATABASE school_db TO devuser;