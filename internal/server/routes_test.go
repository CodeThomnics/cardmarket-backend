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
	ListCardsFunc   func() ([]database.Card, error)
	GetCardByIDFunc func() (database.Card, error)
	CreateCardFunc  func(card database.CardRequest) error
	UpdateCardFunc  func(cardID int, card database.CardRequest) error
	DeleteCardFunc  func(cardID int) error
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
		{ID: 1, Name: "Black Lotus", ImageURL: "https://example.com/black_lotus.jpg", Description: "Adds 3 mana of any single color to your mana pool, then is discarded.", SetName: "Alpha", CardNumber: "232", Rarity: "Mythic Rare", TCGGame: "Magic: The Gathering", Language: "English", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: 2, Name: "Charizard", ImageURL: "https://example.com/charizard.jpg", Description: "Spits fire that is hot enough to melt boulders. Known to cause forest fires unintentionally.", SetName: "Base Set", CardNumber: "4", Rarity: "Rare Holo", TCGGame: "Pokémon", Language: "English", CreatedAt: time.Now(), UpdatedAt: time.Now()},
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
	singleCard := database.Card{ID: 1, Name: "Black Lotus", ImageURL: "https://example.com/black_lotus.jpg", Description: "Adds 3 mana of any single color to your mana pool, then is discarded.", SetName: "Alpha", CardNumber: "232", Rarity: "Mythic Rare", TCGGame: "Magic: The Gathering", Language: "English", CreatedAt: time.Now(), UpdatedAt: time.Now()}
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
		LanguageID:  1,
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
		Rarity:      "Mythic Rare",
		LanguageID:  2,
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
