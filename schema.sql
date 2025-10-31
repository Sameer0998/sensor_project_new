-- Database schema for Sensor Data Processing System

-- Create database
CREATE DATABASE IF NOT EXISTS sensor_data;
USE sensor_data;

-- Create sensor_types table
CREATE TABLE IF NOT EXISTS sensor_types (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    description VARCHAR(255),
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL
);

-- Create sensor_data table
CREATE TABLE IF NOT EXISTS sensor_data (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    sensor_value FLOAT NOT NULL,
    sensor_type_id INT NOT NULL,
    id1 VARCHAR(10) NOT NULL,
    id2 INT NOT NULL,
    created_at BIGINT NOT NULL,
    FOREIGN KEY (sensor_type_id) REFERENCES sensor_types(id),
    INDEX idx_id1_id2 (id1, id2),
    INDEX idx_created_at (created_at),
    INDEX idx_sensor_type (sensor_type_id)
);

-- Create users table for authentication
CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    email VARCHAR(100) NOT NULL UNIQUE,
    created_at BIGINT NOT NULL DEFAULT (UNIX_TIMESTAMP()*1000),
    updated_at BIGINT NOT NULL DEFAULT (UNIX_TIMESTAMP()*1000)
);

-- Create roles table
CREATE TABLE IF NOT EXISTS roles (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    description VARCHAR(255)
);

-- Create user_roles table (many-to-many relationship)
CREATE TABLE IF NOT EXISTS user_roles (
    user_id INT NOT NULL,
    role_id INT NOT NULL,
    PRIMARY KEY (user_id, role_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE
);

-- Insert default roles
INSERT INTO roles (name, description) VALUES 
('admin', 'Administrator with full access'),
('user', 'Regular user with limited access'),

-- Insert default sensor types
INSERT INTO sensor_types (name, description, created_at, updated_at) VALUES 
('temperature', 'Temperature sensor in Celsius', 1761908040000, 1761908040000),
('humidity', 'Humidity sensor in percentage', 1761908040000, 1761908040000),
('pressure', 'Pressure sensor in hPa', 1761908040000, 1761908040000),
('light', 'Light sensor in lux', 1761908040000, 1761908040000);