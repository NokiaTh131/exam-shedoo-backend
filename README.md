# Project shedoo-backend

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

## MakeFile

Run build make command with tests
```bash
make all
```

Build the application
```bash
make build
```

Run the application
```bash
make run
```

Live reload the application:
```bash
make watch
```

Run the test suite:
```bash
make test
```

Clean up binary from the last build:
```bash
make clean
```

# Shedoo Backend API Routes

A comprehensive guide to all API routes in the Shedoo backend system.

## Base URL
All routes are prefixed with your server's base URL (e.g., `http://localhost:8080` or your domain).

## Authentication

### POST /auth/signin
**Purpose**: Sign in with CMU EntraID authorization code  
**Authentication**: None required  
**Request Body**:
```json
{
  "authorizationCode": "your_authorization_code_here"
}
```
**Response**: 
- Success: `{"ok": true}` + sets authentication cookie
- Error: `{"ok": false, "message": "error_message"}`

### GET /auth/profile
**Purpose**: Get current user's profile information  
**Authentication**: Required (cookie-based)  
**Response**: Returns JWT claims with user information
```json
{
  "cmuitaccount_name": "john.doe",
  "student_id": "123456789",
  "firstname_EN": "John",
  "lastname_EN": "Doe",
  // ... other user fields
}
```

## Admin Routes
*Routes for administrative functions*

### Scraping Jobs

#### POST /admin/scrape/course/start
**Purpose**: Start a new course scraping job  
**Request Body**:
```json
{
  "start": "261100",
  "end": "261999", 
  "workers": 4
}
```
**Response**: `{"job_id": 1, "status": "pending"}`

#### GET /admin/scrape/course/status/:id
**Purpose**: Check status of a course scraping job  
**Parameters**: `id` - Job ID  
**Response**: Job details with current status

#### POST /admin/scrape/exams/start/:term
**Purpose**: Start exam data scraping for a specific term  
**Parameters**: `term` - Academic term (e.g., "251")  
**Response**: `{"job_id": 2, "status": "pending"}`

#### GET /admin/scrape/exams/status/:id
**Purpose**: Check status of an exam scraping job  
**Parameters**: `id` - Job ID  
**Response**: Job details with current status

### Enrollment Management

#### POST /admin/enrollments/upload
**Purpose**: Upload and import enrollment data from Excel/CSV file  
**Content-Type**: `multipart/form-data`  
**Form Data**: 
- `file` - Excel or CSV file containing enrollment data
**Response**: `{"message": "Imported X records"}`

#### DELETE /admin/enrollments/:id
**Purpose**: Delete a specific enrollment record  
**Parameters**: `id` - Enrollment ID  
**Response**: `{"message": "deleted successfully"}`

## Student Routes
*Routes for student-specific data*

### GET /students/enrollments/:studentCode
**Purpose**: Get all courses a student is enrolled in  
**Parameters**: `studentCode` - Student's code/ID  
**Response**: Array of enrollment details
```json
[
  {
    "id": 1,
    "course_code": "261200",
    "course_name": "Data Structures",
    "lec_section": "001",
    "lab_section": "000", 
    "credit": 3,
    "instructors": ["Dr. Smith", "Dr. John"],
    "room": "E12-101",
    "days": "MWF",
    "start_time": "09:00",
    "end_time": "10:00"
  }
]
```

### GET /students/exams/:studentCode
**Purpose**: Get all exam schedules for a student  
**Parameters**: `studentCode` - Student's code/ID  
**Response**: Array of exam details
```json
[
  {
    "id": 1,
    "course_code": "261200",
    "course_name": "Data Structures",
    "lec_section": "1",
    "lab_section": "1",
    "midterm_date": "2024-03-15",
    "midterm_start_time": "09:00",
    "midterm_end_time": "11:00",
    "final_date": "2024-05-20",
    "final_start_time": "13:00", 
    "final_end_time": "15:00"
  }
]
```

## Professor Routes
*Routes for professor/lecturer functions*

### Course Management

#### GET /professors/courses?lecturer=lecturerName
**Purpose**: Get all courses taught by a specific lecturer  
**Query Parameters**: `lecturer` - Lecturer's name  
**Response**: Array of courses with exam information
```json
[
  {
    "course_id": 1,
    "course_code": "261200",
    "course_name": "Data Structures", 
    "lec_section": "1",
    "lab_section": "1",
    "midterm_date": "2024-03-15",
    "midterm_start_time": "09:00",
    "midterm_end_time": "11:00",
    "final_date": "2024-05-20",
    "final_start_time": "13:00",
    "final_end_time": "15:00"
  }
]
```

#### GET /professors/courses/enrolled_students/:course_id
**Purpose**: Get list of students enrolled in a specific course  
**Parameters**: `course_id` - Course ID  
**Response**: Array of enrolled students
```json
[
  {
    "enrollment_id": 1,
    "student_code": "650610123"
  },
  {
    "enrollment_id": 2, 
    "student_code": "650610124"
  }
]
```

### Exam Management

#### POST /professors/course_exams/examdate
**Purpose**: Create new exam schedule for a course  
**Request Body**:
```json
{
  "course_code": "261200",
  "lec_section": "1",
  "lab_section": "1",
  "midterm_exam_date": "25   MAR",
  "midterm_exam_start_time": "09:00",
  "midterm_exam_end_time": "11:00",
  "final_exam_date": "25   OCT", 
  "final_exam_start_time": "13:00",
  "final_exam_end_time": "15:00"
}
```
**Response**: Created exam object

