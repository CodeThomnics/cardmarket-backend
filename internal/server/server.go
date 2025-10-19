package server

import (
	"github.com/gofiber/fiber/v2"

	"cardmarket_backend/internal/database"
)

type FiberServer struct {
	*fiber.App

	db database.Service
}

func New() *FiberServer {
	server := &FiberServer{
		App: fiber.New(fiber.Config{
			ServerHeader: "cardmarket_backend",
			AppName:      "cardmarket_backend",
		}),

		db: database.New(),
	}

	return server
}
