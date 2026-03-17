package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

func main() {
	err := InitDB("book_collection.db")
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	seedUser("admin", "admin123")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "<h1>Book Collection API is Running!</h1><p>Please use Postman to test the endpoints.</p>")
	})

	http.HandleFunc("/login", handleLogin)

	http.HandleFunc("/books", AuthMiddleware(handleBooks))
	http.HandleFunc("/books/", AuthMiddleware(handleBookByID))
	http.HandleFunc("/authors", AuthMiddleware(handleAuthors))
	http.HandleFunc("/authors/", AuthMiddleware(handleAuthorByID))
	http.HandleFunc("/genres", AuthMiddleware(handleGenres))
	http.HandleFunc("/genres/", AuthMiddleware(handleGenreByID))

	log.Println("Server starting on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func seedUser(username, password string) {
	var exists bool
	err := db.QueryRow("SELECT exists(SELECT 1 FROM users WHERE username=?)", username).Scan(&exists)
	if err != nil {
		log.Printf("Error checking user: %v", err)
		return
	}

	if !exists {
		hashed, _ := HashPassword(password)
		_, err = db.Exec("INSERT INTO users (username, password) VALUES (?, ?)", username, hashed)
		if err != nil {
			log.Printf("Error seeding user: %v", err)
		} else {
			log.Printf("Seeded user: %s", username)
		}
	}
}

var _ = strings.Contains
