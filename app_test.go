package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

var a App

func TestMain(m *testing.M) {
	err := a.Initialise(DbUser, DbPassword, "test")
	if err != nil {
		log.Fatal("Error while initialising app: ")
	}
	createTable()
	m.Run()
}

func createTable() {
	createTableQuery := `CREATE TABLE IF NOT EXISTS products (
	id int NOT NULL AUTO_INCREMENT ,
	name varchar(255) NOT NULL,
	quantity int,
	price float(10,7),
	PRIMARY KEY (id)
);`
	_, err := a.DB.Exec(createTableQuery)
	if err != nil {
		log.Fatal(err)
	}
}

func clearTable() {
	a.DB.Exec("DELETE from products")
	a.DB.Exec("ALTER table products AUTO_INCREMENT = 1")
	log.Println("Deleted all rows from table")
	// a.DB.Exec("ALTER TABLE products AUTO_INCREMENT = 1")
}

func addProduct(name string, quantity int, price float64) {
	query := fmt.Sprintf("INSERT into products(name, quantity, price) VALUES('%v',%v,%v)", name, quantity, price)

	_, err := a.DB.Exec(query)
	if err != nil {
		log.Println(err)
	}
}

func TestGetProduct(t *testing.T) {
	clearTable()
	addProduct("keyboard", 100, 100.00)
	request, _ := http.NewRequest("GET", "/product/1", nil)
	response := sendRequest(request)
	checkStatusCode(t, http.StatusOK, response.Code)
}

func checkStatusCode(t *testing.T, expectedStatusCode int, actualStatusCode int) {
	if expectedStatusCode != actualStatusCode {
		t.Errorf("Expected status code %v but got %v", expectedStatusCode, actualStatusCode)
	}
}

func sendRequest(request *http.Request) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()
	a.Router.ServeHTTP(recorder, request)
	return recorder
}

func TestCreateProduct(t *testing.T) {
	clearTable()
	product := []byte(`{"name":"keyboard", "quantity":100, "price":100.00"}`)
	req, _ := http.NewRequest("POST", "/product", bytes.NewBuffer(product))

	req.Header.Set("Content-Type", "application/json")
	response := sendRequest(req)
	checkStatusCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["name"] != "chair" {
		t.Errorf("Expected product name to be 'chair'. Got '%v'", m["name"])
	}

	if m["quantity"] != 1.0 {
		t.Errorf("Expected product quantity to be 1.0, Got '%v'", m["quantity"])
	}
}
