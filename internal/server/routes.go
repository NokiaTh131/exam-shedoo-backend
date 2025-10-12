package server

import (
	"shedoo-backend/internal/middlewares"

	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func (s *FiberServer) RegisterFiberRoutes() {
	// Apply CORS middleware
	s.App.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5173",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS,PATCH",
		AllowHeaders:     "Accept,Authorization,Content-Type",
		AllowCredentials: true,
		MaxAge:           300,
	}))
	// Apply logger middleware
	s.App.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${url}\n",
	}))

	// === Public routes ===
	auth := s.App.Group("/auth")
	auth.Post("/signin", s.AuthHandler.SignIn)

	s.App.Use(s.AuthMiddleware.AuthRequired())

	auth.Get("/profile", s.UserHandler.GetProfile)
	auth.Post("/signout", s.AuthHandler.SignOut)

	// === Admin routes ===
	admin := s.App.Group("/admin", s.AuthMiddleware.RequireRoles("admin"))

	admin.Get("/", s.AdminHandler.ListAdmins)
	admin.Post("/", s.AdminHandler.AddAdmin)
	admin.Delete("/data/all", s.AdminHandler.DeleteAllData)
	admin.Delete("/:account", s.AdminHandler.RemoveAdmin)
	admin.Post("/exampdf", s.CourseExamHandler.UploadPDF)

	scrape := admin.Group("/scrape")
	scrape.Post("/course/start", s.ScrapeJobHandler.CreateScrapeJob)
	scrape.Get("/course/status/:id", s.ScrapeJobHandler.GetScrapeJobByID)
	scrape.Post("/exams/start/:term", s.ScrapeJobHandler.CreateExamScrapeJob)
	scrape.Get("/exams/status/:id", s.ScrapeJobHandler.GetExamScrapeJobByID)

	enroll := admin.Group("/enrollments")
	enroll.Post("/upload", s.EnrollmentHandler.UploadEnrollments)

	// === Student routes ===
	student := s.App.Group("/students", s.AuthMiddleware.RequireRoles("student", "admin"))
	student.Get("/enrollments/:studentCode", middlewares.StudentOwnsResource(), s.EnrollmentHandler.GetEnrollmentsByStudent)
	student.Get("/exams/:studentCode", middlewares.StudentOwnsResource(), s.CourseExamHandler.GetExams)

	// === Professor routes ===
	professor := s.App.Group("/professors", s.AuthMiddleware.RequireRoles("professor", "admin"))

	exam := professor.Group("/course_exams")
	exam.Post("/examdate", s.CourseExamHandler.CreateExam)
	exam.Put("/:id", s.CourseExamHandler.UpdateExam)
	exam.Get("/report/:courseId", s.CourseExamHandler.GetExamReport)

	course := professor.Group("/courses")
	course.Get("/", middlewares.ProfessorOwnsResource(), s.CourseHandler.GetCoursesByLecturer)
	course.Get("/enrolled_students/:course_id", s.CourseHandler.GetEnrolledStudents)
}
