🧾 Project Documentation — HRMS (Human Resource Management System)
📌 Project Title
HRMS Backend – A Scalable Employee Management & Attendance System (Built with Go + PostgreSQL)

🔍 Overview
The HRMS Backend is a fully functional, production-grade system that enables organizations to manage employees, attendance, and leave applications efficiently.
It’s built using Go (Golang) for high performance, PostgreSQL for reliable data persistence, and includes JWT-based authentication, role-based access control, and Swagger documentation for interactive API exploration.
All modules have been tested using Postman and a comprehensive Go smoke test suite, ensuring the backend meets enterprise reliability and concurrency standards.

🧰 Tech Stack
Layer
Technology
Programming Language
Go 1.23
Framework / Router
Gorilla Mux
Database
PostgreSQL 16
ORM
GORM
Authentication
JWT (golang-jwt v5)
Documentation
Swagger (swaggo/http-swagger)
Containerization
Docker + Docker Compose
Testing
Go testing framework + tests/smoke_test.go
Environment Management
.env + github.com/joho/godotenv


⚙️ System Architecture
+--------------------------------------------------------------+
|                       HRMS Backend (Go)                      |
|--------------------------------------------------------------|
|  Controllers (REST APIs)                                     |
|   ├── Auth Controller      → Handles registration & login     |
|   ├── Employee Controller  → CRUD operations for employees    |
|   ├── Attendance Controller→ Tracks attendance logs           |
|   └── Leave Controller     → Manages leave requests/approvals |
|--------------------------------------------------------------|
|  Services (Business Logic)                                   |
|   ├── JWT & Role Validation                                   |
|   ├── Concurrency-safe attendance updates                     |
|   └── Transaction-based leave approvals                       |
|--------------------------------------------------------------|
|  Database (PostgreSQL via GORM ORM)                          |
|   ├── users         (Stores login credentials & roles)        |
|   ├── employees     (Employee details)                        |
|   ├── attendance    (Attendance logs)                         |
|   └── leaves        (Leave requests with status)              |
|--------------------------------------------------------------|
|  Infrastructure                                              |
|   ├── Docker Compose (App + DB containers)                   |
|   ├── Environment Config (.env)                              |
|   └── Swagger UI (/docs)                                     |
+--------------------------------------------------------------+


🧩 Key Features
🔐 Authentication & Authorization
JWT-based authentication for secure session management


Role-based access control (HR vs Employee)


Token verification middleware for route protection



👥 Employee Management (HR Role)
HR can add, list, update, or delete employees


Maintains department, email, and salary information


Auto-linked with users table for login mapping



🕒 Attendance Management
Employees can mark daily attendance


HR can view all attendance records


Concurrent submissions handled using Go goroutines and Mutex



🌴 Leave Management
Employees can apply for leaves


HR can view and approve/reject leave requests


ACID-compliant operations via transactional updates



🧠 Concurrency Management (Advanced Go Concepts)
Goroutines handle simultaneous employee operations (attendance, leave requests)


sync.Mutex / sync.RWMutex prevent race conditions


Database connection pooling ensures efficient performance under load



🧾 Database Schema
Table
Columns
Description
users
id, username, password, role
Stores login details and roles
employees
id, name, email, department, salary
Employee details
attendance
id, employee_id, date, status
Attendance records
leaves
id, employee_id, start_date, end_date, reason, status
Leave requests


🧪 Testing & Quality Assurance
A complete smoke test suite (tests/smoke_test.go) validates:
✅ Database connection


✅ Table existence


✅ Authentication (Register/Login)


✅ Employee CRUD


✅ Attendance and Leave workflows


✅ Data persistence


✅ Concurrency handling


Example command:
go test ./tests -v

Output:
✅ Database connected successfully
✅ Auth endpoints OK
✅ Employee CRUD OK
✅ Attendance OK
✅ Leave OK
✅ Data persisted correctly
PASS


🌐 Swagger API Documentation
URL: http://localhost:8082/docs/


Interactive testing of endpoints


Includes BearerAuth security schema


Generated and validated Swagger JSON (/docs/doc.json)



🐳 Docker Setup
docker-compose.yml
services:
  db:
    image: postgres:16
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: somu9866254149
      POSTGRES_DB: hrmssmall
    ports:
      - "7000:5432"
    volumes:
      - db_data:/var/lib/postgresql/data

  app:
    build: .
    ports:
      - "8082:8082"
    env_file:
      - .env
    depends_on:
      - db

volumes:
  db_data:

Commands
docker-compose up --build

✅ Starts both app and DB containers
 ✅ Auto-connects via environment variables
 ✅ Exposes backend on port 8082 and DB on 7000

🧠 Concurrency Example
var mu sync.Mutex

func (s *AttendanceService) UpdateAttendance(att Attendance) error {
    mu.Lock()
    defer mu.Unlock()
    return s.db.Save(&att).Error
}

Ensures multiple users updating attendance concurrently don’t cause data races.

🧑‍🔬 Manual Testing (Postman)
Verified all major endpoints:
/auth/register → 201 Created


/auth/login → 200 OK (returns JWT)


/employees → CRUD ✅


/attendance → Employee submissions ✅


/leaves → Apply + Approve ✅



✅ Final Project Status
Module
Status
Dependency Setup
✅ Completed
Database Connection
✅ Verified
Swagger Documentation
✅ Functional
Auth & Role Middleware
✅ Working
CRUD Operations
✅ Working
Concurrency Handling
✅ Tested
Docker Integration
✅ Stable
Smoke Test Suite
✅ Passed
Postman Verification
✅ Successful


🧩 Future Enhancements
Add Payroll module (salary slips, bonuses)


Add Admin dashboard (analytics, KPIs)


Implement email notifications for leave status


Deploy on Render / Railway with CI/CD integration



💼 How to Run Locally
# 1. Clone repo
git clone https://github.com/somanadh98/hrms-backend.git
cd hrms-backend

# 2. Configure environment
cp .env.example .env

# 3. Start containers
docker-compose up --build

# 4. Visit Swagger
http://localhost:8082/docs/


🧠 Key Takeaways
Strong understanding of Go backend architecture


Practical application of goroutines, mutex, and transactions


Real-world use of JWT security and RESTful best practices


Full-stack deployment pipeline with Docker


Professional-grade testing and documentation



👤 Author
Somanadh K.
 B.Tech Student – Chaitanya (Deemed to be University)
 💼 Interests: Backend Engineering, Cloud Deployment, Full-stack Development
 📧 Email: [k.somanadh98@gmail.com]
 🌐 GitHub: https://github.com/somanadh98/hrms-small-.git
