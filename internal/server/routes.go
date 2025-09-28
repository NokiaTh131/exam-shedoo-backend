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

	// === Public routes ===
	auth := s.App.Group("/auth")
	auth.Post("/signin", s.AuthHandler.SignIn)

	s.App.Use(middlewares.AuthRequired(s.RoleService))

	// profile all authenticated users
	auth.Get("/profile", s.UserHandler.GetProfile)

	// === Admin routes ===
	admin := s.App.Group("/admin", middlewares.RequireRoles("admin"))

	scrape := admin.Group("/scrape")
	scrape.Post("/course/start", s.ScrapeJobHandler.CreateScrapeJob)
	scrape.Get("/course/status/:id", s.ScrapeJobHandler.GetScrapeJobByID)
	scrape.Post("/exams/start/:term", s.ScrapeJobHandler.CreateExamScrapeJob)
	scrape.Get("/exams/status/:id", s.ScrapeJobHandler.GetExamScrapeJobByID)

	enroll := admin.Group("/enrollments")
	enroll.Post("/upload", s.EnrollmentHandler.UploadEnrollments)
	enroll.Delete("/:id", s.EnrollmentHandler.DeleteByID)

	// === Student routes ===
	student := s.App.Group("/students", middlewares.RequireRoles("student"))
	student.Get("/enrollments/:studentCode", middlewares.StudentOwnsResource(), s.EnrollmentHandler.GetEnrollmentsByStudent)
	student.Get("/exams/:studentCode", middlewares.StudentOwnsResource(), s.CourseExamHandler.GetExams)

	// === Professor routes ===
	professor := s.App.Group("/professors", middlewares.RequireRoles("professor", "admin"))

	exam := professor.Group("/course_exams")
	exam.Post("/examdate", s.CourseExamHandler.CreateExam)
	exam.Put("/:id", s.CourseExamHandler.UpdateExam)
	exam.Get("/report/:courseId", s.CourseExamHandler.GetExamReport)

	course := professor.Group("/courses")
	course.Get("/", middlewares.ProfessorOwnsResource(), s.CourseHandler.GetCoursesByLecturer)
	course.Get("/enrolled_students/:course_id", s.CourseHandler.GetEnrolledStudents)
}
