package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func EnsureTableExists(a *App, table string) {
	if _, err := a.DB.Exec(table); err != nil {
		log.Fatal(err)
	}
}

func ClearTable(a *App, table string) {
	q := fmt.Sprintf("DELETE FROM %s", table)
	a.DB.Exec(q)
	q = fmt.Sprintf("ALTER SEQUENCE %s_id_seq RESTART WITH 1", table)
	a.DB.Exec(q)
}

func ExecuteRequest(req *http.Request, a *App) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	return rr
}

func CheckResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}
