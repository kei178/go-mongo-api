package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var a App

// test setups
func TestMain(m *testing.M) {
	a = App{}
	dbName = "gomongoapi_test"
	a.Initialize("root", "")

	code := m.Run()
	clearDB()

	os.Exit(code)
}

// helpers
func clearDB() {
	ctx := dbContext(10)
	err := a.DB.Drop(ctx)
	if err != nil {
		log.Fatalf("teardownDatabase failed: %v", err)
	}
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func addBooks(count int) []string {
	if count < 1 {
		count = 1
	}

	ids := []string{}
	for i := 0; i < count; i++ {
		col := a.DB.Collection(bookCollectionName)
		ctx := dbContext(5)
		b := Book{Title: ("Title " + strconv.Itoa(i+1)), Author: ("Author " + strconv.Itoa(i+1))}
		result, err := col.InsertOne(ctx, b)
		if err != nil {
			log.Fatalf("insertTestData failed: %v", err)
		}

		ids = append(ids, result.InsertedID.(primitive.ObjectID).Hex())
	}
	return ids
}

// tests
func TestGetNonExistentBook(t *testing.T) {
	clearDB()

	req, _ := http.NewRequest("GET", "/book/1", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "Book not found" {
		t.Errorf("Expected the 'error' key of the response to be set to 'User not found'. Got '%s'", m["error"])
	}
}

func TestGetBooks(t *testing.T) {
	clearDB()
	addBooks(10)

	req, _ := http.NewRequest("GET", "/books", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestGetBook(t *testing.T) {
	clearDB()
	id := addBooks(1)[0]

	req, _ := http.NewRequest("GET", "/book/"+id, nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestCreateBook(t *testing.T) {
	clearDB()

	payload := []byte(`{"title":"Golang for beginners","author":"kei178"}`)

	req, _ := http.NewRequest("POST", "/book", bytes.NewBuffer(payload))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusCreated, response.Code)
}

func TestUpdateBook(t *testing.T) {
	clearDB()
	id := addBooks(1)[0]

	payload := []byte(`{"title":"Golang for beginners","author":"kei178"}`)
	req, _ := http.NewRequest("PUT", "/book/"+id, bytes.NewBuffer(payload))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestDeleteBook(t *testing.T) {
	clearDB()
	id := addBooks(1)[0]

	req, _ := http.NewRequest("DELETE", "/book/"+id, nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("GET", "/book/"+id, nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)
}
