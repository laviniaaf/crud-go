// models = data structures

package main

// to send struct items to http, it must be in  JSON format, 
// since HTTP requests work with JSON as text

type Item struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}
