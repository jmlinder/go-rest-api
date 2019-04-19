package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Book Struct (Model)
type Book struct {
	ID     string  `json:"id"`
	Isbn   string  `json:"isbn"`
	Title  string  `json:"title"`
	Author *Author `json:"author"`
}

// Author Struct
type Author struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

// Init books var as a slice Book struct
var books []Book

// Get All Books // ~app.get('', (req, res))
func getBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

// Create Book
func createBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var book Book
	_ = json.NewDecoder(r.Body).Decode(&book)
	book.ID = strconv.Itoa(rand.Intn(1000000)) // Mock ID - not safe
	books = append(books, book)
	json.NewEncoder(w).Encode(book)
}

// Read Book
func getBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r) // GET params
	// Loop through books and find with id
	for _, item := range books {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Book{})
}

//Update Book @todo - fix; only deletes the item
func updateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range books {
		if item.ID == params["id"] {
			books = append(books[:index], books[index+1:]...) // Delete syntax
			var book Book
			_ = json.NewDecoder(r.Body).Decode(&book)
			book.ID = params["id"]
			books = append(books, book)
			json.NewEncoder(w).Encode(book)
			return
		}

	}
}

// Delete Book
func deleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range books {
		if item.ID == params["id"] {
			books = append(books[:index], books[index+1:]...) // Delete syntax
			break
		}
	}
	json.NewEncoder(w).Encode(books)
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the homepage!")

}

func main() {
	// Init Router
	r := mux.NewRouter()
	// Mock Data @todo - implement DB
	books = append(books, Book{ID: "1", Isbn: "49750", Title: "The Art of Being Cool", Author: &Author{Firstname: "Bender", Lastname: "Rodriguez"}})
	books = append(books, Book{ID: "2", Isbn: "36433", Title: "How to Lose Friends and Alienate People", Author: &Author{Firstname: "Reggie", Lastname: "Jackson"}})
	books = append(books, Book{ID: "3", Isbn: "68009", Title: "Book Eight Point Three", Author: &Author{Firstname: "Franz", Lastname: "Liszt"}})
	books = append(books, Book{ID: "4", Isbn: "16378", Title: "Dial M for Magrudergrind", Author: &Author{Firstname: "John", Lastname: "Jakobzon"}})

	// Route Handlers / Endpoints
	r.HandleFunc("/", homePage)
	r.HandleFunc("/api/books", getBooks).Methods("GET")     // Index
	r.HandleFunc("/api/books/{id}", getBook).Methods("GET") // Show
	r.HandleFunc("/api/books", createBook).Methods("POST")  // Create
	// r.HandleFunc("/api/books/{id}", updateBook).Methods("PUT")    // Update
	r.HandleFunc("/api/books/{id}", deleteBook).Methods("DELETE") // Delete
	log.Fatal(http.ListenAndServe(":8000", r))
}
