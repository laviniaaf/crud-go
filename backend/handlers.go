// handlers = logic of manipulation

package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func createItem(w http.ResponseWriter, r *http.Request) {
	var bill Bill

	// Decode the JSON request body into the bill struct
	if err := json.NewDecoder(r.Body).Decode(&bill); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	bill.ID = uuid.New()
	now := time.Now()

	if bill.CreatedAt.IsZero() {
		bill.CreatedAt = now
	}
	bill.UpdatedAt = now

	_, err := db.Exec(
		"INSERT INTO bills (id, embasa, coelba, created_at, updated_at) VALUES (UNHEX(REPLACE(?, '-', '')), ?, ?, ?, ?)",
		bill.ID.String(),
		bill.Embasa,
		bill.Coelba,
		bill.CreatedAt,
		bill.UpdatedAt,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bill)
}

func readItems(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(
		`SELECT
		LOWER(CONCAT(
			SUBSTR(HEX(id), 1, 8), '-',
			SUBSTR(HEX(id), 9, 4), '-',
			SUBSTR(HEX(id), 13, 4), '-',
			SUBSTR(HEX(id), 17, 4), '-',
			SUBSTR(HEX(id), 21, 12)
		)) as id,
		embasa,
		coelba,
		created_at,
		updated_at
	FROM bills`)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	var bills []Bill

	for rows.Next() {
		var b Bill
		var idStr string
		err := rows.Scan(
			&idStr,
			&b.Embasa,
			&b.Coelba,
			&b.CreatedAt,
			&b.UpdatedAt,
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		b.ID, err = uuid.Parse(idStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		bills = append(bills, b)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bills)
}

func updateItem(w http.ResponseWriter, r *http.Request) {

	idStr := chi.URLParam(r, "id") // the id (comes from the URL) as a string, but in the database it is saved as BINARY(16) so
	// need to parse it to a uuid object

	id, err := uuid.Parse(idStr)

	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	var bill Bill

	if err := json.NewDecoder(r.Body).Decode(&bill); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Update the updated_at for the timestamp
	bill.UpdatedAt = time.Now()

	_, err = db.Exec(
		"UPDATE bills SET embasa = ?, coelba = ?, updated_at = ? WHERE id = UNHEX(REPLACE(?, '-', ''))",
		bill.Embasa,
		bill.Coelba,
		bill.UpdatedAt,
		id.String(),
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	bill.ID = id

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bill)
}

func deleteItem(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	_, err = db.Exec("DELETE FROM bills WHERE id = UNHEX(REPLACE(?, '-', ''))", id.String())

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
}

func getBillsByDateRange(w http.ResponseWriter, r *http.Request) {
	// Log
	queryParams := r.URL.Query()
	startDateStr := queryParams.Get("start")
	endDateStr := queryParams.Get("end")

	//log.Printf("=== NEW FILTER REQUEST ===")
	//log.Printf("URL parameters: %v", queryParams)
	//log.Printf("Start date (raw): %s", startDateStr)
	//log.Printf("End date (raw): %s", endDateStr)

	if startDateStr == "" && endDateStr == "" {
		log.Println("Not have a date range, returning all items")
		readItems(w, r)
		return
	}

	var startDate, endDate time.Time
	var err error

	if startDateStr != "" {
		startDate, err = time.Parse("2006-01-02", startDateStr)
		if err != nil {
			http.Error(w, "Invalid start date format. Use YYYY-MM-DD", http.StatusBadRequest)
			return
		}
	} else {
		// If not have a start date, use a very old date
		startDate = time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	}

	if endDateStr != "" {
		endDate, err = time.Parse("2006-01-02", endDateStr)
		if err != nil {
			http.Error(w, "Invalid end date format. Use YYYY-MM-DD", http.StatusBadRequest)
			return
		}
		// Add one day to include the end date
		endDate = endDate.Add(24 * time.Hour)
	} else {
		// If not have an end date, use the current date
		endDate = time.Now()
	}

	//log.Printf("Processed dates - Start: %v (%T), End: %v (%T)\n", startDate, startDate, endDate, endDate)

	// Log the SQL query
	query := `SELECT 
            LOWER(CONCAT(
                SUBSTR(HEX(id), 1, 8), '-',
                SUBSTR(HEX(id), 9, 4), '-',
                SUBSTR(HEX(id), 13, 4), '-',
                SUBSTR(HEX(id), 17, 4), '-',
                SUBSTR(HEX(id), 21, 12)
            )) as id,
            embasa, 
            coelba, 
            created_at, 
            updated_at 
        FROM bills 
        WHERE created_at >= ? AND created_at < ?
        ORDER BY created_at DESC`

	log.Printf("Parameters: startDate=%v (%T), endDate=%v (%T)\n",
		startDate, startDate, endDate, endDate)

	rows, err := db.Query(query, startDate, endDate)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	log.Println("Query executed successfully")

	var bills []Bill

	for rows.Next() {
		var bill Bill
		var idStr string

		err := rows.Scan(
			&idStr,
			&bill.Embasa,
			&bill.Coelba,
			&bill.CreatedAt,
			&bill.UpdatedAt,
		)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		bill.ID, err = uuid.Parse(idStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		bills = append(bills, bill)
	}

	if err = rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(bills)
}
