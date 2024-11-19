package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/notes", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getNotesHandler(w, r)
		case http.MethodPost:
			createNoteHandler(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/notes/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getNoteByIDHandler(w, r)
		case http.MethodPut:
			updateNoteHandler(w, r)
		case http.MethodDelete:
			deleteNoteHandler(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	log.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
