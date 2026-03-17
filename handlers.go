package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

func sendError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(APIError{Message: message})
}

func sendJSON(w http.ResponseWriter, data interface{}, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(data)
}


func handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var user User
	err := db.QueryRow("SELECT id, username, password FROM users WHERE username = ?", req.Username).Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		sendError(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	if !CheckPasswordHash(req.Password, user.Password) {
		sendError(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	token, err := GenerateToken(user.Username)
	if err != nil {
		sendError(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	sendJSON(w, LoginResponse{Token: token}, http.StatusOK)
}


func handleBooks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getBooks(w, r)
	case http.MethodPost:
		createBook(w, r)
	default:
		sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleBookByID(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimSuffix(r.URL.Path, "/")
	idStr := strings.TrimPrefix(path, "/books/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		sendError(w, "Invalid book ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		getBook(w, r, id)
	case http.MethodPut:
		updateBook(w, r, id)
	case http.MethodDelete:
		deleteBook(w, r, id)
	default:
		sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}


func getBooks(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(`
		SELECT b.id, b.title, b.author_id, b.genre_id, b.published_at, b.publisher, b.isbn,
		       COALESCE(a.name, ''), COALESCE(a.bio, ''), 
		       COALESCE(g.name, ''), COALESCE(g.description, '')
		FROM books b
		LEFT JOIN authors a ON b.author_id = a.id
		LEFT JOIN genres g ON b.genre_id = g.id
	`)
	if err != nil {
		sendError(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	books := []Book{}
	for rows.Next() {
		var b Book
		var a Author
		var g Genre
		err := rows.Scan(&b.ID, &b.Title, &b.AuthorID, &b.GenreID, &b.PublishedAt, &b.Publisher, &b.ISBN,
			&a.Name, &a.Bio, &g.Name, &g.Description)
		if err != nil {
			sendError(w, "Error scanning books", http.StatusInternalServerError)
			return
		}
		a.ID = b.AuthorID
		g.ID = b.GenreID
		b.Author = &a
		b.Genre = &g
		books = append(books, b)
	}
	sendJSON(w, books, http.StatusOK)
}

func createBook(w http.ResponseWriter, r *http.Request) {
	var b Book
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		sendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	res, err := db.Exec("INSERT INTO books (title, author_id, genre_id, published_at, publisher, isbn) VALUES (?, ?, ?, ?, ?, ?)",
		b.Title, b.AuthorID, b.GenreID, b.PublishedAt, b.Publisher, b.ISBN)
	if err != nil {
		sendError(w, "Database error", http.StatusInternalServerError)
		return
	}

	id, _ := res.LastInsertId()
	b.ID = int(id)
	sendJSON(w, b, http.StatusCreated)
}

func getBook(w http.ResponseWriter, r *http.Request, id int) {
	var b Book
	var a Author
	var g Genre
	err := db.QueryRow(`
		SELECT b.id, b.title, b.author_id, b.genre_id, b.published_at, b.publisher, b.isbn,
		       COALESCE(a.name, ''), COALESCE(a.bio, ''), 
		       COALESCE(g.name, ''), COALESCE(g.description, '')
		FROM books b
		LEFT JOIN authors a ON b.author_id = a.id
		LEFT JOIN genres g ON b.genre_id = g.id
		WHERE b.id = ?`, id).Scan(&b.ID, &b.Title, &b.AuthorID, &b.GenreID, &b.PublishedAt, &b.Publisher, &b.ISBN,
		&a.Name, &a.Bio, &g.Name, &g.Description)

	if err != nil {
		sendError(w, "Book not found or database error", http.StatusNotFound)
		return
	}
	a.ID = b.AuthorID
	g.ID = b.GenreID
	b.Author = &a
	b.Genre = &g
	sendJSON(w, b, http.StatusOK)
}

func updateBook(w http.ResponseWriter, r *http.Request, id int) {
	var b Book
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		sendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	_, err := db.Exec("UPDATE books SET title = ?, author_id = ?, genre_id = ?, published_at = ?, publisher = ?, isbn = ? WHERE id = ?",
		b.Title, b.AuthorID, b.GenreID, b.PublishedAt, b.Publisher, b.ISBN, id)
	if err != nil {
		sendError(w, "Database error", http.StatusInternalServerError)
		return
	}

	b.ID = id
	sendJSON(w, b, http.StatusOK)
}

func deleteBook(w http.ResponseWriter, r *http.Request, id int) {
	_, err := db.Exec("DELETE FROM books WHERE id = ?", id)
	if err != nil {
		sendError(w, "Database error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}


func handleAuthors(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		rows, _ := db.Query("SELECT id, name, bio FROM authors")
		defer rows.Close()
		authors := []Author{}
		for rows.Next() {
			var a Author
			rows.Scan(&a.ID, &a.Name, &a.Bio)
			authors = append(authors, a)
		}
		sendJSON(w, authors, http.StatusOK)
	case http.MethodPost:
		var a Author
		json.NewDecoder(r.Body).Decode(&a)
		res, _ := db.Exec("INSERT INTO authors (name, bio) VALUES (?, ?)", a.Name, a.Bio)
		id, _ := res.LastInsertId()
		a.ID = int(id)
		sendJSON(w, a, http.StatusCreated)
	default:
		sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleAuthorByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/authors/")
	id, _ := strconv.Atoi(idStr)
	switch r.Method {
	case http.MethodGet:
		var a Author
		err := db.QueryRow("SELECT id, name, bio FROM authors WHERE id = ?", id).Scan(&a.ID, &a.Name, &a.Bio)
		if err != nil {
			sendError(w, "Author not found", http.StatusNotFound)
			return
		}
		sendJSON(w, a, http.StatusOK)
	case http.MethodPut:
		var a Author
		json.NewDecoder(r.Body).Decode(&a)
		db.Exec("UPDATE authors SET name = ?, bio = ? WHERE id = ?", a.Name, a.Bio, id)
		a.ID = id
		sendJSON(w, a, http.StatusOK)
	case http.MethodDelete:
		db.Exec("DELETE FROM authors WHERE id = ?", id)
		w.WriteHeader(http.StatusNoContent)
	}
}


func handleGenres(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		rows, _ := db.Query("SELECT id, name, description FROM genres")
		defer rows.Close()
		genres := []Genre{}
		for rows.Next() {
			var g Genre
			rows.Scan(&g.ID, &g.Name, &g.Description)
			genres = append(genres, g)
		}
		sendJSON(w, genres, http.StatusOK)
	case http.MethodPost:
		var g Genre
		json.NewDecoder(r.Body).Decode(&g)
		res, _ := db.Exec("INSERT INTO genres (name, description) VALUES (?, ?)", g.Name, g.Description)
		id, _ := res.LastInsertId()
		g.ID = int(id)
		sendJSON(w, g, http.StatusCreated)
	default:
		sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleGenreByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/genres/")
	id, _ := strconv.Atoi(idStr)
	switch r.Method {
	case http.MethodGet:
		var g Genre
		err := db.QueryRow("SELECT id, name, description FROM genres WHERE id = ?", id).Scan(&g.ID, &g.Name, &g.Description)
		if err != nil {
			sendError(w, "Genre not found", http.StatusNotFound)
			return
		}
		sendJSON(w, g, http.StatusOK)
	case http.MethodPut:
		var g Genre
		json.NewDecoder(r.Body).Decode(&g)
		db.Exec("UPDATE genres SET name = ?, description = ? WHERE id = ?", g.Name, g.Description, id)
		g.ID = id
		sendJSON(w, g, http.StatusOK)
	case http.MethodDelete:
		db.Exec("DELETE FROM genres WHERE id = ?", id)
		w.WriteHeader(http.StatusNoContent)
	}
}