#### PUT /professors/course_exams/:id
**Purpose**: Update exam schedule for a specific exam  
**Parameters**: `id` - Exam ID  
**Request Body** (all fields optional):
```json
{
  "midterm_exam_date": "25   MAR",
  "final_exam_date": "25   OCT", 
  "midterm_exam_start_time": "10:00",
  "final_exam_start_time": "14:00",
  "midterm_exam_end_time": "12:00",
  "final_exam_end_time": "16:00"
}
```
**Response**: `{"message": "updated examdate"}`

#### GET /professors/course_exams/midterm/:lecturerName
**Purpose**: Get midterm exam report for all courses taught by a lecturer  
**Parameters**: `lecturerName` - Lecturer's name  
**Response**: Array of midterm exam details with student counts
```json
[
  {
    "course_id": 1,
    "course_code": "261200",
    "course_name": "Data Structures",
    "lec_section": "1", 
    "lab_section": "1",
    "number_of_relevant_students": 45,
    "exam_date": "20   OCT",
    "start_time": "09:00",
    "end_time": "11:00"
  }
]
```

#### GET /professors/course_exams/report/:courseId
**Purpose**: Get all exam_datetimes of students who enrolled in a specific course  
**Parameters**: `courseId` - Course ID
**Response**: 
```json
{
    "midterm": [
        {
            "date": "AUG  31",
            "start_time": "0800",
            "end_time": "1100",
            "courses": [
                {
                    "course_id": 256,
                    "course_code": "291494",
                    "course_name": "SELECT TOPIC",
                    "lec_section": "006",
                    "lab_section": "000",
                    "student_count": 2
                }
            ]
        },
        {
            "date": "AUG  26",
            "start_time": "0800",
            "end_time": "1100",
            "courses": [
                {
                    "course_id": 257,
                    "course_code": "291494",
                    "course_name": "SELECT TOPIC",
                    "lec_section": "003",
                    "lab_section": "000",
                    "student_count": 2
                }
            ]
        },
        {
            "date": "AUG  25",
            "start_time": "1200",
            "end_time": "1500",
            "courses": [
                {
                    "course_id": 198,
                    "course_code": "001101",
                    "course_name": "Fundamental English 1",
                    "lec_section": "717",
                    "lab_section": "000",
                    "student_count": 1
                }
            ]
        },
        {
            "date": "AUG  25",
            "start_time": "0800",
            "end_time": "1100",
            "courses": [
                {
                    "course_id": 223,
                    "course_code": "001102",
                    "course_name": "Fundamental English 2",
                    "lec_section": "003",
                    "lab_section": "000",
                    "student_count": 1
                },
                {
                    "course_id": 226,
                    "course_code": "001102",
                    "course_name": "Fundamental English 2",
                    "lec_section": "008",
                    "lab_section": "000",
                    "student_count": 1
                }
            ]
        }
    ],
    "final": [
        {
            "date": "OCT  20",
            "start_time": "0800",
            "end_time": "1100",
            "courses": [
                {
                    "course_id": 256,
                    "course_code": "291494",
                    "course_name": "SELECT TOPIC",
                    "lec_section": "006",
                    "lab_section": "000",
                    "student_count": 2
                },
                {
                    "course_id": 257,
                    "course_code": "291494",
                    "course_name": "SELECT TOPIC",
                    "lec_section": "003",
                    "lab_section": "000",
                    "student_count": 2
                },
                {
                    "course_id": 223,
                    "course_code": "001102",
                    "course_name": "Fundamental English 2",
                    "lec_section": "003",
                    "lab_section": "000",
                    "student_count": 1
                },
                {
                    "course_id": 226,
                    "course_code": "001102",
                    "course_name": "Fundamental English 2",
                    "lec_section": "008",
                    "lab_section": "000",
                    "student_count": 1
                }
            ]
        },
        {
            "date": "OCT  20",
            "start_time": "1200",
            "end_time": "1500",
            "courses": [
                {
                    "course_id": 198,
                    "course_code": "001101",
                    "course_name": "Fundamental English 1",
                    "lec_section": "717",
                    "lab_section": "000",
                    "student_count": 1
                }
            ]
        }
    ]
}```

## Error Responses

All endpoints may return error responses in this format:
```json
{
  "error": "Error description here"
}
```

Common HTTP status codes:
- `200` - Success
- `201` - Created successfully  
- `400` - Bad request (invalid input)
- `401` - Unauthorized (authentication required)
- `404` - Not found
- `500` - Internal server error

## Authentication Notes

- Authentication is handled via HTTP-only cookies
- The `/auth/profile` route and some admin routes require authentication
- Student and professor routes currently don't seem to require authentication based on the route configuration
- Authentication cookies are set with domain and security configurations based on the `PROD` environment variable

## File Upload Notes

- The `/admin/enrollments/upload` endpoint accepts Excel (.xlsx) or CSV files
- Files are temporarily stored in `./tmp/` directory
- Supported file formats are automatically detected
- Excel files are processed using the `excelize` library
- CSV files are processed with Go's built-in CSV reader

## Database Integration

- Uses GORM as the ORM
- PostgreSQL as the database
- Supports automatic migrations for all models
- Uses JSON data types for storing arrays (like lecturer lists)
- Implements proper foreign key relationships between models
