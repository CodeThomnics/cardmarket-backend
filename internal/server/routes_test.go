package server

import (
	"cardmarket_backend/internal/database"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
)

type MockDBService struct {
	ListCardsFunc      func() ([]database.Card, error)
	GetCardByIDFunc    func() (database.Card, error)
	CreateCardFunc     func(card database.CardRequest) error
	UpdateCardFunc     func(cardID int, card database.CardRequest) error
	DeleteCardFunc     func(cardID int) error
	ListProductsFunc   func() ([]database.Product, error)
	GetProductByIDFunc func(productID int) (database.Product, error)
	CreateProductFunc  func(product database.ProductRequest) error
	UpdateProductFunc  func(productID int, product database.ProductRequest) error
	DeleteProductFunc  func(productID int) error
	ListOrdersFunc     func() ([]database.Order, error)
	GetOrderByIDFunc   func(orderID int) (database.Order, error)
	CreateOrderFunc    func(order database.OrderRequest) error
	UpdateOrderFunc    func(orderID int, order database.OrderRequest) error
	DeleteOrderFunc    func(orderID int) error
	ListUsersFunc      func() ([]database.User, error)
	GetUserByIDFunc    func(userID int) (database.User, error)
	CreateUserFunc     func(user database.UserRequest) error
	UpdateUserFunc     func(userID int, user database.UserRequest) error
	DeleteUserFunc     func(userID int) error
}

func (m *MockDBService) Close() error {
	return nil
}

func (m *MockDBService) Health() map[string]string {
	return map[string]string{"status": "up"}
}

func (m *MockDBService) ListCards() ([]database.Card, error) {
	if m.ListCardsFunc != nil {
		return m.ListCardsFunc()
	}
	return []database.Card{}, nil
}

func (m *MockDBService) GetCardByID(id int) (database.Card, error) {
	if m.GetCardByIDFunc != nil {
		return m.GetCardByIDFunc()
	}
	return database.Card{}, nil
}

func (m *MockDBService) CreateCard(card database.CardRequest) error {
	if m.CreateCardFunc != nil {
		return m.CreateCardFunc(card)
	}
	return nil
}

func (m *MockDBService) UpdateCard(cardID int, card database.CardRequest) error {
	if m.UpdateCardFunc != nil {
		return m.UpdateCardFunc(cardID, card)
	}
	return nil
}

func (m *MockDBService) DeleteCard(cardID int) error {
	if m.DeleteCardFunc != nil {
		return m.DeleteCardFunc(cardID)
	}
	return nil
}

func (m *MockDBService) ListProducts() ([]database.Product, error) {
	if m.ListProductsFunc != nil {
		return m.ListProductsFunc()
	}
	return []database.Product{}, nil
}

func (m *MockDBService) GetProductByID(productID int) (database.Product, error) {
	if m.GetProductByIDFunc != nil {
		return m.GetProductByIDFunc(productID)
	}
	return database.Product{}, nil
}

func (m *MockDBService) CreateProduct(product database.ProductRequest) error {
	if m.CreateProductFunc != nil {
		return m.CreateProductFunc(product)
	}
	return nil
}

func (m *MockDBService) UpdateProduct(productID int, product database.ProductRequest) error {
	if m.UpdateProductFunc != nil {
		return m.UpdateProductFunc(productID, product)
	}
	return nil
}

func (m *MockDBService) DeleteProduct(productID int) error {
	if m.DeleteProductFunc != nil {
		return m.DeleteProductFunc(productID)
	}
	return nil
}

func (m *MockDBService) ListOrders() ([]database.Order, error) {
	if m.ListOrdersFunc != nil {
		return m.ListOrdersFunc()
	}
	return []database.Order{}, nil
}

func (m *MockDBService) GetOrderByID(orderID int) (database.Order, error) {
	if m.GetOrderByIDFunc != nil {
		return m.GetOrderByIDFunc(orderID)
	}
	return database.Order{}, nil
}

func (m *MockDBService) CreateOrder(order database.OrderRequest) error {
	if m.CreateOrderFunc != nil {
		return m.CreateOrderFunc(order)
	}
	return nil
}

