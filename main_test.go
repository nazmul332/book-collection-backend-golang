package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func setupTests() {
	err := InitDB(":memory:")
	if err != nil {
		panic(err)
	}
	seedUser("testuser", "password123")
}

func TestLogin(t *testing.T) {
	setupTests()

	loginReq := LoginRequest{Username: "testuser", Password: "password123"}
	body, _ := json.Marshal(loginReq)
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()
	
	handleLogin(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var resp LoginResponse
	json.Unmarshal(rr.Body.Bytes(), &resp)
	if resp.Token == "" {
		t.Error("expected token in login response, got empty string")
	}

	invalidReq := LoginRequest{Username: "testuser", Password: "wrongpassword"}
	body, _ = json.Marshal(invalidReq)
	req, _ = http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	rr = httptest.NewRecorder()
	
	handleLogin(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code for invalid login: got %v want %v", status, http.StatusUnauthorized)
	}
}

func TestBooksCRUD(t *testing.T) {
	setupTests()
	
	token, _ := GenerateToken("testuser")
	authHeader := "Bearer " + token

	newBook := Book{
		Title:     "The Go Programming Language",
		AuthorID:  1,
		GenreID:   1,
		Publisher: "Addison-Wesley",
		ISBN:      "978-0134190440",
	}
	body, _ := json.Marshal(newBook)
	req, _ := http.NewRequest("POST", "/books", bytes.NewBuffer(body))
	req.Header.Set("Authorization", authHeader)
	rr := httptest.NewRecorder()
	
	AuthMiddleware(handleBooks).ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	var createdBook Book
	json.Unmarshal(rr.Body.Bytes(), &createdBook)
	if createdBook.ID == 0 {
		t.Error("expected book ID to be set after creation")
	}

	req, _ = http.NewRequest("GET", "/books", nil)
	req.Header.Set("Authorization", authHeader)
	rr = httptest.NewRecorder()
	AuthMiddleware(handleBooks).ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var books []Book
	json.Unmarshal(rr.Body.Bytes(), &books)
	if len(books) == 0 {
		t.Error("expected at least one book in list")
	}

	createdBook.Title = "Updated Book Title"
	body, _ = json.Marshal(createdBook)
	req, _ = http.NewRequest("PUT", "/books/1", bytes.NewBuffer(body))
	req.Header.Set("Authorization", authHeader)
	rr = httptest.NewRecorder()
	AuthMiddleware(handleBookByID).ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	req, _ = http.NewRequest("DELETE", "/books/1", nil)
	req.Header.Set("Authorization", authHeader)
	rr = httptest.NewRecorder()
	AuthMiddleware(handleBookByID).ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNoContent)
	}
}

func TestUnauthorizedAccess(t *testing.T) {
	setupTests()

	req, _ := http.NewRequest("GET", "/books", nil)
	rr := httptest.NewRecorder()
	AuthMiddleware(handleBooks).ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code for unauthorized access: got %v want %v", status, http.StatusUnauthorized)
	}
}

func TestAuthorsCRUD(t *testing.T) {
	setupTests()
	token, _ := GenerateToken("testuser")
	authHeader := "Bearer " + token

	author := Author{Name: "J.K. Rowling", Bio: "British author"}
	body, _ := json.Marshal(author)
	req, _ := http.NewRequest("POST", "/authors", bytes.NewBuffer(body))
	req.Header.Set("Authorization", authHeader)
	rr := httptest.NewRecorder()
	AuthMiddleware(handleAuthors).ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("wrong status: got %v", rr.Code)
	}

	req, _ = http.NewRequest("GET", "/authors/1", nil)
	req.Header.Set("Authorization", authHeader)
	rr = httptest.NewRecorder()
	AuthMiddleware(handleAuthorByID).ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("wrong status: got %v", rr.Code)
	}

	author.Name = "Updated Name"
	body, _ = json.Marshal(author)
	req, _ = http.NewRequest("PUT", "/authors/1", bytes.NewBuffer(body))
	req.Header.Set("Authorization", authHeader)
	rr = httptest.NewRecorder()
	AuthMiddleware(handleAuthorByID).ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("wrong status: got %v", rr.Code)
	}

	req, _ = http.NewRequest("DELETE", "/authors/1", nil)
	req.Header.Set("Authorization", authHeader)
	rr = httptest.NewRecorder()
	AuthMiddleware(handleAuthorByID).ServeHTTP(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Errorf("wrong status: got %v", rr.Code)
	}
}

func TestGenresCRUD(t *testing.T) {
	setupTests()
	token, _ := GenerateToken("testuser")
	authHeader := "Bearer " + token

	genre := Genre{Name: "Fantasy", Description: "Magic and wonder"}
	body, _ := json.Marshal(genre)
	req, _ := http.NewRequest("POST", "/genres", bytes.NewBuffer(body))
	req.Header.Set("Authorization", authHeader)
	rr := httptest.NewRecorder()
	AuthMiddleware(handleGenres).ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("wrong status: got %v", rr.Code)
	}

	genre.Name = "Sci-Fi"
	body, _ = json.Marshal(genre)
	req, _ = http.NewRequest("PUT", "/genres/1", bytes.NewBuffer(body))
	req.Header.Set("Authorization", authHeader)
	rr = httptest.NewRecorder()
	AuthMiddleware(handleGenreByID).ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("wrong status: got %v", rr.Code)
	}

	req, _ = http.NewRequest("DELETE", "/genres/1", nil)
	req.Header.Set("Authorization", authHeader)
	rr = httptest.NewRecorder()
	AuthMiddleware(handleGenreByID).ServeHTTP(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Errorf("wrong status: got %v", rr.Code)
	}
}

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}
