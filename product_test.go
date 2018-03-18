package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"
)

var a App

func TestMain(m *testing.M) {
	debug = true
	a = App{}
	a.Initialize(
		os.Getenv("TEST_DB_USERNAME"),
		os.Getenv("TEST_DB_PASSWORD"),
		os.Getenv("TEST_DB_NAME"))

	EnsureTableExists(&a, productsTableCreationQuery)
	EnsureTableExists(&a, eventsTableCreationQuery)
	EnsureTableExists(&a, venuesTableCreationQuery)

	code := m.Run()

	ClearTable(&a, "products")
	ClearTable(&a, "events")

	os.Exit(code)
}

const productsTableCreationQuery = `CREATE TABLE IF NOT EXISTS products
(
id SERIAL,
name TEXT NOT NULL,
price NUMERIC(10,2) NOT NULL DEFAULT 0.00,
CONSTRAINT products_pkey PRIMARY KEY (id)
)`

const eventsTableCreationQuery = `CREATE TABLE IF NOT EXISTS events
(
id bigserial primary key,
name TEXT NOT NULL,
date timestamp NOT NULL
)`

const venuesTableCreationQuery = `CREATE TABLE IF NOT EXISTS events
(
id bigserial primary key,
name TEXT NOT NULL,
location TEXT NOT NULL,
capacity INT,
event_id integer REFERENCES events
)`

func TestProductsEmptyTable(t *testing.T) {
	ClearTable(&a, "products")

	req, _ := http.NewRequest("GET", "/products", nil)
	response := ExecuteRequest(req, &a)

	CheckResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body != "[]" {
		t.Errorf("Expected an empty array. Got %s", body)
	}
}

func TestGetNonExistentProduct(t *testing.T) {
	ClearTable(&a, "products")

	req, _ := http.NewRequest("GET", "/product/11", nil)
	response := ExecuteRequest(req, &a)

	CheckResponseCode(t, http.StatusNotFound, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "Product not found" {
		t.Errorf("Expected the 'error' key of the response to be set to 'Product not found'. Got '%s'", m["error"])
	}
}

func TestCreateProduct(t *testing.T) {
	ClearTable(&a, "products")

	payload := []byte(`{"name":"test product","price":11.22}`)

	req, _ := http.NewRequest("POST", "/product", bytes.NewBuffer(payload))
	response := ExecuteRequest(req, &a)

	CheckResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["name"] != "test product" {
		t.Errorf("Expected product name to be 'test product'. Got '%v'", m["name"])
	}

	if m["price"] != 11.22 {
		t.Errorf("Expected product price to be '11.22'. Got '%v'", m["price"])
	}

	// the id is compared to 1.0 because JSON unmarshaling converts numbers to
	// floats, when the target is a map[string]interface{}
	if m["id"] != 1.0 {
		t.Errorf("Expected product ID to be '1'. Got '%v'", m["id"])
	}
}
func TestGetProduct(t *testing.T) {
	ClearTable(&a, "products")
	addProducts(1)

	req, _ := http.NewRequest("GET", "/product/1", nil)
	response := ExecuteRequest(req, &a)

	CheckResponseCode(t, http.StatusOK, response.Code)
}
func addProducts(count int) {
	if count < 1 {
		count = 1
	}

	for i := 0; i < count; i++ {
		a.DB.Exec("INSERT INTO products(name, price) VALUES($1, $2)", "Product "+strconv.Itoa(i), (i+1.0)*10)
	}
}
func TestUpdateProduct(t *testing.T) {
	ClearTable(&a, "products")
	addProducts(1)

	req, _ := http.NewRequest("GET", "/product/1", nil)
	response := ExecuteRequest(req, &a)
	var originalProduct map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &originalProduct)

	payload := []byte(`{"name":"test product - updated name","price":11.22}`)

	req, _ = http.NewRequest("PUT", "/product/1", bytes.NewBuffer(payload))
	response = ExecuteRequest(req, &a)

	CheckResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["id"] != originalProduct["id"] {
		t.Errorf("Expected the id to remain the same (%v). Got %v", originalProduct["id"], m["id"])
	}

	if m["name"] == originalProduct["name"] {
		t.Errorf("Expected the name to change from '%v' to '%v'. Got '%v'", originalProduct["name"], m["name"], m["name"])
	}

	if m["price"] == originalProduct["price"] {
		t.Errorf("Expected the price to change from '%v' to '%v'. Got '%v'", originalProduct["price"], m["price"], m["price"])
	}
}
func TestDeleteProduct(t *testing.T) {
	ClearTable(&a, "products")
	addProducts(1)

	req, _ := http.NewRequest("GET", "/product/1", nil)
	response := ExecuteRequest(req, &a)
	CheckResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("DELETE", "/product/1", nil)
	response = ExecuteRequest(req, &a)

	CheckResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("GET", "/product/1", nil)
	response = ExecuteRequest(req, &a)
	CheckResponseCode(t, http.StatusNotFound, response.Code)
}

