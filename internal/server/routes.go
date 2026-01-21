package server

import (
	"cardmarket_backend/internal/database"
	"database/sql"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/golang-jwt/jwt/v5"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func parseParam(c *fiber.Ctx, param string) (int, error) {
	idStr := c.Params(param)
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func createToken(username string) (string, error) {
	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username":   username,
		"expiration": time.Now().Add(time.Hour * 24).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *FiberServer) RegisterFiberRoutes() {
	// Create a session store instance
	store := session.New()

	s.App.Use(logger.New())

	// Apply CSRF middleware
	s.App.Use(csrf.New(csrf.Config{
		CookieName:        "csrf_",
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

	api.Get("/products", s.listProductsHandler)
	api.Post("/products", s.createProductHandler)
	api.Get("/products/:id", s.getProductByIDHandler)
	api.Put("/products/:id", s.updateProductHandler)
	api.Delete("/products/:id", s.deleteProductHandler)

	api.Get("/orders", s.listOrdersHandler)
	api.Get("/orders/:id", s.getOrderByIDHandler)
	api.Post("/orders", s.createOrderHandler)
	api.Put("/orders/:id", s.updateOrderHandler)
	api.Delete("/orders/:id", s.deleteOrderHandler)

	api.Get("/users", s.listUsersHandler)
	api.Post("/users", s.createUserHandler)
	api.Get("/users/:id", s.getUserByIDHandler)
	api.Put("/users/:id", s.updateUserHandler)
	api.Delete("/users/:id", s.deleteUserHandler)

	api.Post("/login", s.loginHandler)
	api.Post("/logout", s.logoutHandler)
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
	if err == database.ErrNoCards {
		return c.Status(fiber.StatusNotFound).
			JSON(fiber.Map{"cards": []database.Card{}, "error": "no cards found"})
	}
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
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
			"error": err.Error(),
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "card created"})
}

func (s *FiberServer) getCardByIDHandler(c *fiber.Ctx) error {
	cardID, err := parseParam(c, "id")

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid card ID",
		})
	}
	card, err := s.db.GetCardByID(cardID)
	if err == database.ErrNoCardFound {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "card not found",
		})
	}
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{"card": card})
}

func (s *FiberServer) updateCardHandler(c *fiber.Ctx) error {
	var card database.CardRequest
	cardID, err := parseParam(c, "id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid card ID",
		})
	}

	if err := c.BodyParser(&card); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	err = s.db.UpdateCard(cardID, card)
	if err == sql.ErrNoRows {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "card could not be updated",
		})
	}
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{"message": "card updated"})
}

func (s *FiberServer) deleteCardHandler(c *fiber.Ctx) error {
	cardID, err := parseParam(c, "id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid card ID",
		})
	}

	err = s.db.DeleteCard(cardID)
	if err == sql.ErrNoRows {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "card could not be deleted",
		})
	}
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "card deleted"})
}

func (s *FiberServer) listProductsHandler(c *fiber.Ctx) error {
	products, err := s.db.ListProducts()
	if err == database.ErrNoProducts {
		return c.JSON(fiber.Map{"products": []database.Product{}, "error": "no products found"})
	}
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{"products": products})
}

func (s *FiberServer) createProductHandler(c *fiber.Ctx) error {
	var product database.ProductRequest
	if err := c.BodyParser(&product); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	err := s.db.CreateProduct(product)

	if err == sql.ErrNoRows {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "product could not be created",
		})
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "product created"})
}

func (s *FiberServer) updateProductHandler(c *fiber.Ctx) error {
	productID, err := parseParam(c, "id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid product ID",
		})
	}

	var product database.ProductRequest

	err = c.BodyParser(&product)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	err = s.db.UpdateProduct(productID, product)

	if err == sql.ErrNoRows {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "product could not be updated",
		})
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{"message": "product updated"})
}

