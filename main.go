// package main

// import (
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"os"
// 	"time"

// 	"github.com/joho/godotenv"
// 	"go.mongodb.org/mongo-driver/bson"
// 	"go.mongodb.org/mongo-driver/mongo"
// 	"go.mongodb.org/mongo-driver/mongo/options"
// )

// type Post struct {
// 	Title   string `json:"title"`
// 	Content string `json:"content"`
// }

// var postCollection *mongo.Collection

// // Connect to MongoDB
// func connectDB() *mongo.Client {
// 	// Load environment variables
// 	err := godotenv.Load()
// 	if err != nil {
// 		log.Fatal("Error loading .env file")
// 	}

// 	// Get the MongoDB URI from the environment
// 	mongoURI := os.Getenv("MONGO_URI")
// 	if mongoURI == "" {
// 		log.Fatal("MONGO_URI not set in environment")
// 	}

// 	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()

// 	err = client.Connect(ctx)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	return client
// }

// // POST handler
// func createPost(w http.ResponseWriter, r *http.Request) {
// 	var newPost Post
// 	err := json.NewDecoder(r.Body).Decode(&newPost)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	// Insert post into MongoDB
// 	_, err = postCollection.InsertOne(context.TODO(), newPost)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	w.WriteHeader(http.StatusCreated)
// 	fmt.Fprintf(w, "Post created")
// }

// // GET handler
// func getPosts(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")

// 	// Find all posts in MongoDB
// 	cur, err := postCollection.Find(context.TODO(), bson.D{})
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	defer cur.Close(context.TODO())

// 	var posts []Post
// 	for cur.Next(context.TODO()) {
// 		var post Post
// 		err := cur.Decode(&post)
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}
// 		posts = append(posts, post)
// 	}

// 	if err := json.NewEncoder(w).Encode(posts); err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 	}
// }

// func main() {
// 	// Load environment variables
// 	err := godotenv.Load()
// 	if err != nil {
// 		log.Fatal("Error loading .env file")
// 	}

// 	// Connect to the MongoDB database
// 	client := connectDB()
// 	defer client.Disconnect(context.TODO())

// 	// Get the database and collection name from the environment
// 	dbName := os.Getenv("DB_NAME")
// 	collectionName := os.Getenv("COLLECTION_NAME")

// 	// Define the database and collection
// 	postCollection = client.Database(dbName).Collection(collectionName)

// 	http.HandleFunc("/posts", func(w http.ResponseWriter, r *http.Request) {
// 		switch r.Method {
// 		case http.MethodPost:
// 			createPost(w, r)
// 		case http.MethodGet:
// 			getPosts(w, r)
// 		default:
// 			w.WriteHeader(http.StatusMethodNotAllowed)
// 		}
// 	})

//		fmt.Println("Server running on port 8080...")
//		log.Fatal(http.ListenAndServe(":8080", nil))
//	}
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var messageHistory []string
var historyMux sync.Mutex

// Reverse a string
func reverseMessage(msg string) string {
	runes := []rune(msg)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// Set CORS headers
func enableCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

// WebSocket handler
func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	if r.Method == http.MethodOptions {
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade to WebSocket:", err)
		return
	}

	defer conn.Close()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Failed to read message:", err)
			break
		}

		message := string(msg)
		fmt.Println("Received:", message)

		reversed := reverseMessage(message)

		historyMux.Lock()
		messageHistory = append(messageHistory, message)
		if len(messageHistory) > 5 {
			messageHistory = messageHistory[1:]
		}
		historyMux.Unlock()

		err = conn.WriteMessage(websocket.TextMessage, []byte(reversed))
		if err != nil {
			log.Println("Failed to send message:", err)
			break
		}
	}
}

func getLastFiveMessages(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	if r.Method == http.MethodOptions {
		return
	}

	historyMux.Lock()
	defer historyMux.Unlock()

	count := len(messageHistory)
	if count > 5 {
		count = 5
	}

	lastFive := messageHistory[len(messageHistory)-count:]

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(lastFive)
}

func main() {
	http.HandleFunc("/ws", handleWebSocket)
	http.HandleFunc("/history", getLastFiveMessages)

	fmt.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
