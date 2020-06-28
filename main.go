package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// Item is a struct that represents a single item
type Item struct {
	Title string `json:"title"`
}

var items []Item = []Item{{"Foo"}, {"Bar"}, {"Goo"}}

func main() {

	// App Constance
	const buildPath string = "build/"
	const port int = 8080

	// Router
	router := mux.NewRouter()

	// API Server Routes
	router.HandleFunc("/items", getAllItems).Methods("GET")
	router.HandleFunc("/items", addNewItem).Methods("POST")
	router.HandleFunc("/clearAll", clearAllItems).Methods("DELETE")

	// File Server
	router.PathPrefix("/").Handler(http.FileServer(http.Dir(buildPath)))

	// CORS Headers
	cors := handlers.CORS(
		handlers.AllowedHeaders([]string{"content-type"}),
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowCredentials(),
		handlers.AllowedMethods([]string{"GET", "POST", "DELETE"}),
	)

	// Run Server
	if err := http.ListenAndServe(":"+strconv.Itoa(port), cors(router)); err != nil {
		log.Fatal("ListenAndServe error: ", err)
	}
}

// Handlers

func getAllItems(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get All Items Requested")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

func addNewItem(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Add New Items Requested")

	var newItem Item
	json.NewDecoder(r.Body).Decode(&newItem)
	fmt.Println("New Items = " + newItem.Title)

	items = append(items, newItem)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

func clearAllItems(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Clearing all Items")

	items = []Item{}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}
