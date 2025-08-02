-- Create database and user for local development
-- Run this script as a PostgreSQL superuser (usually 'postgres')

-- Create the database
CREATE DATABASE school_db;

-- Create the user
CREATE USER devuser WITH PASSWORD 'devpass123';

-- Grant all privileges on the database to the user
GRANT ALL PRIVILEGES ON DATABASE school_db TO devuser;

-- Connect to the school_db database
\c school_db;

-- Grant schema permissions
GRANT ALL ON SCHEMA public TO devuser;