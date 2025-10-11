package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func main() {

	// starts the database and creates the table
	db = InitDB()

	defer db.Close()

	r := chi.NewRouter()

	// middleware CORS
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders: []string{"Link"},
		MaxAge:         300,
	}))

	r.Route("/bills", func(r chi.Router) {
		r.Post("/", createItem)
		r.Get("/", getBillsByDateRange)
		r.Get("/all", readItems)
		r.Put("/{id}", updateItem)
		r.Delete("/{id}", deleteItem)
	})

	frontendPath := "./frontend"
	fs := http.FileServer(http.Dir(frontendPath))

	r.Get("/*", func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path == "/" {
			http.ServeFile(w, req, frontendPath+"/index.html")
			return
		}
		fs.ServeHTTP(w, req)
	})

	log.Printf("Server running at http://localhost:8082 <---")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