func (m *MockDBService) UpdateOrder(orderID int, order database.OrderRequest) error {
	if m.UpdateOrderFunc != nil {
		return m.UpdateOrderFunc(orderID, order)
	}
	return nil
}

func (m *MockDBService) DeleteOrder(orderID int) error {
	if m.DeleteOrderFunc != nil {
		return m.DeleteOrderFunc(orderID)
	}
	return nil
}

func (m *MockDBService) ListUsers() ([]database.User, error) {
	if m.ListUsersFunc != nil {
		return m.ListUsersFunc()
	}
	return []database.User{}, nil
}

func (m *MockDBService) GetUserByID(userID int) (database.User, error) {
	if m.GetUserByIDFunc != nil {
		return m.GetUserByIDFunc(userID)
	}
	return database.User{}, nil
}

func (m *MockDBService) CreateUser(user database.UserRequest) error {
	if m.CreateUserFunc != nil {
		return m.CreateUserFunc(user)
	}
	return nil
}

func (m *MockDBService) UpdateUser(userID int, user database.UserRequest) error {
	if m.UpdateUserFunc != nil {
		return m.UpdateUserFunc(userID, user)
	}
	return nil
}

func (m *MockDBService) DeleteUser(userID int) error {
	if m.DeleteUserFunc != nil {
		return m.DeleteUserFunc(userID)
	}
	return nil
}

func TestHandler(t *testing.T) {
	// Create a Fiber app for testing
	app := fiber.New()
	// Inject the Fiber app into the server
	s := &FiberServer{App: app}
	// Define a route in the Fiber app
	app.Get("/", s.HelloWorldHandler)
	// Create a test HTTP request
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatalf("error creating request. Err: %v", err)
	}
	// Perform the request
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("error making request to server. Err: %v", err)
	}
	// Your test assertions...
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status OK; got %v", resp.Status)
	}
	expected := "{\"message\":\"Hello World\"}"
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error reading response body. Err: %v", err)
	}
	if expected != string(body) {
		t.Errorf("expected response body to be %v; got %v", expected, string(body))
	}
}

func TestHealthHandler(t *testing.T) {
	app := fiber.New()
	s := &FiberServer{App: app, db: database.New()}
	app.Get("/health", s.healthHandler)

	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatalf("error creating request. Err: %v", err)
	}

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("error making request to server. Err: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status OK; got %v", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error reading response body. Err: %v", err)
	}

	expectedSubstring := "\"status\":\"up\""
	if strings.Contains(string(body), expectedSubstring) {
		t.Logf("health check passed")
	} else {
		t.Errorf("expected response body to contain %v; got %v", expectedSubstring, string(body))
	}
}

