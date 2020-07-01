package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// Item is a struct that represents a single item
type Item struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

var items []Item = []Item{{Title: "Foo"}, {Title: "Bar"}, {Title: "Goo"}}

func main() {

	// Assign random ID to default Items
	for index := range items {
		items[index].ID = getUUID()
	}

	// App Constance
	const buildPath string = "build/"
	const port int = 8080

	// Router
	router := mux.NewRouter()

	// API Server Routes
	router.HandleFunc("/items", getAllItems).Methods("GET")
	router.HandleFunc("/items", addNewItem).Methods("POST")
	router.HandleFunc("/removeItem", removeItem).Methods("POST")
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

	// Register Middleware
	router.Use(loggingMiddleware)

	// Run Server
	srv := &http.Server{
		Handler: cors(router),
		Addr:    ":" + strconv.Itoa(port),
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal("ListenAndServe error: ", err)
	}
}

// Handlers

func getAllItems(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

func addNewItem(w http.ResponseWriter, r *http.Request) {

	var newItem Item
	json.NewDecoder(r.Body).Decode(&newItem)
	newItem.ID = getUUID()
	fmt.Println("New Items = " + newItem.Title + ", ID = " + newItem.ID)

	items = append(items, newItem)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

func removeItem(w http.ResponseWriter, r *http.Request) {

	var itemToRemove Item
	json.NewDecoder(r.Body).Decode(&itemToRemove)
	fmt.Println("Item to Remove = " + itemToRemove.ID)

	for index, item := range items {
		if item.ID == itemToRemove.ID {
			items[index] = items[len(items)-1]
			items = items[:len(items)-1]
			break
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

func clearAllItems(w http.ResponseWriter, r *http.Request) {

	items = []Item{}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

// Helper Methods

func getUUID() string {
	uuid, err := uuid.NewRandom()

	if err != nil {
		log.Fatal(err)
	}
	id := uuid.String()

	return id
}

// Middleware Methods

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mes := "URI: " + r.RequestURI + ", RemoteAddr: " + r.RemoteAddr + ", Method:" + r.Method
		log.Println(mes)
		next.ServeHTTP(w, r)
	})
}