func TestEventsEmptyTable(t *testing.T) {
	ClearTable(&a, "events")

	req, _ := http.NewRequest("GET", "/events", nil)
	response := ExecuteRequest(req, &a)

	CheckResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body != "[]" {
		t.Errorf("Expected an empty array. Got %s", body)
	}
}

func TestCreateEvent(t *testing.T) {
	ClearTable(&a, "events")

	payload := []byte(`{"name":"test event","date":"2018-03-14T15:55:58Z","venue":{"name":"Code Node","id":123}}`)
	// payload := []byte(`{"name":"test event","date":"2018-03-14T15:55:58Z"}`)
	req, _ := http.NewRequest("POST", "/event", bytes.NewBuffer(payload))
	response := ExecuteRequest(req, &a)

	CheckResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["name"] != "test event" {
		t.Errorf("Expected event name to be 'test event'. Got '%v'", m["name"])
	}

	if m["date"] != "2018-03-14T15:55:58Z" {
		t.Errorf("Expected event date to be 2018-03-14T15:55:58Z. Got '%v'", m["date"])
	}

	// the id is compared to 1.0 because JSON unmarshaling converts numbers to
	// floats, when the target is a map[string]interface{}
	if m["id"] != 1.0 {
		t.Errorf("Expected event ID to be '1'. Got '%v'", m["id"])
	}
}

func TestGetEvent(t *testing.T) {
	ClearTable(&a, "events")
	addEvents(1)

	req, _ := http.NewRequest("GET", "/event/1", nil)
	response := ExecuteRequest(req, &a)

	CheckResponseCode(t, http.StatusOK, response.Code)
}
func addEvents(count int) {
	if count < 1 {
		count = 1
	}
	for i := 0; i < count; i++ {
		_, err := a.DB.Exec("INSERT INTO events(name, date) VALUES($1, $2)", "Event "+strconv.Itoa(i), time.Now())
		if err != nil {
			fmt.Println(err)
		}
	}
}

func TestUpdateEvent(t *testing.T) {
	ClearTable(&a, "events")
	addEvents(1)

	req, _ := http.NewRequest("GET", "/event/1", nil)
	response := ExecuteRequest(req, &a)
	var originalEvent map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &originalEvent)

	payload := []byte(`{"name":"test event - updated name","date":"2018-03-14T15:55:58Z"}`)

	req, _ = http.NewRequest("PUT", "/event/1", bytes.NewBuffer(payload))
	response = ExecuteRequest(req, &a)

	CheckResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["id"] != originalEvent["id"] {
		t.Errorf("Expected the id to remain the same (%v). Got %v", originalEvent["id"], m["id"])
	}

	if m["name"] == originalEvent["name"] {
		t.Errorf("Expected the name to change from '%v' to '%v'. Got '%v'", originalEvent["name"], m["name"], m["name"])
	}

	if m["date"] == originalEvent["date"] {
		t.Errorf("Expected the date to change from '%v' to '%v'. Got '%v'", originalEvent["date"], m["date"], m["date"])
	}
}

func TestGetNonExistentEvent(t *testing.T) {
	ClearTable(&a, "events")

	req, _ := http.NewRequest("GET", "/event/11", nil)
	response := ExecuteRequest(req, &a)

	CheckResponseCode(t, http.StatusNotFound, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "Event not found" {
		t.Errorf("Expected the 'error' key of the response to be set to 'Event not found'. Got '%s'", m["error"])
	}
}

func TestDeleteEvent(t *testing.T) {
	ClearTable(&a, "events")
	addEvents(1)

	req, _ := http.NewRequest("GET", "/event/1", nil)
	response := ExecuteRequest(req, &a)
	CheckResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("DELETE", "/event/1", nil)
	response = ExecuteRequest(req, &a)

	CheckResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("GET", "/event/1", nil)
	response = ExecuteRequest(req, &a)
	CheckResponseCode(t, http.StatusNotFound, response.Code)
}
