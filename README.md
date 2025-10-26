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
  "cmuitaccount": "john.doe_acc",
  "cmuitaccount_name": "john.doe",
  "firstname_EN": "John",
  "firstname_TH": "จอห์น",
  "itaccounttype_EN": "Student",
  "itaccounttype_TH": "นักศึกษา",
  "itaccounttype_id": "StdAcc",
  "lastname_EN": "Doe",
  "lastname_TH": "โด",
  "organization_name_EN": "Faculty of Engineering",
  "organization_name_TH": "คณะวิศวกรรมศาสตร์",
  "role": "student",
  "student_id": "650610XXX"
}
```

### GET /auth/entraidurl
**Purpose**: Get the URL for CMU EntraID sign-in page
**Authentication**: None required
**Response**: Returns JWT claims with user information
```json
{
  "ok": true,
  "url": "https://login.microsoftonline.com/..."
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
**Response**:
```json
{
  "ID": 1,
  "StartCode": "261100",
  "EndCode": "261999",
  "Workers": 4,
  "Status": "running", // (or pending, completed, failed)
  "CreatedAt": "2024-01-01T12:00:00Z",
  "UpdatedAt": "2024-01-01T12:01:00Z"
}
```

#### POST /admin/scrape/exams/start/:term
**Purpose**: Start exam data scraping for a specific term  
**Parameters**: `term` - Academic term (e.g., "251")  
**Response**: `{"job_id": 2, "status": "pending"}`

#### GET /admin/scrape/exams/status/:id
**Purpose**: Check status of an exam scraping job  
**Parameters**: `id` - Job ID  
**Response**: 
```json
{
  "ID": 2,
  "Term": "251",
  "Status": "running", // (or pending, completed, failed)
  "CreatedAt": "2024-01-01T12:00:00Z",
  "UpdatedAt": "2024-01-01T12:01:00Z"
}
```

### Admin Role Management

#### GET /admin/
**Purpose**: List all admin accounts 
**Response**: Array of admin objects
```json
[
  {
    "ID": 1,
    "Account": "admin.user",
    "CreatedAt": "2024-01-01T10:00:00Z"
  }
]
```

#### POST /admin/
**Purpose**: Add a new admin by account name
**Request Body**:
```json
{
  "account": "new.admin"
}
```

#### DELETE /admin/:account
**Purpose**: Remove an admin by account name
**Request Body**:
```json
{"ok": true}
```
DELETE /admin/data/all
**Purpose**: Delete all main data (Enrollments, Courses, Exams, Jobs) from the database
**Response**: 
```json
{"ok": true}
```

### Enrollment Management

#### POST /admin/enrollments/upload
**Purpose**: Upload and import enrollment data from Excel/CSV file  
**Content-Type**: `multipart/form-data`  
**Form Data**: 
- `file` - Excel or CSV file containing enrollment data
- `exam_type` - "MIDTERM" or "FINAL"
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
    "lec_section": "001",
    "course_name": "Data Structures",
    "lab_section": "000",
    "lec_credit": 2,
    "lab_credit": 1,
    "instructors": [
      { "name": "Dr. Smith" },
      { "name": "Dr. John" }
    ],
    "room": "E12-101",
    "days": "MWF",
    "start_time": "09:00",
    "semester": "1",
    "year": "2567",
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
    "exam_id": 10,
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
  "CourseCode": "261200",
  "LecSection": "1",
  "LabSection": "1",
  "CourseID": 1,
  "MidtermExamDate": "25  MAR",
  "MidtermExamStartTime": "0900",
  "MidtermExamEndTime": "1100",
  "FinalExamDate": "25  OCT",
  "FinalExamStartTime": "1300",
  "FinalExamEndTime": "1500"
}
```
**Response**: Created exam object

#### PUT /professors/course_exams/:id
**Purpose**: Update exam schedule for a specific exam  
**Parameters**: `id` - Exam ID  
**Request Body** (all fields optional):
```json
{
  "midtermExamDate": "25  MAR",
  "finalExamDate": "25  OCT",
  "midtermExamStartTime": "1000",
  "finalExamStartTime": "1400",
  "midtermExamEndTime": "1200",
  "finalExamEndTime": "1600"
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
    "exam_id": 10,
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
          "students": [
            { "student_code": "650610111" },
            { "student_code": "650610112" }
          ]
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
          "students": [
            { "student_code": "650610111" },
            { "student_code": "650610112" }
          ]
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
          "students": [
            { "student_code": "650610111" },
            { "student_code": "650610112" }
          ]
        },
        {
          "course_id": 257,
          "course_code": "291494",
          "course_name": "SELECT TOPIC",
          "lec_section": "003",
          "lab_section": "000",
          "students": [
            { "student_code": "650610111" },
            { "student_code": "650610112" }
          ]
        }
      ]
    }
  ]
}
```


## Error Responses

All endpoints may return error responses in this format:
```json
{
  "error": "Error description here"
}
```

## Installation

### Docker

```bash
docker build -t course-scraper web-scraper/course_scraper.py
```

```bash
docker build -t exam-scraper web-scraper/exam_scraper.py
```

```bash
docker-compose up -d
```

Then **run** the application
```bash
make run
```

### Manual

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
