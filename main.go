package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

const (
	POST_METHOD     = "POST"
	GET_METHOD      = "GET"
	SUCCESS_MESSAGE = "data successfuly"
	SUCCESS_CODE    = 200
	PUT_METHOD      = "PUT"
	DELETE_METHOD   = "DELETE"
)

type Books struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Author      string    `json:"author"`
	PublishDate time.Time `json:"publish_date"`
	Rating      int64     `json:"rating"`
}

type Response struct {
	Message string
	Status  int
}

func main() {
	log.Println("Server started on: http://127.0.0.1:8080")
	router := mux.NewRouter()
	router.HandleFunc("/", createBook).Methods(POST_METHOD)
	router.HandleFunc("/books", getBooks).Methods(GET_METHOD)
	router.HandleFunc("/books/{id}", updateBook).Methods(PUT_METHOD)
	router.HandleFunc("/books/{id}", deleteBook).Methods(DELETE_METHOD)
	http.ListenAndServe(":8080", router)
}

func dbConn() (db *sql.DB) {
	dbSourceHost := "YOUR_HOST"
	dbSourcePort := "YOUR_PORT"
	dbSourceUser := "YOUR_USER"
	dbSourceName := "YOUR_DB_NAME"
	dbSourceSslMode := "SSL_MODE"
	dbSourcePassword := "YOUR_PASSWORD"

	dbDriver := "postgres"
	dbSource := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s sslmode=%s password=%s",
		dbSourceHost,
		dbSourcePort,
		dbSourceUser,
		dbSourceName,
		dbSourceSslMode,
		dbSourcePassword,
	)

	db, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func createBook(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	defer db.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	var books Books

	json.Unmarshal(body, &books)

	title := books.Title
	author := books.Author
	rating := books.Rating

	stmt, err := db.Prepare("INSERT INTO Books(title, author, rating) VALUES($1, $2, $3)")
	if err != nil {
		log.Fatal(err)
	}
	if _, err := stmt.Exec(title, author, rating); err != nil {
		log.Fatal(err)
	}

	var response Response
	response.Status = SUCCESS_CODE
	response.Message = fmt.Sprintf("Create %s", SUCCESS_MESSAGE)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
	log.Println(response.Message)
}

func getBooks(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	defer db.Close()

	rows, err := db.Query("SELECT * FROM books")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	books := make([]Books, 0)

	for rows.Next() {
		book := Books{}
		err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.PublishDate, &book.Rating)
		if err != nil {
			log.Fatal(err)
		}
		books = append(books, book)
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

func updateBook(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	defer db.Close()

	params := mux.Vars(r)

	stmt, err := db.Prepare("UPDATE books SET title = $1, author = $2, rating = $3 WHERE id = $4")
	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	var books Books

	json.Unmarshal(body, &books)

	title := books.Title
	author := books.Author
	rating := books.Rating

	if _, err = stmt.Exec(title, author, rating, params["id"]); err != nil {
		log.Fatal(err)
	}

	var response Response
	response.Status = SUCCESS_CODE
	response.Message = fmt.Sprintf("Update %s", SUCCESS_MESSAGE)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
	log.Println(response.Message)
}

func deleteBook(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	defer db.Close()

	params := mux.Vars(r)

	stmt, err := db.Prepare("DELETE FROM books WHERE id = $1")
	if err != nil {
		log.Fatal(err)
	}
	if _, err := stmt.Exec(params["id"]); err != nil {
		log.Fatal(err)
	}

	var response Response
	response.Status = SUCCESS_CODE
	response.Message = fmt.Sprintf("DELETE %s", SUCCESS_MESSAGE)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
	log.Println(response.Message)
}