func (s *FiberServer) deleteProductHandler(c *fiber.Ctx) error {
	productID, err := parseParam(c, "id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid product ID",
		})
	}

	err = s.db.DeleteProduct(productID)

	if err == sql.ErrNoRows {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "product could not be deleted",
		})
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "product deleted"})
}

func (s *FiberServer) getProductByIDHandler(c *fiber.Ctx) error {
	productID, err := parseParam(c, "id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid product ID",
		})
	}
	product, err := s.db.GetProductByID(productID)

	if err == database.ErrNoProductFound {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "product not found",
		})
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{"product": product})
}

func (s *FiberServer) listOrdersHandler(c *fiber.Ctx) error {
	orders, err := s.db.ListOrders()

	if err == database.ErrNoOrders {
		return c.JSON(fiber.Map{"orders": []database.Order{}, "error": "no orders found"})
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{"orders": orders})
}

func (s *FiberServer) getOrderByIDHandler(c *fiber.Ctx) error {
	orderID, err := parseParam(c, "id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid order ID",
		})
	}
	order, err := s.db.GetOrderByID(orderID)

	if err == database.ErrNoOrderFound {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "order not found",
		})
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{"order": order})
}

func (s *FiberServer) createOrderHandler(c *fiber.Ctx) error {
	var order database.OrderRequest
	if err := c.BodyParser(&order); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	err := s.db.CreateOrder(order)

	if err == sql.ErrNoRows {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "order could not be created",
		})
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "order accepted"})
}

func (s *FiberServer) updateOrderHandler(c *fiber.Ctx) error {
	orderID, err := parseParam(c, "id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid order ID",
		})
	}

	var order database.OrderRequest
	if err := c.BodyParser(&order); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	err = s.db.UpdateOrder(orderID, order)

	if err == sql.ErrNoRows {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "order could not be updated",
		})
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{"message": "order updated"})
}

func (s *FiberServer) deleteOrderHandler(c *fiber.Ctx) error {
	orderID, err := parseParam(c, "id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid order ID",
		})
	}

	err = s.db.DeleteOrder(orderID)

	if err == sql.ErrNoRows {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "order could not be deleted",
		})
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "order deleted"})
}

func (s *FiberServer) listUsersHandler(c *fiber.Ctx) error {
	users, err := s.db.ListUsers()
	if err == database.ErrNoUsers {
		return c.JSON(fiber.Map{"users": []database.User{}, "error": "no users found"})
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{"users": users})
}

func (s *FiberServer) getUserByIDHandler(c *fiber.Ctx) error {
	userID, err := parseParam(c, "id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}
	user, err := s.db.GetUserByID(userID)

	if err == database.ErrNoUserFound {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "user not found",
		})
	}
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{"user": user})
}

func (s *FiberServer) createUserHandler(c *fiber.Ctx) error {
	var user database.UserRequest
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}
	if err := s.db.CreateUser(user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "user created"})
}

func (s *FiberServer) updateUserHandler(c *fiber.Ctx) error {
	userID, err := parseParam(c, "id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	var user database.UserRequest

	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	err = s.db.UpdateUser(userID, user)

	if err == sql.ErrNoRows {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "user could not be updated",
		})
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{"message": "user updated"})
}

func (s *FiberServer) deleteUserHandler(c *fiber.Ctx) error {
	userID, err := parseParam(c, "id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	err = s.db.DeleteUser(userID)

	if err == sql.ErrNoRows {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "user could not be deleted",
		})
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "user deleted"})
}

func (s *FiberServer) loginHandler(c *fiber.Ctx) error {
	var request LoginRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	check, err := s.db.CheckLoginCredentials(request.Username, request.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	if !check {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid credentials"})
	}

	token, err := createToken(request.Username)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "login successful", "token": token})
}

func (s *FiberServer) logoutHandler(c *fiber.Ctx) error {
	// Placeholder for logout logic
	return c.JSON(fiber.Map{"message": "logout successful"})
}
