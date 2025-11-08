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

type Order struct {
	OrderID         int        `json:"order_id"`
	Buyer           string     `json:"buyer"`
	Seller          string     `json:"seller"`
	Quantity        int        `json:"quantity"`
	Product         string     `json:"product"`
	OrderDate       time.Time  `json:"order_date"`
	ShippingAddress string     `json:"shipping_address"`
	ShippingCost    float64    `json:"shipping_cost"`
	Total           float64    `json:"total"`
	Status          string     `json:"status"`
	TrackingNumber  *string    `json:"tracking_number,omitempty"`
	ShippedAt       *time.Time `json:"shipped_at,omitempty"`
	DeliveredAt     *time.Time `json:"delivered_at,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

type OrderRequest struct {
	BuyerID         int        `json:"buyer_id"`
	SellerID        int        `json:"seller_id"`
	Quantity        int        `json:"quantity"`
	ProductID       int        `json:"product_id"`
	OrderDate       time.Time  `json:"order_date"`
	Total           float64    `json:"total"`
	ShippingAddress string     `json:"shipping_address"`
	ShippingCost    float64    `json:"shipping_cost"`
	TrackingNumber  *string    `json:"tracking_number,omitempty"`
	ShippedAt       *time.Time `json:"shipped_at,omitempty"`
	DeliveredAt     *time.Time `json:"delivered_at,omitempty"`
	Status          string     `json:"status"`
}

type User struct {
	UserID       int       `json:"user_id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	StreetName   string    `json:"street_name"`
	StreetNumber string    `json:"street_number"`
	City         string    `json:"city"`
	State        string    `json:"state"`
	ZipCode      string    `json:"zip_code"`
	SellerType   string    `json:"seller_type"`
	Country      string    `json:"country"`
	Language     string    `json:"language"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type UserRequest struct {
	Username     string `json:"username"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	StreetName   string `json:"street_name"`
	StreetNumber string `json:"street_number"`
	City         string `json:"city"`
	State        string `json:"state"`
	ZipCode      string `json:"zip_code"`
	SellerType   string `json:"seller_type"`
	CountryID    int    `json:"country_id"`
	LanguageID   int    `json:"language_id"`
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

	ListOrders() ([]Order, error)
	GetOrderByID(orderID int) (Order, error)
	CreateOrder(order OrderRequest) error
	UpdateOrder(orderID int, order OrderRequest) error
	DeleteOrder(orderID int) error

	ListUsers() ([]User, error)
	GetUserByID(userID int) (User, error)
	CreateUser(user UserRequest) error
	UpdateUser(userID int, user UserRequest) error
	DeleteUser(userID int) error
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
	if err := rows.Err(); err != nil {
		return nil, err
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
	if err := rows.Err(); err != nil {
		return nil, err
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

func (s *service) ListOrders() ([]Order, error) {
	query := `SELECT o.order_id, buyers.username AS buyer, sellers.username AS seller, o.quantity, c.name AS product, o.order_date, o.shipping_address, o.shipping_cost, o.total_amount, o.tracking_number, o.shipped_at, o.delivered_at, o.status, o.created_at, o.updated_at FROM orders o JOIN users buyers ON o.buyer_id = buyers.user_id JOIN users sellers ON o.seller_id = sellers.user_id JOIN products product ON o.product_id = product.product_id JOIN cards c ON product.card_id = c.card_id`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []Order
	for rows.Next() {
		var order Order
		if err := rows.Scan(&order.OrderID, &order.Buyer, &order.Seller, &order.Quantity, &order.Product, &order.OrderDate, &order.ShippingAddress, &order.ShippingCost, &order.Total, &order.TrackingNumber, &order.ShippedAt, &order.DeliveredAt, &order.Status, &order.CreatedAt, &order.UpdatedAt); err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return orders, nil
}

func (s *service) GetOrderByID(orderID int) (Order, error) {
	var order Order
	query := `SELECT o.order_id, buyers.username AS buyer, sellers.username AS seller, o.quantity, c.name AS product, o.order_date, o.shipping_address, o.shipping_cost, o.total_amount, o.tracking_number, o.shipped_at, o.delivered_at, o.status, o.created_at, o.updated_at FROM orders o JOIN users buyers ON o.buyer_id = buyers.user_id JOIN users sellers ON o.seller_id = sellers.user_id JOIN products product ON o.product_id = product.product_id JOIN cards c ON product.card_id = c.card_id WHERE o.order_id = $1`
	err := s.db.QueryRow(query, orderID).Scan(&order.OrderID, &order.Buyer, &order.Seller, &order.Quantity, &order.Product, &order.OrderDate, &order.ShippingAddress, &order.ShippingCost, &order.Total, &order.TrackingNumber, &order.ShippedAt, &order.DeliveredAt, &order.Status, &order.CreatedAt, &order.UpdatedAt)
	if err != nil {
		return Order{}, err
	}
	return order, nil
}

func (s *service) CreateOrder(order OrderRequest) error {
	query := `INSERT INTO orders (buyer_id, seller_id, product_id, quantity, order_date, shipping_address, shipping_cost, total_amount, tracking_number, shipped_at, delivered_at, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`
	_, err := s.db.Exec(query, order.BuyerID, order.SellerID, order.ProductID, order.Quantity, order.OrderDate, order.ShippingAddress, order.ShippingCost, order.Total, order.TrackingNumber, order.ShippedAt, order.DeliveredAt, order.Status)
	return err
}

func (s *service) UpdateOrder(orderID int, order OrderRequest) error {
	query := `UPDATE orders SET buyer_id = $1, seller_id = $2, product_id = $3, quantity = $4, order_date = $5, shipping_address = $6, shipping_cost = $7, total_amount = $8, tracking_number = $9, shipped_at = $10, delivered_at = $11, status = $12, updated_at = CURRENT_TIMESTAMP WHERE order_id = $13`

	result, err := s.db.Exec(query, order.BuyerID, order.SellerID, order.ProductID, order.Quantity, order.OrderDate, order.ShippingAddress, order.ShippingCost, order.Total, order.TrackingNumber, order.ShippedAt, order.DeliveredAt, order.Status, orderID)

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

func (s *service) DeleteOrder(orderID int) error {
	query := `DELETE FROM orders WHERE order_id = $1`

	result, err := s.db.Exec(query, orderID)

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

func (s *service) ListUsers() ([]User, error) {
	rows, err := s.db.Query("SELECT u.user_id, u.username, u.email, u.first_name, u.last_name, u.street_name, u.street_number, u.city, u.state, u.zip_code, u.seller_type, c.country_name, l.language_name, u.created_at, u.updated_at FROM users u JOIN countries c ON u.country_id = c.country_id JOIN languages l ON u.language_id = l.language_id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.UserID, &user.Username, &user.Email, &user.FirstName, &user.LastName, &user.StreetName, &user.StreetNumber, &user.City, &user.State, &user.ZipCode, &user.SellerType, &user.Country, &user.Language, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

func (s *service) GetUserByID(userID int) (User, error) {
	var user User
	err := s.db.QueryRow("SELECT u.user_id, u.username, u.email, u.first_name, u.last_name, u.street_name, u.street_number, u.city, u.state, u.zip_code, u.seller_type, c.country_name, l.language_name, u.created_at, u.updated_at FROM users u JOIN countries c ON u.country_id = c.country_id JOIN languages l ON u.language_id = l.language_id WHERE u.user_id = $1", userID).Scan(&user.UserID, &user.Username, &user.Email, &user.FirstName, &user.LastName, &user.StreetName, &user.StreetNumber, &user.City, &user.State, &user.ZipCode, &user.SellerType, &user.Country, &user.Language, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func (s *service) CreateUser(user UserRequest) error {
	query := `INSERT INTO users (username, email, password, first_name, last_name, street_name, street_number, city, state, zip_code, seller_type, country_id, language_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`
	_, err := s.db.Exec(query, user.Username, user.Email, user.Password, user.FirstName, user.LastName, user.StreetName, user.StreetNumber, user.City, user.State, user.ZipCode, user.SellerType, user.CountryID, user.LanguageID)
	return err
}

func (s *service) UpdateUser(userID int, user UserRequest) error {
	query := `UPDATE users SET username = $1, email = $2, password = $3, first_name = $4, last_name = $5, street_name = $6, street_number = $7, city = $8, state = $9, zip_code = $10, seller_type = $11, country_id = $12, language_id = $13, updated_at = CURRENT_TIMESTAMP WHERE user_id = $14`

	result, err := s.db.Exec(query, user.Username, user.Email, user.Password, user.FirstName, user.LastName, user.StreetName, user.StreetNumber, user.City, user.State, user.ZipCode, user.CountryID, user.LanguageID, userID)

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

func (s *service) DeleteUser(userID int) error {
	query := `DELETE FROM users WHERE user_id = $1`

	result, err := s.db.Exec(query, userID)

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
