-- PostgreSQL initialization script for Timelith
-- This file creates necessary extensions and initial settings

-- Enable required extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Grant necessary privileges
GRANT ALL PRIVILEGES ON DATABASE timelith_production TO timelith;
