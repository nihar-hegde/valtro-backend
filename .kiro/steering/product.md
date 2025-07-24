# Valtro Backend API

Valtro is a backend API service that manages users, organizations, and projects. The system integrates with Clerk for authentication and provides RESTful endpoints for managing these core entities.

## Core Features
- User management with Clerk authentication integration
- Organization management
- Project management
- Health monitoring endpoints

## Key Characteristics
- RESTful API design with versioned endpoints (`/api/v1/`)
- PostgreSQL database with GORM ORM
- Soft deletion support for data integrity
- UUID-based primary keys for better scalability
- Environment-based configuration