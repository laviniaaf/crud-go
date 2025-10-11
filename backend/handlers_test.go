package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-chi/chi/v5"
	_ "github.com/go-sql-driver/mysql"
)

// the handler will return the item data in JSON format in the HTTP response body
// to be able to check the bank values (id, name and price)
// in the accurate test converter the JSON back to a struct using marshal
// without this conversation would not give to make the checks of the item fields

// any type of data in go can be converted to []byte through the marshal function
// the marshal = converts any value (struct, map) to a byte array ([]byte) that becomes JSON format

// mockDB = global variable for database mock
var mockDB sqlmock.Sqlmock

func TestMain(m *testing.M) {

	var err error

	db, mockDB, err = sqlmock.New() // creates database mock

	if err != nil {
		panic("Error creating database mock: " + err.Error())
	}

	m.Run()
}

func TestCreateItem(t *testing.T) {

	t.Run("success", func(t *testing.T) {

		item := Item{
			Name:  "Computer",
			Price: 990.50,
		}

		body, _ := json.Marshal(item) // transforms the item into JSON

		mockDB.ExpectExec("INSERT INTO items").
			WithArgs(item.Name, item.Price).
			WillReturnResult(sqlmock.NewResult(1, 1))

		req := httptest.NewRequest("POST", "/itens", bytes.NewReader(body))

		rr := httptest.NewRecorder()

		createItem(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("Unexpected status code: received %v, expected %v", status, http.StatusOK)
		}

		var createdItem Item

		err := json.Unmarshal(rr.Body.Bytes(), &createdItem) // decodes response JSON

		if err != nil {
			t.Fatal("Error decoding JSON response:", err)
		}

		if createdItem.ID != 1 {
			t.Errorf("Unexpected ID: received %d, expected %d", createdItem.ID, 1)
		}

		if err := mockDB.ExpectationsWereMet(); err != nil {
			t.Errorf("Mock expectations not met: %s", err)
		}
	})

	t.Run("json_invalid", func(t *testing.T) {

		req := httptest.NewRequest("POST", "/itens", bytes.NewBufferString("json-invalid"))

		rr := httptest.NewRecorder()

		createItem(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("Unexpected status code: received %v, expected %v", status, http.StatusBadRequest)
		}
	})
}

func TestReadItems(t *testing.T) {

	rows := sqlmock.NewRows([]string{"id", "name", "price"}).
		AddRow(1, "Computer", 990.50).
		AddRow(2, "Keyboard Redragon", 145.99)

	mockDB.ExpectQuery("SELECT id, name, price FROM items").
		WillReturnRows(rows)

	req := httptest.NewRequest("GET", "/itens", nil)

	rr := httptest.NewRecorder()

	readItems(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Unexpected status code: received %v, expected %v", status, http.StatusOK)
	}

	var items []Item

	err := json.Unmarshal(rr.Body.Bytes(), &items)

	if err != nil {
		t.Fatal("Error decoding JSON response:", err)
	}

	if len(items) != 2 {
		t.Errorf("Unexpected number of items: received %d, expected %d", len(items), 2)
	}

	if items[0].Name != "Computer" {
		t.Errorf("Unexpected item name: received %s, expected %s", items[0].Name, "Computer")
	}

	if err := mockDB.ExpectationsWereMet(); err != nil {
		t.Errorf("Mock expectations not met: %s", err)
	}
}

func TestUpdateItem(t *testing.T) {

	t.Run("success", func(t *testing.T) {
		id := 1
		item := Item{
			Name:  "Mouse Update",
			Price: 75.00,
		}

		body, _ := json.Marshal(item)

		r := chi.NewRouter()

		r.Put("/itens/{id}", updateItem)

		req := httptest.NewRequest("PUT", "/itens/"+strconv.Itoa(id), bytes.NewReader(body))

		rr := httptest.NewRecorder()

		mockDB.ExpectExec("UPDATE items").
			WithArgs(item.Name, item.Price, id).
			WillReturnResult(sqlmock.NewResult(0, 1))

		r.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("Unexpected status code: received %v, expected %v", status, http.StatusOK)
		}

		if err := mockDB.ExpectationsWereMet(); err != nil {
			t.Errorf("Mock expectations not met: %s", err)
		}
	})

	t.Run("id_invalid", func(t *testing.T) {
		req := httptest.NewRequest("PUT", "/itens/not-is-a-id", nil)

		rr := httptest.NewRecorder()

		updateItem(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("Unexpected status code: received %v, expected %v", status, http.StatusBadRequest)
		}
	})
}

func TestDeleteItem(t *testing.T) {

	t.Run("success", func(t *testing.T) {

		id := 1

		r := chi.NewRouter()

		r.Delete("/itens/{id}", deleteItem)

		req := httptest.NewRequest("DELETE", "/itens/"+strconv.Itoa(id), nil)

		rr := httptest.NewRecorder()

		mockDB.ExpectExec("DELETE FROM items").
			WithArgs(id).
			WillReturnResult(sqlmock.NewResult(0, 1))

		r.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("Unexpected status code: received %v, expected %v", status, http.StatusOK)
		}

		if err := mockDB.ExpectationsWereMet(); err != nil {
			t.Errorf("Mock expectations not met: %s", err)
		}
	})

	t.Run("id_invalid", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/itens/not-is-a-id", nil)

		rr := httptest.NewRecorder()

		deleteItem(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("Unexpected status code: received %v, expected %v", status, http.StatusBadRequest)
		}
	})
}
