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

	// admin routes
	admin := s.App.Group("/admin")
	scrape := admin.Group("/scrape")
	enroll := admin.Group("/enrollments")
	scrape.Post("/course/start", s.ScrapeJobHandler.CreateScrapeJob)
	scrape.Get("/course/status/:id", s.ScrapeJobHandler.GetScrapeJobByID)
	scrape.Post("/exams/start/:term", s.ScrapeJobHandler.CreateExamScrapeJob)
	scrape.Get("/exams/status/:id", s.ScrapeJobHandler.GetExamScrapeJobByID)
	enroll.Post("/upload", s.EnrollmentHandler.UploadEnrollments)
	enroll.Delete("/:id", s.EnrollmentHandler.DeleteByID)

	// student routes
	student := s.App.Group("/students")
	student.Get("/enrollments/:studentCode", s.EnrollmentHandler.GetEnrollmentsByStudent)
	student.Get("/exams/:studentCode", s.CourseExamHandler.GetExams)

	// professor routes
	professor := s.App.Group("/professors")
	exam := professor.Group("/course_exams")
	course := professor.Group("/courses")
	course.Get("/", s.CourseHandler.GetCoursesByLecturer)
	course.Get("/enrolled_students/:course_id", s.CourseHandler.GetEnrolledStudents)
	exam.Post("/examdate", s.CourseExamHandler.CreateExam)
	exam.Put("/:id", s.CourseExamHandler.UpdateExam)
	exam.Get("/midterm/:lecturerName", s.CourseExamHandler.GetMidtermExamReport)
}
