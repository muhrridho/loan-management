# Loan Management

## Overview
Loan Engine Example

## Tech Stack
- Golang 1.21
- Fiber (Web Framework)
- SQLite (Database)

## Getting Started

### Prerequisites
- Go 1.21
- SQLite3

### Installation

1. Clone the repository
```bash
git clone https://github.com/muhrridho/loan-management.git
cd loan-management
```

2. Install dependencies
```bash
go mod download
```

3. Run database migrations
```bash
go run main.go migrate
```

4. Seed the database with initial data
```bash
go run main.go seed
```

5. Start the application
```bash
go run main.go
```

The server will start on `http://localhost:8080` by default.

## Database Schema
See the `schema.txt` file for detailed database structure.

## API Documentation
A Postman collection is included with this repository for testing the API endpoints.

## Test Cases

### Test Case 1: Making a Payment

1. Create a loan with current date or past date
```bash
curl --location 'http://localhost:3000/api/loans/create' \
  --header 'Content-Type: application/json' \
  --data '{
    "user_id": 1,
    "amount": 5000000.00,
    "interest": 10.00,
    "interest_type": 0,
    "tenure": 52,
    "tenure_type": 0,
    "billing_start_date": "2025-02-18T00:00:00Z"
}'
```

2. Inquiry for due payment using loan ID
```bash
curl --location 'http://localhost:3000/api/transaction/inquiry?loan_id=1'
```

3. Create transaction using the loan_id and amount from inquiry
```bash
curl --location 'http://localhost:3000/api/transaction/create' \
  --header 'Content-Type: application/json' \
  --data '{
    "loan_id": 1,
    "amount": 317307.6923076923
  }'
```

4. Inquiry for due payment using loan ID (should return "No billing available")
```bash
curl --location 'http://localhost:3000/api/transaction/inquiry?loan_id=1'
```


### Test Case 2: Checking Outstanding Balance

1. Create a loan with current date or past date
```bash
curl --location 'http://localhost:3000/api/loans/create' \
  --header 'Content-Type: application/json' \
  --data '{
    "user_id": 1,
    "amount": 5000000.00,
    "interest": 10.00,
    "interest_type": 0,
    "tenure": 52,
    "tenure_type": 0,
    "billing_start_date": "2025-02-18T00:00:00Z"
}'
```

2. Check current outstanding balance
```bash
curl --location 'http://localhost:3000/api/loans/2'
```

3. Inquiry for due payment
```bash
curl --location 'http://localhost:3000/api/transaction/inquiry?loan_id=2'
```

4. Create transaction using the loan_id and amount from inquiry
```bash
curl --location 'http://localhost:3000/api/transaction/create' \
  --header 'Content-Type: application/json' \
  --data '{
      "loan_id": 1,
      "amount": 317307.6923076923
  }'
```

5. Check current outstanding balance again
```bash
curl --location 'http://localhost:3000/api/loans/2'
```

### Test Case 3: Checking Delinquent Status

0. Create new fresh user*
```bash
curl --location 'http://localhost:3000/api/users/register' \
--header 'Content-Type:  application/json' \
--form 'email="newuser@test"' \
--form 'name="new user"'
```

1. Check if user is delinquent
```bash
curl --location 'http://localhost:3000/api/users/2/delinquent-status'
```

2. Create a loan with past date
```bash
curl --location 'http://localhost:3000/api/loans/create' \
  --header 'Content-Type: application/json' \
  --data '{
    "user_id": 2,
    "amount": 5000000.00,
    "interest": 10.00,
    "interest_type": 0,
    "tenure": 52,
    "tenure_type": 0,
    "billing_start_date": "2025-02-18T00:00:00Z"
  }'
```

3. Check if user is delinquent again (should be true)
```bash
curl --location 'http://localhost:3000/api/users/2/delinquent-status'
```

4. Pay the due payment
```bash
curl --location 'http://localhost:3000/api/transaction/create' \
  --header 'Content-Type: application/json' \
  --data '{
    "loan_id": 3,
    "amount": 317307.6923076923
  }'
```

5. Check if user is delinquent after payment (should be false)
```bash
curl --location 'http://localhost:3000/api/users/2/delinquent-status'
```
