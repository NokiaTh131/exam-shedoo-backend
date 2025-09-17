package server

import (
	"shedoo-backend/internal/app/enrollment"
	"shedoo-backend/internal/config"
	"shedoo-backend/internal/handlers"
	"shedoo-backend/internal/repositories"

	"github.com/gofiber/fiber/v2"
)

type FiberServer struct {
	*fiber.App
	EnrollmentHandler *handlers.EnrollmentHandler
}

func New() *FiberServer {
	dbService := config.New()
	enrollmentRepo := repositories.NewEnrollmentRepository(dbService.DB)
	enrollmentService := enrollment.NewEnrollmentService(enrollmentRepo)
	enrollmentHandler := handlers.NewEnrollmentHandler(enrollmentService)

	server := &FiberServer{
		App: fiber.New(fiber.Config{
			ServerHeader: "shedoo-backend",
			AppName:      "shedoo-backend",
		}),
		EnrollmentHandler: enrollmentHandler,
	}

	return server
}
