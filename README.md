# book-collection-backend-golang

A simple RESTful API built with Go (Golang) for managing a book collection.
This project includes authentication using JWT, full CRUD operations, SQLite database, and unit tests.

## Features

* JWT Authentication (Login system)
* Manage Books (CRUD)
* Manage Authors (CRUD)
* Manage Genres (CRUD)
* SQLite Database
* Unit Testing
* Docker Support

## Tech Stack

* Go (Golang)
* SQLite
* JWT (github.com/golang-jwt/jwt/v5)
* bcrypt (password hashing)
* net/http (standard library)

## Project Structure

├── main.go
├── handlers.go
├── models.go
├── database.go
├── auth.go
├── main_test.go
├── Dockerfile
├── go.mod


## Setup & Run

### 1. Clone the repository

https://github.com/nazmul332/book-collection-backend-golang.git

### 2. Run the project

go run main.go

Server will start at:

http://localhost:8080

## Authentication

### Default user (auto seeded)

username: admin
password: admin123

### Login

POST /login

Request:

json
{
  "username": "admin",
  "password": "admin123"
}

Response:

json
{
  "token": "your-jwt-token"
}

 Use this token in headers:

Authorization: Bearer <token>

## API Endpoints

### Books

| Method | Endpoint    | Description    |
| ------ | ----------- | -------------- |
| GET    | /books      | Get all books  |
| POST   | /books      | Create a book  |
| GET    | /books/{id} | Get book by ID |
| PUT    | /books/{id} | Update book    |
| DELETE | /books/{id} | Delete book    |


### Authors

| Method | Endpoint      | Description     |
| ------ | ------------- | --------------- |
| GET    | /authors      | Get all authors |
| POST   | /authors      | Create author   |
| GET    | /authors/{id} | Get author      |
| PUT    | /authors/{id} | Update author   |
| DELETE | /authors/{id} | Delete author   |


### Genres

| Method | Endpoint     | Description    |
| ------ | ------------ | -------------- |
| GET    | /genres      | Get all genres |
| POST   | /genres      | Create genre   |
| GET    | /genres/{id} | Get genre      |
| PUT    | /genres/{id} | Update genre   |
| DELETE | /genres/{id} | Delete genre   |


## Running Tests

go test ./...

## Database

* Uses SQLite (`book_collection.db`)
* Tables:

  * users
  * books
  * authors
  * genres

## Notes

* JWT secret key is hardcoded → change it in production
* Passwords are hashed using bcrypt
* API is protected (except `/login`)


## Future Improvements

* Add pagination
* Add search/filter
* Add Swagger documentation
* Improve error handling
* Use environment variables

## Author

Nazmul Hasan


## License

This project is open-source and free to use.
