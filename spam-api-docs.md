# Spam Detection API Documentation

## Overview
A RESTful API service that allows users to create accounts, report spam calls/messages, and search for spam reports by name or phone number.

## Technical Requirements
- Docker & Docker Compose
- Go (v1.18 or higher)
- PostgreSQL (handled by Docker)

## Quick Start

### 1. Using Docker (Recommended)
```bash
# Start all services
docker-compose up --build
```
The API will be available at `http://localhost:8080`

### 2. Development Mode
```bash
# Run directly with Go
go run main.go -e development
```

## API Reference

### Base URL
All API endpoints are prefixed with:
```
http://localhost:8080/api/v1
```

### 1. User Management

#### Create Account
- **Endpoint**: `/users/create`
- **Method**: `POST`
- **Body**:
```json
{
  "name": "John Doe",
  "email": "john.doe@example.com",
  "phone_number": "+123456789",
  "password": "yourpassword"
}
```
- **Success Response**: `200 OK`
```json
{
  "message": "User created successfully"
}
```
- **Error Cases**:
  - `409`: User already exists
  - `400`: Invalid input data

#### Login
- **Endpoint**: `/users/login`
- **Method**: `POST`
- **Body**:
```json
{
  "phone_number": "+123456789",
  "password": "yourpassword"
}
```
- **Success Response**: `200 OK`
```json
{
  "message": "login successful",
  "response": {
    "access_token": "eyJhbG..."
  }
}
```
- **Error Cases**:
  - `401`: Invalid credentials

### 2. Spam Management
**Note**: All spam endpoints require authentication. Include this header:
```
Authorization: Bearer your_access_token
```

#### Report Spam
- **Endpoint**: `/spam/report`
- **Method**: `POST`
- **Body**:
```json
{
  "name": "Spam Caller",
  "phone_number": "+123456789",
  "spam_likelihood": 75  // 0-100 scale
}
```
- **Success Response**: `200 OK`
```json
{
  "message": "Spam report created successfully"
}
```

#### Search by Name
- **Endpoint**: `/spam/search/name`
- **Method**: `GET`
- **Query Parameter**: `?name=John`
- **Success Response**: `200 OK`
```json
{
  "results": [
    {
      "name": "John Spammer",
      "phone_number": "+123456789",
      "spam_likelihood": 75,
      "reported_at": "2024-01-01T12:00:00Z"
    }
  ]
}
```

#### Search by Phone Number
- **Endpoint**: `/spam/search/phone`
- **Method**: `GET`
- **Query Parameter**: `?phone_number=+123456789`
- **Response Format**: Same as name search

## Example Usage

### Create a New User
```bash
curl -X POST http://localhost:8080/api/v1/users/create \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "phone_number": "+123456789",
    "password": "password123"
  }'
```

### Login
```bash
curl -X POST http://localhost:8080/api/v1/users/login \
  -H "Content-Type: application/json" \
  -d '{
    "phone_number": "+123456789",
    "password": "password123"
  }'
```

### Report Spam (Authenticated)
```bash
# Replace YOUR_TOKEN with the token received from login
curl -X POST http://localhost:8080/api/v1/spam/report \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Spam Caller",
    "phone_number": "+123456789",
    "spam_likelihood": 80
  }'
```

### Search for Spam Reports
```bash
# Search by name
curl -X GET "http://localhost:8080/api/v1/spam/search/name?name=John" \
  -H "Authorization: Bearer YOUR_TOKEN"

# Search by phone number
curl -X GET "http://localhost:8080/api/v1/spam/search/phone?phone_number=+123456789" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## Project Structure
```
├── controllers/
│   ├── user_controller.go
│   └── spam_controller.go
├── middleware/
│   └── auth_middleware.go
├── models/
│   ├── user.go
│   └── spam_report.go
├── main.go
└── docker-compose.yml
```

## Security Notes
- All passwords are hashed before storage
- JWT tokens are used for authentication
- Rate limiting is implemented on all endpoints
- Input validation is performed on all requests
