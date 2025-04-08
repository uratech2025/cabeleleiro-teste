package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	_ "modernc.org/sqlite"
)

type Service struct {
	ID         int     `json:"id"`
	Name       string  `json:"name"`
	Price      float64 `json:"price"`
	Commission float64 `json:"commission"`
}

type Client struct {
	ID      int     `json:"id"`
	Name    string  `json:"name"`
	Balance float64 `json:"balance"`
}

type Order struct {
	ID        int       `json:"id"`
	ClientID  int       `json:"client_id"`
	ServiceID int       `json:"service_id"`
	Date      time.Time `json:"date"`
	Status    string    `json:"status"`
	Total     float64   `json:"total"`
}

type CashFlow struct {
	ID          int       `json:"id"`
	Type        string    `json:"type"` // "entrada" or "saida"
	Description string    `json:"description"`
	Amount      float64   `json:"amount"`
	Date        time.Time `json:"date"`
}

type InventoryItem struct {
	ID       int     `json:"id"`
	ItemName string  `json:"item_name"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("sqlite", "./salon.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create tables
	createTables()

	r := gin.Default()

	// Servir arquivos estáticos do frontend
	r.Static("/static", "./frontend")
	r.StaticFile("/", "./frontend/index.html")
	r.StaticFile("/index.html", "./frontend/index.html")
	r.StaticFile("/styles.css", "./frontend/styles.css")
	r.StaticFile("/script.js", "./frontend/script.js")

	// Rota principal
	r.GET("/", func(c *gin.Context) {
		c.File("./frontend/index.html")
	})

	// Enable CORS
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Routes
	// Services
	r.GET("/services", getServices)
	r.POST("/services", createService)
	r.PUT("/services/:id", updateService)
	r.DELETE("/services/:id", deleteService)

	// Clients
	r.GET("/clients", getClients)
	r.POST("/clients", createClient)
	r.PUT("/clients/:id", updateClient)
	r.DELETE("/clients/:id", deleteClient)

	// Orders
	r.GET("/orders", getOrders)
	r.POST("/orders", createOrder)
	r.PUT("/orders/:id/status", updateOrderStatus)
	r.DELETE("/orders/:id", deleteOrder)

	// CashFlow
	r.GET("/cashflow", getCashFlow)
	r.POST("/cashflow", createCashFlow)
	r.DELETE("/cashflow/:id", deleteCashFlow)

	// Inventory
	r.GET("/inventory", getInventory)
	r.POST("/inventory", createInventoryItem)
	r.PUT("/inventory/:id", updateInventoryItem)
	r.DELETE("/inventory/:id", deleteInventoryItem)

	// Obter a porta da variável de ambiente ou usar o padrão
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Iniciar o servidor no host 0.0.0.0
	r.Run("0.0.0.0:" + port)
}

func createTables() {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS services (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			price REAL NOT NULL,
			commission REAL NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS clients (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			balance REAL DEFAULT 0
		)`,
		`CREATE TABLE IF NOT EXISTS orders (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			client_id INTEGER,
			service_id INTEGER,
			date DATETIME DEFAULT CURRENT_TIMESTAMP,
			status TEXT,
			total REAL,
			FOREIGN KEY (client_id) REFERENCES clients(id),
			FOREIGN KEY (service_id) REFERENCES services(id)
		)`,
		`CREATE TABLE IF NOT EXISTS cashflow (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			type TEXT NOT NULL,
			description TEXT,
			amount REAL NOT NULL,
			date DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS inventory (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			item_name TEXT NOT NULL,
			quantity INTEGER NOT NULL,
			price REAL NOT NULL
		)`,
	}

	for _, query := range queries {
		_, err := db.Exec(query)
		if err != nil {
			log.Fatal(err)
		}
	}
}

// Service handlers
func getServices(c *gin.Context) {
	rows, err := db.Query("SELECT id, name, price, commission FROM services")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var services []Service
	for rows.Next() {
		var s Service
		if err := rows.Scan(&s.ID, &s.Name, &s.Price, &s.Commission); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		services = append(services, s)
	}

	c.JSON(http.StatusOK, services)
}

func createService(c *gin.Context) {
	var service Service
	if err := c.ShouldBindJSON(&service); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := db.Exec("INSERT INTO services (name, price, commission) VALUES (?, ?, ?)",
		service.Name, service.Price, service.Commission)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	id, _ := result.LastInsertId()
	service.ID = int(id)
	c.JSON(http.StatusCreated, service)
}

func updateService(c *gin.Context) {
	id := c.Param("id")
	var service Service
	if err := c.ShouldBindJSON(&service); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := db.Exec("UPDATE services SET name = ?, price = ?, commission = ? WHERE id = ?",
		service.Name, service.Price, service.Commission, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Service updated successfully"})
}

func deleteService(c *gin.Context) {
	id := c.Param("id")
	_, err := db.Exec("DELETE FROM services WHERE id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Service deleted successfully"})
}

// Client handlers
func getClients(c *gin.Context) {
	rows, err := db.Query("SELECT id, name, balance FROM clients")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var clients []Client
	for rows.Next() {
		var client Client
		if err := rows.Scan(&client.ID, &client.Name, &client.Balance); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		clients = append(clients, client)
	}

	c.JSON(http.StatusOK, clients)
}

func createClient(c *gin.Context) {
	var client Client
	if err := c.ShouldBindJSON(&client); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := db.Exec("INSERT INTO clients (name, balance) VALUES (?, ?)",
		client.Name, client.Balance)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	id, _ := result.LastInsertId()
	client.ID = int(id)
	c.JSON(http.StatusCreated, client)
}

