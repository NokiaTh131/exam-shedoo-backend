package server

import (
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

	// routes enrollments
	enroll := s.App.Group("/enrollments")
	enroll.Post("/upload", s.EnrollmentHandler.UploadEnrollments)
	enroll.Get("/course", s.EnrollmentHandler.GetByCourseSections)
	enroll.Get("/:studentCode", s.EnrollmentHandler.GetByStudentCode)
	enroll.Delete("/:id", s.EnrollmentHandler.DeleteByID)

	// routes courseexam
	exam := s.App.Group("/course_exams")
	exam.Post("/examdate", s.CourseExamHandler.CreateExam)
	exam.Get("/course", s.CourseExamHandler.GetByCourseSections)
	exam.Put("/:id", s.CourseExamHandler.UpdateExam)

	// routes courses
	course := s.App.Group("/courses")
	course.Get("/lecturer", s.CourseHandler.GetCoursesByLecturer)
	course.Get("/code-sec", s.CourseHandler.GetCourseByCodeSec)

	// routes scrape jobs
	scrape := s.App.Group("/scrape")
	scrape.Post("/course/start", s.ScrapeJobHandler.CreateScrapeJob)
	scrape.Get("/course/status/:id", s.ScrapeJobHandler.GetScrapeJobByID)
	scrape.Post("/exams/start/:term", s.ScrapeJobHandler.CreateExamScrapeJob)
	scrape.Get("/exams/status/:id", s.ScrapeJobHandler.GetExamScrapeJobByID)
}
