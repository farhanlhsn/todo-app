# Todo App API

> A RESTful API for todo/task management built with **Go**, **Gin**, **GORM**, and **JWT authentication**.

[![Go Version](https://img.shields.io/badge/Go-1.23+-blue.svg)](https://golang.org)
[![Gin Framework](https://img.shields.io/badge/Framework-Gin-green.svg)](https://gin-gonic.com)
[![Database](https://img.shields.io/badge/Database-MySQL-orange.svg)](https://mysql.com)

## 📋 Table of Contents

- [Features](#-features)
- [Architecture & Design Patterns](#️-architecture--design-patterns)
- [Technology Stack](#️-technology-stack)
- [Prerequisites](#-prerequisites)
- [Installation & Setup](#-installation--setup)
- [Testing](#-testing)
- [API Documentation](#-api-documentation)
- [Database Schema](#-database-schema)
- [Security Features](#-security-features)

## ✨ Features

### Core Features
- **User Authentication**: Register, Login, Logout with JWT
- **Task Management**: Full CRUD operations for tasks
- **Category Management**: Organize tasks by categories
- **Task Status Management**: Mark tasks as complete/incomplete
- **Advanced Filtering**: Filter by status, category, due date, priority
- **Search Functionality**: Search tasks by title/description
- **User Statistics**: Track completion rates and overdue tasks

### Technical Features
- **RESTful API**: Clean and consistent API design
- **JWT Authentication**: Secure token-based authentication
- **Rate Limiting**: Prevent API abuse and DDoS attacks
- **CORS Support**: Cross-origin resource sharing
- **Input Validation**: Comprehensive request validation
- **Database Optimization**: Indexed queries for better performance
- **Auto Migration**: Automatic database schema migration
- **E2E Testing**: Comprehensive end-to-end testing

## 🏗️ Architecture & Design Patterns

### **Primary Pattern: MVC (Model-View-Controller) with Layered Architecture**

#### **Why This Pattern?**

1. **🔍 Separation of Concerns**
   - Clear boundaries between data, business logic, and presentation
   - Each layer has specific responsibilities
   - Reduces coupling between components

2. **🔧 Maintainability**
   - Easy to locate and fix bugs
   - Simple to modify individual components
   - Clear code organization

3. **🧪 Testability**
   - Each layer can be tested independently
   - Mock dependencies easily
   - Comprehensive test coverage

4. **📈 Scalability**
   - Add new features without affecting existing code
   - Horizontal scaling capabilities
   - Performance optimization per layer

5. **👥 Team Collaboration**
   - Multiple developers can work on different layers
   - Clear interfaces between components
   - Reduced merge conflicts

### **Detailed Project Structure**

```
todo-app/
├── 📁 controllers/              # 🎮 Business Logic Layer
│   ├── userControllers.go           # User authentication & management
│   └── taskControllers.go           # Task & category operations
├── 📁 models/                   # 🗃️ Data Access Layer
│   ├── userModels.go                # User entity with relationships
│   ├── taskModels.go                # Task entity with business rules
│   └── taskCategoryModels.go        # Category entity with constraints
├── 📁 middlewares/              # 🛡️ Cross-Cutting Concerns
│   ├── requiredAuth.go              # JWT authentication middleware
│   ├── cors.go                      # CORS configuration
│   └── rateLimit.go                 # Rate limiting protection
├── 📁 initializers/             # ⚙️ Configuration Layer
│   ├── connectToDB.go               # Database connection setup
│   ├── syncDatabase.go              # Schema migration & seeding
│   ├── addIndexes.go                # Database optimization
│   └── loadEnvVariables.go          # Environment configuration
├── 📁 helpers/                  # 🔧 Utility Layer
│   └── responseHelpers.go           # Standardized API responses
├── 📁 tests/                    # 🧪 Testing Layer
│   ├── login_e2e_with_layers_test.go   # E2E layered testing
│   ├── 📁 assertions/               # Test assertion helpers
│   ├── 📁 services/                 # Test service layer
│   └── 📁 helpers/                  # Test utility functions
├── 📄 main.go                   # 🚀 Application Entry Point
├── 📄 go.mod                    # 📦 Dependency Management
├── 📄 .env                      # 🔐 Environment Variables
└── 📄 README.md                 # 📖 Documentation
```

### **Architecture Layers Explained**

#### 1. **Presentation Layer** (Gin Router + Middlewares)
- HTTP request handling
- Route definition and grouping
- Middleware chain execution
- Response formatting

#### 2. **Business Logic Layer** (Controllers)
- Request validation and sanitization
- Business rule implementation
- Data transformation
- Error handling

#### 3. **Data Access Layer** (Models + GORM)
- Database entity definitions
- Relationship management
- Query optimization
- Data validation

#### 4. **Infrastructure Layer** (Initializers + Helpers)
- Database connection management
- Configuration management
- Utility functions
- External service integration

### **Design Patterns Applied**

1. **🔄 Repository Pattern (Implicit)**
   - GORM acts as the repository layer
   - Abstracts database operations
   - Consistent data access interface

2. **🎯 Middleware Pattern**
   - Cross-cutting concerns (auth, CORS, rate limiting)
   - Request/response interception
   - Modular and reusable components

3. **💉 Dependency Injection**
   - Database connection injected through initializers
   - Testable and mockable dependencies
   - Loose coupling between components

4. **📋 Response Standardization**
   - Consistent API response format
   - Error handling standardization
   - Client-friendly responses

## 🛠️ Technology Stack

### **Backend Framework**
- **[Gin](https://gin-gonic.com)**: High-performance HTTP web framework
- **[GORM](https://gorm.io)**: Object-Relational Mapping library
- **[MySQL](https://mysql.com)**: Relational database management system

### **Security & Authentication**
- **[JWT](https://jwt.io)**: JSON Web Tokens for stateless authentication
- **[bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt)**: Password hashing algorithm
- **Rate Limiting**: API abuse prevention

### **Development & Testing**
- **[Testify](https://github.com/stretchr/testify)**: Testing toolkit
- **[godotenv](https://github.com/joho/godotenv)**: Environment variable management
- **[CompileDaemon](https://github.com/githubnemo/CompileDaemon)**: Hot reload development

### **Performance & Optimization**
- **Connection Pooling**: Optimized database connections
- **Database Indexing**: Query performance optimization
- **Caching Strategy**: Middleware caching implementation

## 📋 Prerequisites

### **Required Software**
- **Go**: Version 1.23 or higher
- **MySQL**: Version 8.0 or higher
- **Git**: For version control

### **Development Tools (Optional)**
- **Postman**: API testing and documentation
- **MySQL Workbench**: Database management
- **VS Code**: Recommended editor with Go extension

## 🚀 Installation & Setup

### **1. Clone Repository**
```bash
git clone https://github.com/your-username/todo-app.git
cd todo-app
```

### **2. Install Dependencies**
```bash
go mod download
go mod tidy
```

### **3. Environment Configuration**
Create a `.env` file in the root directory:

```env
# Database Configuration (Production)
DATABASE_URL=username:password@tcp(localhost:3306)/todo_app?charset=utf8mb4&parseTime=True&loc=Local

# Database Configuration (Testing)
TEST_DATABASE_URL=username:password@tcp(localhost:3306)/todo_app_test?charset=utf8mb4&parseTime=True&loc=Local

# JWT Configuration
JWT_SECRET_KEY=your-super-secure-jwt-secret-key-change-this-in-production

# Server Configuration
PORT=3000

# Environment Mode
GIN_MODE=debug
```

### **4. Database Setup**

#### **Production Database**
```bash
mysql -u root -p
CREATE DATABASE todo_app CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

#### **Testing Database**
```bash
mysql -u root -p
CREATE DATABASE todo_app_test CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

### **5. Run Application**

#### **Development Mode**
```bash
# With hot reload
go run main.go

# Or with CompileDaemon for auto-restart
CompileDaemon -command="./todo-app"
```

#### **Production Mode**
```bash
# Build and run
go build -o todo-app
./todo-app
```

### **6. Verify Installation**
- Application: `http://localhost:3000`
- Health Check: `http://localhost:3000/api/v1/health`

## 🧪 Testing

### **Testing Strategy**

This project implements **layered E2E testing** following best practices:

1. **Layer 1**: Registration with Database Integration
2. **Layer 2**: Login Flow with JWT Generation
3. **Layer 3**: Protected Route Access with Authentication
4. **Layer 4**: Middleware Authentication Verification
5. **Layer 5**: Invalid Credentials & Error Handling

### **Running Tests**

#### **Run All Tests**
```bash
go test ./tests/ -v
```

#### **Run Specific Test Layer**
```bash
# Test only Layer 1 (Registration)
go test ./tests/ -run TestLoginE2E_WithLayeringPattern/Layer_1 -v

# Test only Layer 2 (Login)
go test ./tests/ -run TestLoginE2E_WithLayeringPattern/Layer_2 -v

# Test only Layer 3 (Protected Routes)
go test ./tests/ -run TestLoginE2E_WithLayeringPattern/Layer_3 -v
```

#### **Test with Coverage**
```bash
go test ./tests/ -cover -v
```

#### **Test with Timeout**
```bash
go test ./tests/ -timeout 30s -v
```

### **Test Architecture**

```
tests/
├── login_e2e_with_layers_test.go    # Main E2E test implementation
├── assertions/                      # Test assertion helpers
│   └── auth_assertions.go
├── services/                        # Test service layer
│   ├── auth_service.go
│   └── user_service.go
└── helpers/                         # Test utility functions
    └── http_client.go
```

### **Testing Best Practices Implemented**

1. **🗃️ Database Isolation**: Each test uses separate test database
2. **🧹 Cleanup Strategy**: Automatic cleanup after each test
3. **🔄 Test Independence**: Each layer can run independently
4. **📊 Comprehensive Coverage**: All authentication flows tested
5. **🛡️ Security Testing**: Invalid credentials and unauthorized access
6. **⚡ Performance Testing**: Database optimization verification

## 📚 API Documentation

### **Base URL**
```
http://localhost:3000/api/v1
```

### **Response Format**
All API responses follow a consistent format:

```json
{
  "success": true,
  "message": "Operation successful",
  "data": {
    // Response data here
  }
}
```

### **Authentication Flow**

#### **1. Register User**
```http
POST /api/v1/auth/register
Content-Type: application/json

{
  "nama": "John Doe",
  "email": "john@example.com",
  "password": "password123"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "User registered successfully"
}
```

#### **2. Login User**
```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "john@example.com",
  "password": "password123"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Login successful"
}
```
*Note: JWT token is set as HTTP-only cookie*

#### **3. Logout User**
```http
POST /api/v1/auth/logout
Cookie: Authorization=<jwt_token>
```

### **User Management**

#### **Get User Profile**
```http
GET /api/v1/users/profile
Cookie: Authorization=<jwt_token>
```

#### **Get User Statistics**
```http
GET /api/v1/users/stats
Cookie: Authorization=<jwt_token>
```

### **Task Management**

#### **Create Task**
```http
POST /api/v1/task/
Cookie: Authorization=<jwt_token>
Content-Type: application/json

{
  "title": "Complete project documentation",
  "description": "Write comprehensive API documentation",
  "due_date": "2024-12-31T23:59:59Z",
  "category_id": 1,
  "priority": "high"
}
```

#### **Get All Tasks**
```http
GET /api/v1/task/
Cookie: Authorization=<jwt_token>
```

#### **Update Task**
```http
PUT /api/v1/task/{id}
Cookie: Authorization=<jwt_token>
Content-Type: application/json

{
  "title": "Updated task title",
  "description": "Updated description",
  "priority": "medium"
}
```

#### **Delete Task**
```http
DELETE /api/v1/task/{id}
Cookie: Authorization=<jwt_token>
```

#### **Task Status Management**
```http
# Mark as complete
PUT /api/v1/task/{id}/complete
Cookie: Authorization=<jwt_token>

# Mark as incomplete
PUT /api/v1/task/{id}/uncomplete
Cookie: Authorization=<jwt_token>
```

#### **Task Filtering & Search**
```http
# Get completed tasks
GET /api/v1/task/completed
Cookie: Authorization=<jwt_token>

# Get pending tasks
GET /api/v1/task/pending
Cookie: Authorization=<jwt_token>

# Get overdue tasks
GET /api/v1/task/overdue
Cookie: Authorization=<jwt_token>

# Search tasks
GET /api/v1/task/search?q=project&priority=high
Cookie: Authorization=<jwt_token>
```

### **Category Management**

#### **Create Category**
```http
POST /api/v1/task/categories/
Cookie: Authorization=<jwt_token>
Content-Type: application/json

{
  "name": "Work Projects"
}
```

#### **Get Tasks by Category**
```http
GET /api/v1/task/categories/{name}
Cookie: Authorization=<jwt_token>
```

### **Import Postman Collection**

1. Download the `Todo_App_API.postman_collection.json` file
2. Open Postman
3. Click "Import" button
4. Select the JSON file
5. Collection will be imported with all endpoints and environment variables

## 🗃️ Database Schema

### **Entity Relationship Diagram**

```
Users (1) ──────── (N) Tasks (N) ──────── (1) TaskCategories
  │                    │                      │
  │ ID                 │ UserID               │ ID
  │ Nama               │ CategoryID           │ Name
  │ Email              │ Title                │ UserID
  │ Password           │ Description          │ IsDefault
  │ IsLoggedIn         │ IsCompleted          │
  │                    │ DueDate              │
  │                    │ Priority             │
```

### **Database Optimizations**

- **Indexes**: Optimized queries with strategic indexing
- **Connection Pooling**: Efficient connection management
- **Foreign Keys**: Data integrity constraints
- **Soft Deletes**: Data preservation with GORM
- **UTF8MB4**: Full Unicode support

## 🔒 Security Features

### **Authentication & Authorization**
- ✅ **JWT Tokens**: Stateless authentication
- ✅ **HTTP-Only Cookies**: XSS protection
- ✅ **Password Hashing**: bcrypt with salt
- ✅ **Session Management**: Login status tracking

### **API Security**
- ✅ **Rate Limiting**: DDoS protection
- ✅ **CORS Configuration**: Cross-origin control
- ✅ **Input Validation**: SQL injection prevention
- ✅ **Error Handling**: Information disclosure prevention

### **Database Security**
- ✅ **Prepared Statements**: SQL injection protection
- ✅ **Connection Encryption**: TLS/SSL support
- ✅ **Access Control**: User-based data isolation
- ✅ **Audit Trail**: Soft delete with timestamps

## 🎯 Performance Features

### **Database Optimization**
- **Connection Pooling**: Max 100 connections, 10 idle connections
- **Strategic Indexing**: Auto-generated indexes for query optimization
- **Preload Relationships**: Efficient data loading with GORM preload
- **Query Optimization**: Composite indexes for common queries

### **API Performance**
- **Rate Limiting**: IP-based rate limiting with memory cleanup
- **Optimized Queries**: Strategic use of WHERE clauses and indexes
- **Efficient Models**: Well-structured entities with proper relationships
- **Memory Management**: Automatic visitor cleanup in rate limiter

## 📧 Contact

- **Email**: farhanlhsn@gmail.com
- **GitHub**: [@farhanlhsn](https://github.com/farhanlhsn)

---


<div align="center">

**⭐ Star this project if you find it helpful!**

Made with ❤️ by [Muhammad Farhan Al Hasan](https://github.com/farhanlhsn)

</div> 