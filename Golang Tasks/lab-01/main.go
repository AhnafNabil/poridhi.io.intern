package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Book struct {
	ID     primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Title  string             `json:"title,omitempty" bson:"title,omitempty"`
	Author string             `json:"author,omitempty" bson:"author,omitempty"`
	Year   int                `json:"year,omitempty" bson:"year,omitempty"`
}

var client *mongo.Client
var booksCollection *mongo.Collection

func main() {
	// Set up MongoDB connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	var err error
	client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	// Get a handle for the books collection
	booksCollection = client.Database("library").Collection("books")

	// Set up the router
	r := mux.NewRouter()

	// Define routes
	r.HandleFunc("/books", getBooks).Methods("GET")
	r.HandleFunc("/books/{id}", getBook).Methods("GET")
	r.HandleFunc("/books", createBook).Methods("POST")
	r.HandleFunc("/books/{id}", updateBook).Methods("PUT")
	r.HandleFunc("/books/{id}", deleteBook).Methods("DELETE")

	// Start the server
	fmt.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func getBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var books []Book

	cur, err := booksCollection.Find(context.Background(), bson.M{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cur.Close(context.Background())

	for cur.Next(context.Background()) {
		var book Book
		err := cur.Decode(&book)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		books = append(books, book)
	}

	json.NewEncoder(w).Encode(books)
}

func getBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	id, _ := primitive.ObjectIDFromHex(params["id"])

	var book Book
	err := booksCollection.FindOne(context.Background(), Book{ID: id}).Decode(&book)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "Book not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	json.NewEncoder(w).Encode(book)
}

func createBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var book Book
	_ = json.NewDecoder(r.Body).Decode(&book)

	result, err := booksCollection.InsertOne(context.Background(), book)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	book.ID = result.InsertedID.(primitive.ObjectID)
	json.NewEncoder(w).Encode(book)
}

func updateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	id, _ := primitive.ObjectIDFromHex(params["id"])

	var book Book
	_ = json.NewDecoder(r.Body).Decode(&book)

	filter := bson.M{"_id": id}
	update := bson.M{"$set": book}

	err := booksCollection.FindOneAndUpdate(context.Background(), filter, update).Decode(&book)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "Book not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	book.ID = id
	json.NewEncoder(w).Encode(book)
}

func deleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	id, _ := primitive.ObjectIDFromHex(params["id"])

	result, err := booksCollection.DeleteOne(context.Background(), Book{ID: id})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if result.DeletedCount == 0 {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"result": "Book deleted successfully"})
}