package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/jmoiron/sqlx"
	"github.com/google/uuid"
)

type ToDoList struct {
	ID   uuid.UUID `json:"id" db:"id"`
	Name string    `json:"name" db:"name"`
}

type ToDoItem struct {
	ID          uuid.UUID `json:"id" db:"id"`
	ListID      uuid.UUID `json:"list_id" db:"list_id"`
	Description string    `json:"description" db:"description"`
	IsChecked   bool      `json:"is_checked" db:"is_checked"`
}

var db *sqlx.DB

func initDB() {
	// Load environment variables for database connection
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))

	var err error
	db, err = sqlx.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Unable to connect to the database:", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Unable to reach the database:", err)
	}
}

func createToDoList(w http.ResponseWriter, r *http.Request) {
	var list ToDoList
	json.NewDecoder(r.Body).Decode(&list)

	list.ID = uuid.New()
	_, err := db.Exec("INSERT INTO todo_lists (id, name) VALUES ($1, $2)", list.ID, list.Name)
	if err != nil {
		http.Error(w, "Failed to create list", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}

func addToDoItem(w http.ResponseWriter, r *http.Request) {
	var item ToDoItem
	json.NewDecoder(r.Body).Decode(&item)

	// Check if the list exists
	var list ToDoList
	err := db.Get(&list, "SELECT id FROM todo_lists WHERE id = $1", item.ListID)
	if err == sql.ErrNoRows {
		http.Error(w, "List not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Failed to fetch list", http.StatusInternalServerError)
		return
	}

	// check if the same description already exists in the list
	var count int
	err = db.Get(&count, "SELECT COUNT(*) FROM todo_items WHERE list_id = $1 AND description = $2", item.ListID, item.Description)
	if err != nil {
		http.Error(w, "Failed to check for duplicate item", http.StatusInternalServerError)
		return
	}

	if count > 0 {
		http.Error(w, "Item already exists in the list", http.StatusBadRequest)
		return
	}

	item.ID = uuid.New()
	_, err = db.Exec("INSERT INTO todo_items (id, list_id, description) VALUES ($1, $2, $3)", item.ID, item.ListID, item.Description)
	if err != nil {
		http.Error(w, "Failed to add item", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}

func getToDoList(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	listID, err := uuid.Parse(params["id"])
	if err != nil {
		http.Error(w, "Invalid list ID", http.StatusBadRequest)
		return
	}

	var list ToDoList
	err = db.Get(&list, "SELECT id, name FROM todo_lists WHERE id = $1", listID)
	if err == sql.ErrNoRows {
		http.Error(w, "List not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Failed to fetch list", http.StatusInternalServerError)
		return
	}

	var items []ToDoItem
	err = db.Select(&items, "SELECT id, list_id, description, is_checked FROM todo_items WHERE list_id = $1", listID)
	if err != nil {
		http.Error(w, "Failed to fetch items", http.StatusInternalServerError)
		return
	}

	response := struct {
		List  ToDoList  `json:"list"`
		Items []ToDoItem `json:"items"`
	}{
		List:  list,
		Items: items,
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(response)
}

func updateToDoItem(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	itemID, err := uuid.Parse(params["id"])
	if err != nil {
		http.Error(w, "Invalid item ID", http.StatusBadRequest)
		return
	}

	var item ToDoItem
	json.NewDecoder(r.Body).Decode(&item)

	_, err = db.Exec("UPDATE todo_items SET description = $1, is_checked = $2 WHERE id = $3", item.Description, item.IsChecked, itemID)
	if err != nil {
		http.Error(w, "Failed to update item", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func deleteToDoItem(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	itemID, err := uuid.Parse(params["id"])
	if err != nil {
		http.Error(w, "Invalid item ID", http.StatusBadRequest)
		return
	}

	_, err = db.Exec("DELETE FROM todo_items WHERE id = $1", itemID)
	if err != nil {
		http.Error(w, "Failed to delete item", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func main() {
	initDB()
	defer db.Close()

	r := mux.NewRouter()

	r.HandleFunc("/lists", createToDoList).Methods("POST")
	r.HandleFunc("/lists/{id}", getToDoList).Methods("GET")
	r.HandleFunc("/items", addToDoItem).Methods("POST")
	r.HandleFunc("/items/{id}", updateToDoItem).Methods("PUT")
	r.HandleFunc("/items/{id}", deleteToDoItem).Methods("DELETE")

	log.Println("API is running on port 8080")
	http.ListenAndServe(":8080", r)
}
