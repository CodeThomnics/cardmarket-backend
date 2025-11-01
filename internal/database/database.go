package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
)

type Card struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	ImageURL    string    `json:"image_url"`
	Description string    `json:"description"`
	SetName     string    `json:"set_name"`
	CardNumber  string    `json:"card_number"`
	Rarity      string    `json:"rarity"`
	TCGGame     string    `json:"tcg_game"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CardRequest struct {
	Name        string `json:"name"`
	ImageURL    string `json:"image_url"`
	Description string `json:"description"`
	SetName     string `json:"set_name"`
	CardNumber  string `json:"card_number"`
	Rarity      string `json:"rarity"`
	TCGGameID   int    `json:"tcg_game_id"`
}

type Product struct {
	ProductID   int       `json:"product_id"`
	Price       float64   `json:"price"`
	Condition   string    `json:"condition"`
	Quantity    int       `json:"quantity"`
	IsAvailable bool      `json:"is_available"`
	Seller      string    `json:"seller"`
	Card        string    `json:"card"`
	Language    string    `json:"language"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ProductRequest struct {
	ProductID   int     `json:"product_id"`
	Price       float64 `json:"price"`
	Condition   string  `json:"condition"`
	Quantity    int     `json:"quantity"`
	IsAvailable bool    `json:"is_available"`
	SellerID    int     `json:"seller_id"`
	CardID      int     `json:"card_id"`
	LanguageID  int     `json:"language_id"`
}

// Service represents a service that interacts with a database.
type Service interface {
	// Health returns a map of health status information.
	// The keys and values in the map are service-specific.
	Health() map[string]string

	// Close terminates the database connection.
	// It returns an error if the connection cannot be closed.
	Close() error

	// ListCards returns a list of card names from the cards table.
	// TODO: Replace with a richer domain model and filtering/pagination as needed.
	ListCards() ([]Card, error)
	GetCardByID(cardID int) (Card, error)
	CreateCard(card CardRequest) error
	UpdateCard(cardID int, card CardRequest) error
	DeleteCard(cardID int) error

	ListProducts() ([]Product, error)
	GetProductByID(productID int) (Product, error)
	CreateProduct(product ProductRequest) error
	UpdateProduct(productID int, product ProductRequest) error
	DeleteProduct(productID int) error
}

type service struct {
	db *sql.DB
}

var (
	database   = os.Getenv("DB_DATABASE")
	password   = os.Getenv("DB_PASSWORD")
	username   = os.Getenv("DB_USERNAME")
	port       = os.Getenv("DB_PORT")
	host       = os.Getenv("DB_HOST")
	schema     = os.Getenv("DB_SCHEMA")
	dbInstance *service
)

func New() Service {
	// Reuse Connection
	if dbInstance != nil {
		return dbInstance
	}
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&search_path=%s", username, password, host, port, database, schema)
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatal(err)
	}
	dbInstance = &service{
		db: db,
	}
	return dbInstance
}

// Health checks the health of the database connection by pinging the database.
// It returns a map with keys indicating various health statistics.
func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	stats := make(map[string]string)

	// Ping the database
	err := s.db.PingContext(ctx)
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("db down: %v", err)
		log.Fatalf("db down: %v", err) // Log the error and terminate the program
		return stats
	}

	// Database is up, add more statistics
	stats["status"] = "up"
	stats["message"] = "It's healthy"

	// Get database stats (like open connections, in use, idle, etc.)
	dbStats := s.db.Stats()
	stats["open_connections"] = strconv.Itoa(dbStats.OpenConnections)
	stats["in_use"] = strconv.Itoa(dbStats.InUse)
	stats["idle"] = strconv.Itoa(dbStats.Idle)
	stats["wait_count"] = strconv.FormatInt(dbStats.WaitCount, 10)
	stats["wait_duration"] = dbStats.WaitDuration.String()
	stats["max_idle_closed"] = strconv.FormatInt(dbStats.MaxIdleClosed, 10)
	stats["max_lifetime_closed"] = strconv.FormatInt(dbStats.MaxLifetimeClosed, 10)

	// Evaluate stats to provide a health message
	if dbStats.OpenConnections > 40 { // Assuming 50 is the max for this example
		stats["message"] = "The database is experiencing heavy load."
	}

	if dbStats.WaitCount > 1000 {
		stats["message"] = "The database has a high number of wait events, indicating potential bottlenecks."
	}

	if dbStats.MaxIdleClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many idle connections are being closed, consider revising the connection pool settings."
	}

	if dbStats.MaxLifetimeClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many connections are being closed due to max lifetime, consider increasing max lifetime or revising the connection usage pattern."
	}

	return stats
}

// Close closes the database connection.
// It logs a message indicating the disconnection from the specific database.
// If the connection is successfully closed, it returns nil.
// If an error occurs while closing the connection, it returns the error.
func (s *service) Close() error {
	log.Printf("Disconnected from database: %s", database)
	return s.db.Close()
}

