package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoURL = "mongodb://localhost:27017"
var dbName = "gomongoapi"

type App struct {
	Router *mux.Router
	DB     *mongo.Database
}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

func (a *App) Initialize(_user, _password string) {
	fmt.Println("Starting the application....")

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	a.DB, _ = a.configDB(ctx)
	fmt.Println("Connected to MongoDB!")

	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

func (a *App) configDB(ctx context.Context) (*mongo.Database, error) {
	clientOptions := options.Client().ApplyURI(mongoURL)
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		return nil, fmt.Errorf("Couldn't connect to mongo: %v", err)
	}
	err = client.Connect(ctx)
	if err != nil {
		return nil, fmt.Errorf("Mongo client couldn't connect with background context: %v", err)
	}
	return client.Database(dbName), nil
}

// routing
func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/books", a.getBooks).Methods("GET")
	a.Router.HandleFunc("/book/{id}", a.getBook).Methods("GET")
	a.Router.HandleFunc("/book", a.ceateBook).Methods("POST")
	a.Router.HandleFunc("/book/{id}", a.updateBook).Methods("PUT")
	a.Router.HandleFunc("/book/{id}", a.deleteBook).Methods("DELETE")
}

func (a *App) getBooks(w http.ResponseWriter, r *http.Request) {
	count, _ := strconv.Atoi(r.FormValue("count"))
	start, _ := strconv.Atoi(r.FormValue("start"))

	if count > 10 || count < 1 {
		count = 10
	}
	if start < 0 {
		start = 0
	}

	bs, err := getBooks(a.DB, start, count)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, bs)
}

func (a *App) getBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := primitive.ObjectIDFromHex(vars["id"])

	b := Book{ID: id}
	if err := b.getBook(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, b)
}

func (a *App) ceateBook(w http.ResponseWriter, r *http.Request) {
}

func (a *App) updateBook(w http.ResponseWriter, r *http.Request) {
}

func (a *App) deleteBook(w http.ResponseWriter, r *http.Request) {
}

// helpers
func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
