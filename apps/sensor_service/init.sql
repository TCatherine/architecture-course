-- Create the database if it doesn't exist
CREATE DATABASE sensor;

-- Connect to the database
\c sensor;

-- Create the sensors table
CREATE TABLE IF NOT EXISTS homes (
    home_id      SERIAL PRIMARY KEY,
    user_id      INT,
    name         VARCHAR(255) NOT NULL,
    city         VARCHAR(255),
    street       VARCHAR(255),
    num          INT,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
    );

CREATE TABLE IF NOT EXISTS sensors (
    id             SERIAL PRIMARY KEY,
    service_id     INT,
    home_id        INT
    );


-- Create indexes for common queries
CREATE INDEX IF NOT EXISTS idx_sensors_home_id ON sensors(home_id);
