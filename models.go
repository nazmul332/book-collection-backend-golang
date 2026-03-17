package main

import "time"

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"-"` 
}

type Author struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Bio  string `json:"bio"`
}

type Genre struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Book struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	AuthorID    int       `json:"author_id"`
	GenreID     int       `json:"genre_id"`
	PublishedAt time.Time `json:"published_at"`
	Publisher   string    `json:"publisher"` 
	ISBN        string    `json:"isbn"`
	
	Author *Author `json:"author,omitempty"`
	Genre  *Genre  `json:"genre,omitempty"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type APIError struct {
	Message string `json:"message"`
}
