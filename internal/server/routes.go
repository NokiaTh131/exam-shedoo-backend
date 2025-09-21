package server

import (
	"github.com/gofiber/fiber/v2"
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

	s.App.Get("/", s.HelloWorldHandler)
	// routes enrollments
	enroll := s.App.Group("/enrollments")
	enroll.Post("/upload", s.EnrollmentHandler.UploadEnrollments)
	enroll.Get("/course", s.EnrollmentHandler.GetByCourseSections)
	enroll.Get("/:studentCode", s.EnrollmentHandler.GetByStudentCode)
	enroll.Delete("/:id", s.EnrollmentHandler.DeleteByID)

	course := s.App.Group("/courses")
	course.Get("/", s.CourseHandler.GetCoursesByLecturer)
}

func (s *FiberServer) HelloWorldHandler(c *fiber.Ctx) error {
	resp := fiber.Map{
		"message": "Hello World",
	}

	return c.JSON(resp)
}
