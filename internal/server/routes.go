package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/session"
)

func (s *FiberServer) RegisterFiberRoutes() {
	// Create a session store instance
	store := session.New()

	// Apply CSRF middleware
	s.App.Use(csrf.New(csrf.Config{
		CookieName:        "__Host-csrf_",
		CookieSameSite:    "Lax",
		CookieSecure:      true,
		CookieSessionOnly: true,
		CookieHTTPOnly:    true,
		Session:           store,
		SessionKey:        "fiber.csrf.token",
		HandlerContextKey: "fiber.csrf.handler",
	}))

	// Apply CORS middleware
	s.App.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS,PATCH",
		AllowHeaders:     "Accept,Authorization,Content-Type",
		AllowCredentials: false, // credentials require explicit origins
		MaxAge:           300,
	}))

	s.App.Get("/", s.HelloWorldHandler)

	s.App.Get("/health", s.healthHandler)

	api := s.App.Group("/api")
	api.Get("/cards", s.listCardsHandler)
	api.Post("/orders", s.createOrderHandler)
}

func (s *FiberServer) HelloWorldHandler(c *fiber.Ctx) error {
	resp := fiber.Map{
		"message": "Hello World",
	}

	return c.JSON(resp)
}

func (s *FiberServer) healthHandler(c *fiber.Ctx) error {
	return c.JSON(s.db.Health())
}

func (s *FiberServer) listCardsHandler(c *fiber.Ctx) error {
	cards, err := s.db.GetCardsFromDB()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"items": cards})
}

func (s *FiberServer) createOrderHandler(c *fiber.Ctx) error {
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "order accepted"})
}
