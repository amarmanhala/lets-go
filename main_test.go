package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
)

func TestHomeHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	homeHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var response map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if response["endpoint"] != "/api/countries" {
		t.Fatalf("expected endpoint /api/countries, got %q", response["endpoint"])
	}
}

func TestHomeHandlerNotFound(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/missing", nil)
	rec := httptest.NewRecorder()

	homeHandler(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
}

func TestCountriesHandlerReturnsCountries(t *testing.T) {
	db, mock := newMockDB(t)
	defer db.Close()

	mock.ExpectQuery(countriesQueryPattern()).
		WillReturnRows(sqlmock.NewRows([]string{"name", "population", "flag_colors", "language"}).
			AddRow("Canada", int64(40097761), pq.Array([]string{"red", "white"}), "English").
			AddRow("Japan", int64(123294513), pq.Array([]string{"white", "red"}), "Japanese"))

	req := httptest.NewRequest(http.MethodGet, "/api/countries", nil)
	rec := httptest.NewRecorder()

	countriesHandler(db)(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
	if rec.Header().Get("Content-Type") != "application/json" {
		t.Fatalf("expected application/json content type, got %q", rec.Header().Get("Content-Type"))
	}

	var response []Country
	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if len(response) != 2 {
		t.Fatalf("expected 2 countries, got %d", len(response))
	}
	if response[0].Name != "Canada" || response[0].Population != 40097761 {
		t.Fatalf("unexpected first country: %+v", response[0])
	}
	if len(response[0].FlagColors) != 2 || response[0].FlagColors[0] != "red" {
		t.Fatalf("unexpected flag colors: %+v", response[0].FlagColors)
	}

	assertNoSQLMockExpectations(t, mock)
}

func TestCountriesHandlerRejectsNonGet(t *testing.T) {
	db, mock := newMockDB(t)
	defer db.Close()

	req := httptest.NewRequest(http.MethodPost, "/api/countries", nil)
	rec := httptest.NewRecorder()

	countriesHandler(db)(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected status %d, got %d", http.StatusMethodNotAllowed, rec.Code)
	}
	if rec.Header().Get("Allow") != http.MethodGet {
		t.Fatalf("expected Allow header %q, got %q", http.MethodGet, rec.Header().Get("Allow"))
	}

	assertNoSQLMockExpectations(t, mock)
}

func TestCountriesHandlerReturnsServerErrorOnDatabaseFailure(t *testing.T) {
	db, mock := newMockDB(t)
	defer db.Close()

	mock.ExpectQuery(countriesQueryPattern()).WillReturnError(errors.New("database unavailable"))

	req := httptest.NewRequest(http.MethodGet, "/api/countries", nil)
	rec := httptest.NewRecorder()

	countriesHandler(db)(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, rec.Code)
	}

	assertNoSQLMockExpectations(t, mock)
}

func TestGetCountriesScansRows(t *testing.T) {
	db, mock := newMockDB(t)
	defer db.Close()

	mock.ExpectQuery(countriesQueryPattern()).
		WillReturnRows(sqlmock.NewRows([]string{"name", "population", "flag_colors", "language"}).
			AddRow("Brazil", int64(216422446), pq.Array([]string{"green", "yellow", "blue", "white"}), "Portuguese"))

	result, err := getCountries(db)
	if err != nil {
		t.Fatalf("getCountries returned error: %v", err)
	}

	if len(result) != 1 {
		t.Fatalf("expected 1 country, got %d", len(result))
	}
	if result[0].Name != "Brazil" || result[0].Language != "Portuguese" {
		t.Fatalf("unexpected result: %+v", result[0])
	}

	assertNoSQLMockExpectations(t, mock)
}

func newMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	t.Helper()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("create sql mock: %v", err)
	}

	return db, mock
}

func countriesQueryPattern() string {
	return `SELECT name, population, flag_colors, language\s+FROM countries\s+ORDER BY name`
}

func assertNoSQLMockExpectations(t *testing.T, mock sqlmock.Sqlmock) {
	t.Helper()

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet SQL mock expectations: %v", err)
	}
}