func updateClient(c *gin.Context) {
	id := c.Param("id")
	var client Client
	if err := c.ShouldBindJSON(&client); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := db.Exec("UPDATE clients SET name = ?, balance = ? WHERE id = ?",
		client.Name, client.Balance, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Client updated successfully"})
}

func deleteClient(c *gin.Context) {
	id := c.Param("id")
	_, err := db.Exec("DELETE FROM clients WHERE id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Client deleted successfully"})
}

// Order handlers
func getOrders(c *gin.Context) {
	rows, err := db.Query(`
		SELECT o.id, o.client_id, o.service_id, o.date, o.status, o.total,
			   c.name as client_name, s.name as service_name
		FROM orders o
		JOIN clients c ON o.client_id = c.id
		JOIN services s ON o.service_id = s.id
		ORDER BY o.date DESC
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var orders []map[string]interface{}
	for rows.Next() {
		var order Order
		var clientName, serviceName string
		if err := rows.Scan(&order.ID, &order.ClientID, &order.ServiceID, &order.Date,
			&order.Status, &order.Total, &clientName, &serviceName); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		orders = append(orders, map[string]interface{}{
			"id":           order.ID,
			"client_name":  clientName,
			"service_name": serviceName,
			"date":         order.Date,
			"status":       order.Status,
			"total":        order.Total,
		})
	}

	c.JSON(http.StatusOK, orders)
}

func createOrder(c *gin.Context) {
	var order Order
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx, err := db.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get service price
	var price float64
	err = tx.QueryRow("SELECT price FROM services WHERE id = ?", order.ServiceID).Scan(&price)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Create order
	result, err := tx.Exec("INSERT INTO orders (client_id, service_id, status, total) VALUES (?, ?, ?, ?)",
		order.ClientID, order.ServiceID, "pending", price)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	id, _ := result.LastInsertId()
	order.ID = int(id)
	order.Total = price

	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, order)
}

func updateOrderStatus(c *gin.Context) {
	id := c.Param("id")
	var order Order
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx, err := db.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Update order status
	_, err = tx.Exec("UPDATE orders SET status = ? WHERE id = ?", order.Status, id)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// If order is completed, update client balance and create cashflow entry
	if order.Status == "completed" {
		// Get order details
		var total float64
		var clientID int
		err = tx.QueryRow("SELECT total, client_id FROM orders WHERE id = ?", id).Scan(&total, &clientID)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Update client balance
		_, err = tx.Exec("UPDATE clients SET balance = balance - ? WHERE id = ?", total, clientID)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Create cashflow entry
		_, err = tx.Exec("INSERT INTO cashflow (type, description, amount) VALUES (?, ?, ?)",
			"entrada", "Pagamento de serviço", total)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order status updated successfully"})
}

func deleteOrder(c *gin.Context) {
	id := c.Param("id")
	_, err := db.Exec("DELETE FROM orders WHERE id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order deleted successfully"})
}

// CashFlow handlers
func getCashFlow(c *gin.Context) {
	rows, err := db.Query("SELECT id, type, description, amount, date FROM cashflow ORDER BY date DESC")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var cashflows []CashFlow
	for rows.Next() {
		var cf CashFlow
		if err := rows.Scan(&cf.ID, &cf.Type, &cf.Description, &cf.Amount, &cf.Date); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		cashflows = append(cashflows, cf)
	}

	c.JSON(http.StatusOK, cashflows)
}

func createCashFlow(c *gin.Context) {
	var cf CashFlow
	if err := c.ShouldBindJSON(&cf); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := db.Exec("INSERT INTO cashflow (type, description, amount) VALUES (?, ?, ?)",
		cf.Type, cf.Description, cf.Amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	id, _ := result.LastInsertId()
	cf.ID = int(id)
	c.JSON(http.StatusCreated, cf)
}

func deleteCashFlow(c *gin.Context) {
	id := c.Param("id")
	_, err := db.Exec("DELETE FROM cashflow WHERE id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "CashFlow entry deleted successfully"})
}

// Inventory handlers
func getInventory(c *gin.Context) {
	rows, err := db.Query("SELECT id, item_name, quantity, price FROM inventory")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var items []map[string]interface{}
	for rows.Next() {
		var id int
		var name string
		var quantity int
		var price float64
		if err := rows.Scan(&id, &name, &quantity, &price); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		items = append(items, map[string]interface{}{
			"id":       id,
			"name":     name,
			"quantity": quantity,
			"price":    price,
		})
	}

	c.JSON(http.StatusOK, items)
}

func createInventoryItem(c *gin.Context) {
	var item InventoryItem
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := db.Exec("INSERT INTO inventory (item_name, quantity, price) VALUES (?, ?, ?)",
		item.ItemName, item.Quantity, item.Price)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	id, _ := result.LastInsertId()
	item.ID = int(id)
	c.JSON(http.StatusCreated, item)
}

func updateInventoryItem(c *gin.Context) {
	id := c.Param("id")
	var item InventoryItem
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := db.Exec("UPDATE inventory SET item_name = ?, quantity = ?, price = ? WHERE id = ?",
		item.ItemName, item.Quantity, item.Price, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Inventory item updated successfully"})
}

func deleteInventoryItem(c *gin.Context) {
	id := c.Param("id")
	_, err := db.Exec("DELETE FROM inventory WHERE id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Inventory item deleted successfully"})
}
