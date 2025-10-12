package server_test

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"shedoo-backend/internal/app/course"
	"shedoo-backend/internal/app/courseexam"
	"shedoo-backend/internal/app/enrollment"
	scrapejobs "shedoo-backend/internal/app/scrape_jobs"
	"shedoo-backend/internal/handlers"
	"shedoo-backend/internal/models"
	"shedoo-backend/internal/repositories"

	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	testApp *fiber.App
	testDB  *gorm.DB
)

func TestMain(m *testing.M) {
	var err error
	testDB, err = gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to open test DB: %v", err)
	}

	err = testDB.AutoMigrate(
		&models.Enrollment{},
		&models.Course{},
		&models.CourseExam{},
		&models.Admin{},
		&models.ScrapeCourseJob{},
		&models.ScrapeExamJob{},
	)
	if err != nil {
		log.Fatalf("failed to migrate test DB: %v", err)
	}

	seedData(testDB)

	testApp = newFeatureTestServerWithDB(testDB)

	os.Exit(m.Run())
}

func seedData(db *gorm.DB) {
	// Add Admin
	db.Create(&models.Admin{
		Account:   "admin1",
		CreatedAt: time.Now(),
	})

	// Add Course
	db.Create(&models.Course{
		CourseCode: "CS101",
		Title:      "Intro to CS",
		LecSection: ptr("001"),
		LabSection: ptr("002"),
		Lecturers:  []string{"prof1"},
	})

	// Add Enrollment
	db.Create(&models.Enrollment{
		StudentCode: "12345",
		CourseID:    1,
		CourseCode:  "CS101",
		LecSection:  "001",
		LabSection:  "002",
		Semester:    "S1",
		Year:        "2025",
	})

	// Add ScrapeCourseJob
	db.Create(&models.ScrapeCourseJob{
		StartCode: "100",
		EndCode:   "200",
		Workers:   4,
		Status:    "pending",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	// Add ScrapeExamJob
	db.Create(&models.ScrapeExamJob{
		Term:      "2025S",
		Status:    "pending",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
}

func ptr[T any](v T) *T { return &v }

func newFeatureTestServerWithDB(db *gorm.DB) *fiber.App {
	enrollmentHandler := handlers.NewEnrollmentHandler(
		enrollment.NewEnrollmentService(
			repositories.NewEnrollmentRepository(db),
		),
	)
	courseHandler := handlers.NewCourseHandler(
		course.NewCourseService(
			repositories.NewCourseRepository(db),
		),
	)
	courseExamHandler := handlers.NewCourseExamHandler(
		courseexam.NewCourseExamService(
			repositories.NewCourseExamRepository(db),
		),
	)
	scrapeJobHandler := handlers.NewScrapeJobHandler(
		scrapejobs.NewScrapeJobService(
			repositories.NewScrapeJobRepository(db),
		),
	)
	adminHandler := handlers.NewAdminHandler(nil)

	app := fiber.New()

	// === Admin routes ===
	admin := app.Group("/admin")
	admin.Get("/", adminHandler.ListAdmins)
	admin.Post("/", adminHandler.AddAdmin)
	admin.Delete("/data/all", adminHandler.DeleteAllData)
	admin.Delete("/:account", adminHandler.RemoveAdmin)
	admin.Post("/exampdf", courseExamHandler.UploadPDF)

	scrape := admin.Group("/scrape")
	scrape.Post("/course/start", scrapeJobHandler.CreateScrapeJob)
	scrape.Get("/course/status/:id", scrapeJobHandler.GetScrapeJobByID)
	scrape.Post("/exams/start/:term", scrapeJobHandler.CreateExamScrapeJob)
	scrape.Get("/exams/status/:id", scrapeJobHandler.GetExamScrapeJobByID)

	enroll := admin.Group("/enrollments")
	enroll.Post("/upload", enrollmentHandler.UploadEnrollments)

	// === Student routes ===
	student := app.Group("/students")
	student.Get("/enrollments/:studentCode", enrollmentHandler.GetEnrollmentsByStudent)
	student.Get("/exams/:studentCode", courseExamHandler.GetExams)

	// === Professor routes ===
	professor := app.Group("/professors")
	exam := professor.Group("/course_exams")
	exam.Post("/examdate", courseExamHandler.CreateExam)
	exam.Put("/:id", courseExamHandler.UpdateExam)
	exam.Get("/report/:courseId", courseExamHandler.GetExamReport)

	course := professor.Group("/courses")
	course.Get("/", courseHandler.GetCoursesByLecturer)
	course.Get("/enrolled_students/:course_id", courseHandler.GetEnrolledStudents)

	return app
}

func TestStudentEnrollments(t *testing.T) {
	// student exists in DB
	req := httptest.NewRequest("GET", "/students/enrollments/12345", nil)
	resp, _ := testApp.Test(req, -1)
	if resp.StatusCode != 200 {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
	// student does not exist in DB
	req = httptest.NewRequest("GET", "/students/enrollments/54321", nil)
	resp, _ = testApp.Test(req, -1)
	if resp.StatusCode != 404 {
		t.Errorf("expected 404, got %d", resp.StatusCode)
	}
}

func TestCourseRoutes(t *testing.T) {
	// student enrolled in course
	req := httptest.NewRequest("GET", "/professors/courses/enrolled_students/1", nil)
	resp, _ := testApp.Test(req, -1)
	if resp.StatusCode != 200 {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}

	// student not enrolled in course
	req = httptest.NewRequest("GET", "/professors/courses/enrolled_students/100", nil)
	resp, _ = testApp.Test(req, -1)
	if resp.StatusCode != 404 {
		t.Errorf("expected 404, got %d", resp.StatusCode)
	}
}

func TestCourseExamRoutes(t *testing.T) {
	correct_body := map[string]any{
		"courseCode":           "CS101",
		"lecSection":           "001",
		"labSection":           "002",
		"courseID":             1,
		"midtermExamDate":      "AUG 27",
		"midtermExamStartTime": "0900",
		"midtermExamEndTime":   "1200",
	}
	incorrect_body := map[string]any{
		"courseID":             1,
		"midtermExamDate":      "AUG 27",
		"midtermExamStartTime": "0900",
		"midtermExamEndTime":   "1200",
	}

	correct_buf, _ := json.Marshal(correct_body)
	incorrect_buf, _ := json.Marshal(incorrect_body)

	// correct request
	req := httptest.NewRequest("POST", "/professors/course_exams/examdate", bytes.NewBuffer(correct_buf))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := testApp.Test(req, -1)
	if resp.StatusCode != 201 {
		t.Errorf("expected 201, got %d", resp.StatusCode)
	}

	// incorrect request
	req = httptest.NewRequest("POST", "/professors/course_exams/examdate", bytes.NewBuffer(incorrect_buf))
	req.Header.Set("Content-Type", "application/json")
	resp, _ = testApp.Test(req, -1)
	if resp.StatusCode != 400 {
		t.Errorf("expected 400, got %d", resp.StatusCode)
	}

	// correct request (courseID exists)
	req = httptest.NewRequest("PUT", "/professors/course_exams/1", bytes.NewBuffer(correct_buf))
	req.Header.Set("Content-Type", "application/json")
	resp, _ = testApp.Test(req, -1)
	if resp.StatusCode != 200 {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}

	// incorrect request (courseID does not exist)
	req = httptest.NewRequest("PUT", "/professors/course_exams/99", bytes.NewBuffer(correct_buf))
	req.Header.Set("Content-Type", "application/json")
	resp, _ = testApp.Test(req, -1)
	if resp.StatusCode != 500 {
		t.Errorf("expected 500, got %d", resp.StatusCode)
	}
}

func TestScrapeRoutes(t *testing.T) {
	// correct request (start code is less than end code)
	correct_reqBody := map[string]any{
		"start":   "100",
		"end":     "200",
		"workers": 4,
	}
	correct_buf, _ := json.Marshal(correct_reqBody)

	req := httptest.NewRequest("POST", "/admin/scrape/course/start", bytes.NewBuffer(correct_buf))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := testApp.Test(req, -1)
	if resp.StatusCode != 200 {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}

	// incorrect request (start code is greater than end code)
	incorrect_reqBody := map[string]any{
		"start":   "200",
		"end":     "100",
		"workers": 4,
	}
	incorrect_buf, _ := json.Marshal(incorrect_reqBody)

	req = httptest.NewRequest("POST", "/admin/scrape/course/start", bytes.NewBuffer(incorrect_buf))
	req.Header.Set("Content-Type", "application/json")
	resp, _ = testApp.Test(req, -1)
	if resp.StatusCode != 400 {
		t.Errorf("expected 400, got %d", resp.StatusCode)
	}
}