func (s *service) ListCards() ([]Card, error) {
	rows, err := s.db.Query("SELECT c.card_id, c.name, c.image_url, c.description, c.set_name, c.card_number, c.rarity, tcg.name AS tcg_game,c.created_at,c.updated_at FROM cards c JOIN tcg_games tcg ON c.tcg_game_id = tcg.tcg_game_id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var cards []Card
	for rows.Next() {
		var card Card
		if err := rows.Scan(&card.ID, &card.Name, &card.ImageURL, &card.Description, &card.SetName, &card.CardNumber, &card.Rarity, &card.TCGGame, &card.CreatedAt, &card.UpdatedAt); err != nil {
			return nil, err
		}
		cards = append(cards, card)
	}
	return cards, nil
}

func (s *service) GetCardByID(cardID int) (Card, error) {
	var card Card
	err := s.db.QueryRow("SELECT c.card_id, c.name, c.image_url, c.description, c.set_name, c.card_number, c.rarity, tcg.name AS tcg_game,c.created_at,c.updated_at FROM cards c JOIN tcg_games tcg ON c.tcg_game_id = tcg.tcg_game_id WHERE c.card_id = $1", cardID).Scan(&card.ID, &card.Name, &card.ImageURL, &card.Description, &card.SetName, &card.CardNumber, &card.Rarity, &card.TCGGame, &card.CreatedAt, &card.UpdatedAt)
	if err != nil {
		return Card{}, err
	}
	return card, nil
}

func (s *service) CreateCard(card CardRequest) error {
	query := `INSERT INTO cards (name, image_url, description, set_name, card_number, rarity, tcg_game_id) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := s.db.Exec(query, card.Name, card.ImageURL, card.Description, card.SetName, card.CardNumber, card.Rarity, card.TCGGameID)
	return err
}

func (s *service) UpdateCard(cardID int, card CardRequest) error {
	query := `UPDATE cards SET name = $1, image_url = $2, description = $3, set_name = $4, card_number = $5, rarity = $6, tcg_game_id = $7, updated_at = CURRENT_TIMESTAMP WHERE card_id = $8`

	result, err := s.db.Exec(query, card.Name, card.ImageURL, card.Description, card.SetName, card.CardNumber, card.Rarity, card.TCGGameID, cardID)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (s *service) DeleteCard(cardID int) error {
	query := `DELETE FROM cards WHERE card_id = $1`

	result, err := s.db.Exec(query, cardID)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (s *service) ListProducts() ([]Product, error) {
	query := `SELECT p.product_id, p.price, p.condition, p.quantity, p.is_available, us.username AS seller, c.name AS card, l.language_name AS language, p.created_at, p.updated_at FROM products p JOIN users us ON p.seller_id = us.user_id JOIN cards c ON p.card_id = c.card_id JOIN languages l ON p.language_id = l.language_id`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var product Product
		if err := rows.Scan(&product.ProductID, &product.Price, &product.Condition, &product.Quantity, &product.IsAvailable, &product.Seller, &product.Card, &product.Language, &product.CreatedAt, &product.UpdatedAt); err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	return products, nil
}

func (s *service) GetProductByID(productID int) (Product, error) {
	var product Product
	query := `SELECT p.product_id, p.price, p.condition, p.quantity, p.is_available, us.username AS seller, c.name AS card, l.language_name AS language, p.created_at, p.updated_at FROM products p JOIN users us ON p.seller_id = us.user_id JOIN cards c ON p.card_id = c.card_id JOIN languages l ON p.language_id = l.language_id WHERE p.product_id = $1`
	err := s.db.QueryRow(query, productID).Scan(&product.ProductID, &product.Price, &product.Condition, &product.Quantity, &product.IsAvailable, &product.Seller, &product.Card, &product.Language, &product.CreatedAt, &product.UpdatedAt)
	if err != nil {
		return Product{}, err
	}
	return product, nil
}

func (s *service) CreateProduct(product ProductRequest) error {
	query := `INSERT INTO products (price, condition, quantity, is_available, seller_id, card_id, language_id) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := s.db.Exec(query, product.Price, product.Condition, product.Quantity, product.IsAvailable, product.SellerID, product.CardID, product.LanguageID)
	return err
}

func (s *service) UpdateProduct(productID int, product ProductRequest) error {
	query := `UPDATE products SET price = $1, condition = $2, quantity = $3, is_available = $4, seller_id = $5, card_id = $6, language_id = $7, updated_at = CURRENT_TIMESTAMP WHERE product_id = $8`

	result, err := s.db.Exec(query, product.Price, product.Condition, product.Quantity, product.IsAvailable, product.SellerID, product.CardID, product.LanguageID, productID)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (s *service) DeleteProduct(productID int) error {
	query := `DELETE FROM products WHERE product_id = $1`

	result, err := s.db.Exec(query, productID)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
