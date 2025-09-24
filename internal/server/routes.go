package server

import (
	"shedoo-backend/internal/middlewares"

	"github.com/gofiber/fiber/v2/middleware/cors"
)

func (s *FiberServer) RegisterFiberRoutes() {
	// Apply CORS middleware
	s.App.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS,PATCH",
		AllowHeaders:     "Accept,Authorization,Content-Type",
		AllowCredentials: false, // credentials require explicit origins
		MaxAge:           300,
	}))

	// routes auth
	auth := s.App.Group("/auth")
	auth.Post("/signin", s.AuthHandler.SignIn)
	auth.Get("/profile", middlewares.AuthRequired, s.UserHandler.GetProfile)

	// routes enrollments
	enroll := s.App.Group("/enrollments")
	enroll.Post("/upload", s.EnrollmentHandler.UploadEnrollments)     // Admin
	enroll.Get("/course", s.EnrollmentHandler.GetByCourseSections)    // Professor
	enroll.Get("/:studentCode", s.EnrollmentHandler.GetByStudentCode) // Professor, Student
	enroll.Delete("/:id", s.EnrollmentHandler.DeleteByID)             // Admin

	// routes courseexam
	exam := s.App.Group("/course_exams")
	exam.Post("/examdate", s.CourseExamHandler.CreateExam)       // Professor
	exam.Get("/course", s.CourseExamHandler.GetByCourseSections) // Professor, Student
	exam.Put("/:id", s.CourseExamHandler.UpdateExam)             // Professor

	// routes courses
	course := s.App.Group("/courses")
	course.Get("/lecturer", s.CourseHandler.GetCoursesByLecturer) // Professor
	course.Get("/code-sec", s.CourseHandler.GetCourseByCodeSec)   // Professor, Student

	// routes scrape jobs
	scrape := s.App.Group("/scrape")
	scrape.Post("/course/start", s.ScrapeJobHandler.CreateScrapeJob)          // Admin
	scrape.Get("/course/status/:id", s.ScrapeJobHandler.GetScrapeJobByID)     // Admin
	scrape.Post("/exams/start/:term", s.ScrapeJobHandler.CreateExamScrapeJob) // Admin
	scrape.Get("/exams/status/:id", s.ScrapeJobHandler.GetExamScrapeJobByID)  // Admin

	// student routes
	student := s.App.Group("/students")
	student.Get("/enrollments/:studentCode", s.StudentHandler.GetEnrollmentsByStudent)
	student.Get("/exams/:studentCode", s.StudentHandler.GetExams)
}
