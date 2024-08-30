package handlers

import (
	"bytes"
	"encoding/json"
	"example/bookstore/database"
	"example/bookstore/models"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// SetupRouter sets up the Gin router with the necessary routes
func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/books", GetBooks)
	r.POST("/books", PostBooks)
	r.GET("/books/:id", GetBookByID)
	return r
}

// TestGetBooks tests the GET /books endpoint
func TestGetBooks(t *testing.T) {
	database.Connect() // Ensure the database is connected
	router := SetupRouter()

	req, err := http.NewRequest(http.MethodGet, "/books", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "Expected status OK, got %v", rr.Code)

	var books []models.Book
	err = json.Unmarshal(rr.Body.Bytes(), &books)
	if err != nil {
		t.Fatalf("Failed to unmarshal response body: %v", err)
	}

	assert.NotNil(t, books, "Expected books to be non-nil")
	assert.GreaterOrEqual(t, len(books), 0, "Expected at least 0 books in response")

	for _, book := range books {
		assert.NotZero(t, book.ID, "Expected book ID to be non-zero")
		assert.NotEmpty(t, book.Title, "Expected book title to be non-empty")
		assert.NotEmpty(t, book.Author, "Expected book author to be non-empty")
		assert.NotZero(t, book.Price, "Expected book price to be non-zero")
	}
}

// TestPostBooks tests the POST /books endpoint
func TestPostBooks(t *testing.T) {
	database.Connect() // Ensure the database is connected
	router := SetupRouter()

	book := models.Book{
		Title:  "Test Book",
		Author: "Test Author",
		Price:  19.99,
	}

	jsonBook, err := json.Marshal(book)
	if err != nil {
		t.Fatalf("Failed to marshal book: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, "/books", bytes.NewBuffer(jsonBook))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code, "Expected status Created, got %v", rr.Code)

	var createdBook models.Book
	err = json.Unmarshal(rr.Body.Bytes(), &createdBook)
	if err != nil {
		t.Fatalf("Failed to unmarshal response body: %v", err)
	}

	assert.NotNil(t, createdBook, "Expected createdBook to be non-nil")
	assert.Equal(t, book.Title, createdBook.Title, "Expected title to match")
	assert.Equal(t, book.Author, createdBook.Author, "Expected author to match")
	assert.Equal(t, book.Price, createdBook.Price, "Expected price to match")
}

// TestGetBookById tests the GET /books/:id endpoint
func TestGetBookById(t *testing.T) {
	database.Connect() // Ensure the database is connected
	router := SetupRouter()

	bookID := 1 // Replace with a valid ID for testing

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/books/%d", bookID), nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "Expected status OK, got %v", rr.Code)

	var book models.Book
	err = json.Unmarshal(rr.Body.Bytes(), &book)
	if err != nil {
		t.Fatalf("Failed to unmarshal response body: %v", err)
	}

	assert.NotNil(t, book, "Expected book to be non-nil")
	assert.Equal(t, int64(bookID), book.ID, "Expected book ID to match")
	assert.NotEmpty(t, book.Title, "Expected book title to be non-empty")
	assert.NotEmpty(t, book.Author, "Expected book author to be non-empty")
	assert.NotZero(t, book.Price, "Expected book price to be non-zero")
}
