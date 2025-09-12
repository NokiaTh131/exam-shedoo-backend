package server

import (
	"shedoo-backend/internal/config"

	"github.com/gofiber/fiber/v2"
)

type FiberServer struct {
	*fiber.App
	db *config.Service
}

func New() *FiberServer {
	server := &FiberServer{
		App: fiber.New(fiber.Config{
			ServerHeader: "shedoo-backend",
			AppName:      "shedoo-backend",
		}),
		db: config.New(),
	}

	return server
}
