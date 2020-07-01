package main

import (
	"encoding/json"
	"fmt"

	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/fatih/color"
	"github.com/google/uuid"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

// Item is a struct that represents a single item
type Item struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

var items []Item = []Item{{Title: "Foo"}, {Title: "Bar"}, {Title: "Goo"}}

// WebSocket
var wsConnections = []*websocket.Conn{}
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

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

	// WebSocket Endpoint
	router.HandleFunc("/ws", wsEndpoint)

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

	color.New(color.BgHiGreen, color.FgBlack, color.Bold).
		Println("Server is Running on PORT " + strconv.Itoa(port))

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
	items = append(items, newItem)

	sendUpdateComand()
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

	sendUpdateComand()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

func clearAllItems(w http.ResponseWriter, r *http.Request) {

	items = []Item{}

	sendUpdateComand()
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

func sendUpdateComand() {
	for _, ws := range wsConnections {
		// color.New(color.BgHiMagenta).Println("Sending Update Message to: " + ws.RemoteAddr().String())
		ws.WriteMessage(1, []byte("itemsUpdated"))
	}
}

// Middleware Methods

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mes := "URI: " + r.RequestURI + ", RemoteAddr: " + r.RemoteAddr + ", Method:" + r.Method
		log.Println(mes)
		next.ServeHTTP(w, r)
	})
}

// WebSocket

func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	log.Println("Client Attempting to connect...")
	fmt.Println(r.Host)

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	wsConnections = append(wsConnections, ws)

	log.Println("Client Successfully Connected...")

	// reader(ws)
}

func reader(conn *websocket.Conn) {
	for {
		// read in a message
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		// print out that message for clarity
		fmt.Println(string(p))
		fmt.Println(messageType)

		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}

	}
}
