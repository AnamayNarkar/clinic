# Clinic Management System

## Overview

The Clinic Management System is a comprehensive healthcare platform designed to facilitate efficient management of medical appointments, patient records, and doctor scheduling. This system enables seamless communication between patients, doctors, and administrative staff, streamlining the entire healthcare process from appointment requests to post-appointment follow-ups.

## Technologies Used

### Backend

- **Go (Golang)** - Main programming language
- **Gin Web Framework** - High-performance HTTP web framework
- **PostgreSQL** - Relational database for persistent storage
- **SQLC** - SQL compiler generating type-safe Go code from SQL
- **Redis** - In-memory data structure store used for session management
- **Flyway** - Database migration tool
- **UUID** - For generating unique identifiers
- **bcrypt** - For secure password hashing

### Security

- **JSON Web Tokens (JWT)** - For secure authentication
- **ABAC (Attribute-Based Access Control)** - For fine-grained authorization
- **bcrypt** with salt - For secure password storage

### Development Tools

- **Go Modules** - For dependency management
- **Docker** - For containerization and deployment
- **Git** - For version control

## System Architecture

The system follows a clean, layered architecture design:

1. **Presentation Layer**: HTTP handlers using the Gin framework
2. **Controller Layer**: Business logic handlers
3. **Data Access Layer**: SQLC-generated code for database operations
4. **Database Layer**: PostgreSQL with proper schema design

### Key Components

#### 1. User Management

The system supports four user roles:

- **Patients**: Can request appointments, view their medical history
- **Doctors**: Can manage appointments, update medical records
- **Receptionists**: Can coordinate appointments, manage patient records
- **Administrators**: Have full access to system management

#### 2. Authentication & Authorization

- **Session-based authentication** using Redis for state management
- **ABAC (Attribute-Based Access Control)** for granular permission control

#### 3. Appointment Management

- **Applications**: Initial requests for appointments
- **Appointments**: Scheduled and confirmed consultations
- **Results**: Medical records from consultations

## Database Schema

### Key Tables

- `admins`: System administrators
- `doctors`: Medical professionals
- `patients`: Healthcare clients
- `receptionists`: Front desk staff
- `applications`: Appointment requests
- `appointments`: Scheduled medical consultations

## Security Implementation - ABAC

The system implements Attribute-Based Access Control (ABAC), a flexible security model that evaluates permissions based on attributes of:

1. **Users**: Role, department, specialty
2. **Resources**: Type of data (appointments, patient records)
3. **Actions**: Create, read, update, delete
4. **Environmental context**: Time, location, etc.

### Permission Structure

The ABAC system is implemented in `securityManager.go`, which defines specific permissions for each role on various resources:

```
Role → Resource → Actions (CRUD)
```

For example:
- Doctors can read and update their own appointments but not create them
- Receptionists can create appointments but not update medical results
- Patients can view their own records but can't modify doctor's notes

### Implementation Details

The security manager intercepts API requests through middleware, evaluating permissions before allowing access to resources. This implementation provides:

- **Fine-grained control**: Specific permissions for each role-resource combination
- **Flexibility**: Easy to add new roles or resources
- **Auditability**: Clear permissions structure for compliance and review

## API Endpoints

### Authentication

- `POST /api/auth/login/:role` - Login for different user roles
- `POST /api/auth/logout` - Logout and invalidate session

### Patient Management

- `POST /api/patient` - Create patient record
- `GET /api/patient/:id` - Get patient details

### Doctor Management

- `POST /api/doctor` - Create doctor record
- `GET /api/doctors/appointments` - Get all appointments for a doctor

### Application Management

- `POST /api/application` - Create appointment application
- `GET /api/application` - List applications
- `PUT /api/application` - Update application status

### Appointment Management

- `GET /api/appointment/upcoming` - Get upcoming appointments
- `PUT /api/appointment/doctor/:id` - Update appointment by doctor

## Usage Flow

1. **Patient Registration**: Patients register in the system
2. **Appointment Request**: Patients submit appointment applications
3. **Receptionist Review**: Receptionists review applications and create appointments
4. **Doctor Consultation**: Doctors access appointments and updates medical results