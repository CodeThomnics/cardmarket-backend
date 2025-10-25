package server

import (
	"cardmarket_backend/internal/database"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/session"
)

func (s *FiberServer) RegisterFiberRoutes() {
	// Create a session store instance
	store := session.New()

	s.App.Use(logger.New())

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
	api.Post("/cards", s.createCardHandler)
	api.Get("/cards/:id", s.getCardByIDHandler)
	api.Put("/cards/:id", s.updateCardHandler)
	api.Delete("/cards/:id", s.deleteCardHandler)

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
	cards, err := s.db.ListCards()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch cards",
		})
	}
	return c.JSON(fiber.Map{"cards": cards})
}

func (s *FiberServer) createCardHandler(c *fiber.Ctx) error {
	var card database.CardRequest
	if err := c.BodyParser(&card); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}
	if err := s.db.CreateCard(card); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create card",
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "card created"})
}

func (s *FiberServer) getCardByIDHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	cardID, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid card ID",
		})
	}
	card, err := s.db.GetCardByID(cardID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Card not found",
		})
	}
	return c.JSON(fiber.Map{"card": card})
}

func (s *FiberServer) updateCardHandler(c *fiber.Ctx) error {
	var card database.CardRequest
	id := c.Params("id")
	cardID, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid card ID",
		})
	}

	if err := c.BodyParser(&card); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{})
	}

	s.db.UpdateCard(cardID, card)

	return c.JSON(fiber.Map{"message": "card updated"})
}

func (s *FiberServer) deleteCardHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	cardID, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid card ID",
		})
	}

	err = s.db.DeleteCard(cardID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete card",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "card deleted"})
}

func (s *FiberServer) createOrderHandler(c *fiber.Ctx) error {
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "order accepted"})
}