func TestListCardsHandler(t *testing.T) {

	cards := []database.Card{
		{ID: 1, Name: "Black Lotus", ImageURL: "https://example.com/black_lotus.jpg", Description: "Adds 3 mana of any single color to your mana pool, then is discarded.", SetName: "Alpha", CardNumber: "232", Rarity: "Mythic Rare", TCGGame: "Magic: The Gathering", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: 2, Name: "Charizard", ImageURL: "https://example.com/charizard.jpg", Description: "Spits fire that is hot enough to melt boulders. Known to cause forest fires unintentionally.", SetName: "Base Set", CardNumber: "4", Rarity: "Rare Holo", TCGGame: "Pokémon", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}
	mockDB := MockDBService{
		ListCardsFunc: func() ([]database.Card, error) {
			return cards, nil
		},
	}
	app := fiber.New()
	s := &FiberServer{App: app, db: &mockDB}
	app.Get("/api/cards", s.listCardsHandler)

	req, err := http.NewRequest("GET", "/api/cards", nil)
	if err != nil {
		t.Fatalf("error creating request. Err: %v", err)
	}

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("error making request to server. Err: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status OK; got %v", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error reading response body. Err: %v", err)
	}
	bytes, err := json.Marshal(cards)
	if err != nil {
		t.Fatalf("error marshalling expected cards. Err: %v", err)
	}

	expected := "{\"cards\":" + string(bytes) + "}"
	if expected != string(body) {
		t.Errorf("expected response body to be %v; got %v", expected, string(body))
	}

}

func TestCardByIdHandler(t *testing.T) {
	singleCard := database.Card{ID: 1, Name: "Black Lotus", ImageURL: "https://example.com/black_lotus.jpg", Description: "Adds 3 mana of any single color to your mana pool, then is discarded.", SetName: "Alpha", CardNumber: "232", Rarity: "Mythic Rare", TCGGame: "Magic: The Gathering", CreatedAt: time.Now(), UpdatedAt: time.Now()}
	mockDB := MockDBService{
		GetCardByIDFunc: func() (database.Card, error) {
			return singleCard, nil
		},
	}
	app := fiber.New()
	s := &FiberServer{App: app, db: &mockDB}
	app.Get("/api/cards/:id", s.getCardByIDHandler)

	req, err := http.NewRequest("GET", "/api/cards/1", nil)
	if err != nil {
		t.Fatalf("error creating request. Err: %v", err)
	}

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("error making request to server. Err: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status OK; got %v", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error reading response body. Err: %v", err)
	}
	bytes, err := json.Marshal(singleCard)
	if err != nil {
		t.Fatalf("error marshalling expected card. Err: %v", err)
	}

	expected := "{\"card\":" + string(bytes) + "}"
	if expected != string(body) {
		t.Errorf("expected response body to be %v; got %v", expected, string(body))
	}
}

func TestCreateCardHandler(t *testing.T) {
	mockDB := MockDBService{
		CreateCardFunc: func(card database.CardRequest) error {
			return nil
		},
	}
	app := fiber.New()
	s := &FiberServer{App: app, db: &mockDB}
	app.Post("/api/cards", s.createCardHandler)

	cardRequest := database.CardRequest{
		Name:        "Black Lotus",
		ImageURL:    "https://example.com/black_lotus.jpg",
		Description: "Adds 3 mana of any single color to your mana pool, then is discarded.",
		SetName:     "Alpha",
		CardNumber:  "232",
		Rarity:      "Mythic Rare",
		TCGGameID:   1,
	}
	cardRequestBytes, err := json.Marshal(cardRequest)
	if err != nil {
		t.Fatalf("error marshalling card request. Err: %v", err)
	}

	req, err := http.NewRequest("POST", "/api/cards", strings.NewReader(string(cardRequestBytes)))
	if err != nil {
		t.Fatalf("error creating request. Err: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("error making request to server. Err: %v", err)
	}

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("expected status Created; got %v", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error reading response body. Err: %v", err)
	}

	expected := "{\"message\":\"card created\"}"
	if expected != string(body) {
		t.Errorf("expected response body to be %v; got %v", expected, string(body))
	}
}

func TestUpdateCardHandler(t *testing.T) {
	mockDB := MockDBService{
		UpdateCardFunc: func(cardID int, card database.CardRequest) error {
			return nil
		},
	}
	app := fiber.New()
	s := &FiberServer{App: app, db: &mockDB}
	app.Put("/api/cards/:id", s.updateCardHandler)

	cardRequest := database.CardRequest{
		Name:        "Black Lotus",
		ImageURL:    "https://example.com/black_lotus.jpg",
		Description: "Adds 3 mana of any single color to your mana pool, then is discarded.",
		SetName:     "Alpha",
		CardNumber:  "232",
		Rarity:      "Rare",
		TCGGameID:   1,
	}
	cardRequestBytes, err := json.Marshal(cardRequest)
	if err != nil {
		t.Fatalf("error marshalling card request. Err: %v", err)
	}

	req, err := http.NewRequest("PUT", "/api/cards/1", strings.NewReader(string(cardRequestBytes)))
	if err != nil {
		t.Fatalf("error creating request. Err: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("error making request to server. Err: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status OK; got %v", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error reading response body. Err: %v", err)
	}

	expected := "{\"message\":\"card updated\"}"
	if expected != string(body) {
		t.Errorf("expected response body to be %v; got %v", expected, string(body))
	}
}

func TestDeleteCardHandler(t *testing.T) {
	mockDB := MockDBService{
		DeleteCardFunc: func(cardID int) error {
			return nil
		},
	}
	app := fiber.New()
	s := &FiberServer{App: app, db: &mockDB}
	app.Delete("/api/cards/:id", s.deleteCardHandler)

	req, err := http.NewRequest("DELETE", "/api/cards/1", nil)
	if err != nil {
		t.Fatalf("error creating request. Err: %v", err)
	}

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("error making request to server. Err: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status OK; got %v", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error reading response body. Err: %v", err)
	}

	expected := "{\"message\":\"card deleted\"}"
	if expected != string(body) {
		t.Errorf("expected response body to be %v; got %v", expected, string(body))
	}
}

func TestListOrdersHandler(t *testing.T) {
	orders := []database.Order{
		{OrderID: 1, Buyer: "john_doe", OrderDate: time.Now(), Total: 99.99, Status: "Processing", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{OrderID: 2, Buyer: "jane_smith", OrderDate: time.Now(), Total: 49.49, Status: "Shipped", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}
	mockDB := MockDBService{
		ListOrdersFunc: func() ([]database.Order, error) {
			return orders, nil
		},
	}
	app := fiber.New()
	s := &FiberServer{App: app, db: &mockDB}
	app.Get("/api/orders", s.ListOrdersHandler)

	req, err := http.NewRequest("GET", "/api/orders", nil)
	if err != nil {
		t.Fatalf("error creating request. Err: %v", err)
	}

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("error making request to server. Err: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status OK; got %v", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error reading response body. Err: %v", err)
	}
	bytes, err := json.Marshal(orders)
	if err != nil {
		t.Fatalf("error marshalling expected orders. Err: %v", err)
	}

	expected := "{\"orders\":" + string(bytes) + "}"
	if expected != string(body) {
		t.Errorf("expected response body to be %v; got %v", expected, string(body))
	}

}

func TestGetOrderByIDHandler(t *testing.T) {
	singleOrder := database.Order{OrderID: 1, Buyer: "john_doe", OrderDate: time.Now(), Total: 99.99, Status: "Processing", CreatedAt: time.Now(), UpdatedAt: time.Now()}
	mockDB := MockDBService{
		GetOrderByIDFunc: func(orderID int) (database.Order, error) {
			return singleOrder, nil
		},
	}
	app := fiber.New()
	s := &FiberServer{App: app, db: &mockDB}
	app.Get("/api/orders/:id", s.GetOrderByIDHandler)

	req, err := http.NewRequest("GET", "/api/orders/1", nil)
	if err != nil {
		t.Fatalf("error creating request. Err: %v", err)
	}

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("error making request to server. Err: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status OK; got %v", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error reading response body. Err: %v", err)
	}
	bytes, err := json.Marshal(singleOrder)
	if err != nil {
		t.Fatalf("error marshalling expected order. Err: %v", err)
	}

	expected := "{\"order\":" + string(bytes) + "}"
	if expected != string(body) {
		t.Errorf("expected response body to be %v; got %v", expected, string(body))
	}
}

func TestCreateOrderHandler(t *testing.T) {
	mockDB := MockDBService{
		CreateOrderFunc: func(order database.OrderRequest) error {
			return nil
		},
	}
	app := fiber.New()
	s := &FiberServer{App: app, db: &mockDB}
	app.Post("/api/orders", s.CreateOrderHandler)

	orderRequest := database.OrderRequest{
		BuyerID:   1,
		OrderDate: time.Now(),
		Total:     99.99,
		Status:    "Processing",
	}
	orderRequestBytes, err := json.Marshal(orderRequest)
	if err != nil {
		t.Fatalf("error marshalling order request. Err: %v", err)
	}

	req, err := http.NewRequest("POST", "/api/orders", strings.NewReader(string(orderRequestBytes)))
	if err != nil {
		t.Fatalf("error creating request. Err: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("error making request to server. Err: %v", err)
	}

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("expected status Created; got %v", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error reading response body. Err: %v", err)
	}

	expected := "{\"message\":\"order accepted\"}"
	if expected != string(body) {
		t.Errorf("expected response body to be %v; got %v", expected, string(body))
	}
}

func TestUpdateOrderHandler(t *testing.T) {
	mockDB := MockDBService{
		UpdateOrderFunc: func(orderID int, order database.OrderRequest) error {
			return nil
		},
	}
	app := fiber.New()
	s := &FiberServer{App: app, db: &mockDB}
	app.Put("/api/orders/:id", s.UpdateOrderHandler)

	orderRequest := database.OrderRequest{
		BuyerID:   1,
		OrderDate: time.Now(),
		Total:     89.99,
		Status:    "Shipped",
	}
	orderRequestBytes, err := json.Marshal(orderRequest)
	if err != nil {
		t.Fatalf("error marshalling order request. Err: %v", err)
	}

	req, err := http.NewRequest("PUT", "/api/orders/1", strings.NewReader(string(orderRequestBytes)))
	if err != nil {
		t.Fatalf("error creating request. Err: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("error making request to server. Err: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status OK; got %v", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error reading response body. Err: %v", err)
	}

	expected := "{\"message\":\"order updated\"}"
	if expected != string(body) {
		t.Errorf("expected response body to be %v; got %v", expected, string(body))
	}
}

func TestDeleteOrderHandler(t *testing.T) {
	mockDB := MockDBService{
		DeleteOrderFunc: func(orderID int) error {
			return nil
		},
	}
	app := fiber.New()
	s := &FiberServer{App: app, db: &mockDB}
	app.Delete("/api/orders/:id", s.DeleteOrderHandler)

	req, err := http.NewRequest("DELETE", "/api/orders/1", nil)
	if err != nil {
		t.Fatalf("error creating request. Err: %v", err)
	}

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("error making request to server. Err: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status OK; got %v", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error reading response body. Err: %v", err)
	}

	expected := "{\"message\":\"order deleted\"}"
	if expected != string(body) {
		t.Errorf("expected response body to be %v; got %v", expected, string(body))
	}
}

func TestListUsersHandler(t *testing.T) {
	users := []database.User{
		{UserID: 1, Username: "john_doe", Email: "john@example.com", FirstName: "John", LastName: "Doe", StreetName: "23 Main St", City: "Anytown", State: "CA", ZipCode: "12345", Country: "USA", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{UserID: 2, Username: "jane_doe", Email: "jane@example.com", FirstName: "Jane", LastName: "Doe", StreetName: "24 Main St", City: "Anytown", State: "CA", ZipCode: "12345", Country: "USA", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}
	mockDB := MockDBService{
		ListUsersFunc: func() ([]database.User, error) {
			return users, nil
		},
	}
	app := fiber.New()
	s := &FiberServer{App: app, db: &mockDB}
	app.Get("/api/users", s.ListUsersHandler)

	req, err := http.NewRequest("GET", "/api/users", nil)
	if err != nil {
		t.Fatalf("error creating request. Err: %v", err)
	}

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("error making request to server. Err: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status OK; got %v", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error reading response body. Err: %v", err)
	}
	bytes, err := json.Marshal(users)
	if err != nil {
		t.Fatalf("error marshalling users. Err: %v", err)
	}

	expected := "{\"users\":" + string(bytes) + "}"
	if expected != string(body) {
		t.Errorf("expected response body to be %v; got %v", expected, string(body))
	}
}

func TestGetUserByIDHandler(t *testing.T) {
	singleUser := database.User{UserID: 1, Username: "john_doe", Email: "john@example.com", FirstName: "John", LastName: "Doe", StreetName: "23 Main St", City: "Anytown", State: "CA", ZipCode: "12345", Country: "USA", CreatedAt: time.Now(), UpdatedAt: time.Now()}
	mockDB := MockDBService{
		GetUserByIDFunc: func(userID int) (database.User, error) {
			return singleUser, nil
		},
	}
	app := fiber.New()
	s := &FiberServer{App: app, db: &mockDB}
	app.Get("/api/users/:id", s.GetUserByIDHandler)

	req, err := http.NewRequest("GET", "/api/users/1", nil)
	if err != nil {
		t.Fatalf("error creating request. Err: %v", err)
	}

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("error making request to server. Err: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status OK; got %v", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error reading response body. Err: %v", err)
	}
	bytes, err := json.Marshal(singleUser)
	if err != nil {
		t.Fatalf("error marshalling user. Err: %v", err)
	}

	expected := "{\"user\":" + string(bytes) + "}"
	if expected != string(body) {
		t.Errorf("expected response body to be %v; got %v", expected, string(body))
	}
}

func TestCreateUserHandler(t *testing.T) {
	mockDB := MockDBService{
		CreateUserFunc: func(user database.UserRequest) error {
			return nil
		},
	}
	app := fiber.New()
	s := &FiberServer{App: app, db: &mockDB}
	app.Post("/api/users", s.CreateUserHandler)

	userRequest := database.UserRequest{
		Username:     "john_doe",
		Email:        "john@example.com",
		Password:     "newpassword",
		FirstName:    "John",
		LastName:     "Doe",
		StreetName:   "Main St",
		StreetNumber: "23",
		City:         "Anytown",
		State:        "CA",
		ZipCode:      "12345",
		SellerType:   "powerseller",
		LanguageID:   1,
		CountryID:    1,
	}
	userRequestBytes, err := json.Marshal(userRequest)
	if err != nil {
		t.Fatalf("error marshalling user request. Err: %v", err)
	}

	req, err := http.NewRequest("POST", "/api/users", strings.NewReader(string(userRequestBytes)))
	if err != nil {
		t.Fatalf("error creating request. Err: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("error making request to server. Err: %v", err)
	}

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("expected status Created; got %v", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error reading response body. Err: %v", err)
	}

	expected := "{\"message\":\"user created\"}"
	if expected != string(body) {
		t.Errorf("expected response body to be %v; got %v", expected, string(body))
	}
}

func TestUpdateUserHandler(t *testing.T) {
	mockDB := MockDBService{
		UpdateUserFunc: func(userID int, user database.UserRequest) error {
			return nil
		},
	}
	app := fiber.New()
	s := &FiberServer{App: app, db: &mockDB}
	app.Put("/api/users/:id", s.UpdateUserHandler)

	userRequest := database.UserRequest{
		Username:     "john_doe_updated",
		Email:        "john_updated@example.com",
		Password:     "updatedpassword",
		FirstName:    "John",
		LastName:     "Doe",
		StreetName:   "Main St",
		StreetNumber: "23",
		City:         "Anytown",
		State:        "CA",
		ZipCode:      "12345",
		SellerType:   "powerseller",
		LanguageID:   1,
		CountryID:    1,
	}
	userRequestBytes, err := json.Marshal(userRequest)
	if err != nil {
		t.Fatalf("error marshalling user request. Err: %v", err)
	}

	req, err := http.NewRequest("PUT", "/api/users/1", strings.NewReader(string(userRequestBytes)))
	if err != nil {
		t.Fatalf("error creating request. Err: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("error making request to server. Err: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status OK; got %v", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error reading response body. Err: %v", err)
	}

	expected := "{\"message\":\"user updated\"}"
	if expected != string(body) {
		t.Errorf("expected response body to be %v; got %v", expected, string(body))
	}
}

func TestDeleteUserHandler(t *testing.T) {
	mockDB := MockDBService{
		DeleteUserFunc: func(userID int) error {
			return nil
		},
	}
	app := fiber.New()
	s := &FiberServer{App: app, db: &mockDB}
	app.Delete("/api/users/:id", s.DeleteUserHandler)

	req, err := http.NewRequest("DELETE", "/api/users/1", nil)
	if err != nil {
		t.Fatalf("error creating request. Err: %v", err)
	}

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("error making request to server. Err: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status OK; got %v", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error reading response body. Err: %v", err)
	}

	expected := "{\"message\":\"user deleted\"}"
	if expected != string(body) {
		t.Errorf("expected response body to be %v; got %v", expected, string(body))
	}
}
