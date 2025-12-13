package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

// Person represents a person in the meibo (directory)
type Person struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

var db *sql.DB

func main() {
	// Get database configuration from environment variables
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	if dbHost == "" || dbUser == "" || dbPassword == "" || dbName == "" {
		log.Fatal("Database environment variables (DB_HOST, DB_USER, DB_PASSWORD, DB_NAME) must be set")
	}

	// Connect to database
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", dbUser, dbPassword, dbHost, dbName)
	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Test database connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	log.Println("Connected to database successfully")

	// Initialize database table
	if err := initDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	r := mux.NewRouter()

	// API routes
	r.HandleFunc("/", healthHandler).Methods("GET")
	r.HandleFunc("/api/health", healthHandler).Methods("GET")
	r.HandleFunc("/api/persons", getPersonsHandler).Methods("GET")
	r.HandleFunc("/api/persons", createPersonHandler).Methods("POST")
	r.HandleFunc("/api/persons/{id}", getPersonHandler).Methods("GET")
	r.HandleFunc("/api/persons/{id}", updatePersonHandler).Methods("PUT")
	r.HandleFunc("/api/persons/{id}", deletePersonHandler).Methods("DELETE")

	// CORS settings
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: false,
	})

	handler := c.Handler(r)

	log.Println("Meibo API Server starting on :80")
	log.Fatal(http.ListenAndServe(":80", handler))
}

// Initialize database table
func initDB() error {
	query := `
	CREATE TABLE IF NOT EXISTS persons (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		email VARCHAR(255),
		phone VARCHAR(50),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
	)`
	_, err := db.Exec(query)
	return err
}

// Health check endpoint
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "ok",
		"service": "meibo-api",
	})
}

// Get all persons
func getPersonsHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, name, email, phone FROM persons ORDER BY id")
	if err != nil {
		http.Error(w, "Failed to fetch persons", http.StatusInternalServerError)
		log.Printf("Error fetching persons: %v", err)
		return
	}
	defer rows.Close()

	var persons []Person
	for rows.Next() {
		var p Person
		var email, phone sql.NullString
		if err := rows.Scan(&p.ID, &p.Name, &email, &phone); err != nil {
			http.Error(w, "Failed to scan person", http.StatusInternalServerError)
			log.Printf("Error scanning person: %v", err)
			return
		}
		p.Email = email.String
		p.Phone = phone.String
		persons = append(persons, p)
	}

	if persons == nil {
		persons = []Person{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(persons)
}

// Get a single person by ID
func getPersonHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var p Person
	var email, phone sql.NullString
	err = db.QueryRow("SELECT id, name, email, phone FROM persons WHERE id = ?", id).Scan(&p.ID, &p.Name, &email, &phone)
	if err == sql.ErrNoRows {
		http.Error(w, "Person not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Failed to fetch person", http.StatusInternalServerError)
		log.Printf("Error fetching person: %v", err)
		return
	}
	p.Email = email.String
	p.Phone = phone.String

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

// Create a new person
func createPersonHandler(w http.ResponseWriter, r *http.Request) {
	var p Person
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if p.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	result, err := db.Exec("INSERT INTO persons (name, email, phone) VALUES (?, ?, ?)", p.Name, p.Email, p.Phone)
	if err != nil {
		http.Error(w, "Failed to create person", http.StatusInternalServerError)
		log.Printf("Error creating person: %v", err)
		return
	}

	id, _ := result.LastInsertId()
	p.ID = int(id)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(p)
}

// Update an existing person
func updatePersonHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var p Person
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if p.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	result, err := db.Exec("UPDATE persons SET name = ?, email = ?, phone = ? WHERE id = ?", p.Name, p.Email, p.Phone, id)
	if err != nil {
		http.Error(w, "Failed to update person", http.StatusInternalServerError)
		log.Printf("Error updating person: %v", err)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Person not found", http.StatusNotFound)
		return
	}

	p.ID = id
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

// Delete a person
func deletePersonHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	result, err := db.Exec("DELETE FROM persons WHERE id = ?", id)
	if err != nil {
		http.Error(w, "Failed to delete person", http.StatusInternalServerError)
		log.Printf("Error deleting person: %v", err)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Person not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
