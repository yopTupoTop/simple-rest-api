package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// Получить все заметки
func getNotesHandler(w http.ResponseWriter, r *http.Request) {
	titleFilter := r.URL.Query().Get("title")

	filtredNotes := notes
	if titleFilter != "" {
		filtredNotes = []Note{}
		for _, note := range notes {
			if strings.Contains(strings.ToLower(note.Title), strings.ToLower(titleFilter)) {
				filtredNotes = append(filtredNotes, note)
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(filtredNotes)
}

// Создать новую заметку
func createNoteHandler(w http.ResponseWriter, r *http.Request) {
	var newNote Note
	if err := json.NewDecoder(r.Body).Decode(&newNote); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	newNote.ID = nextID
	nextID++
	notes = append(notes, newNote)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newNote)
}

func updateNoteHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/notes/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	var updatedNote Note
	if err := json.NewDecoder(r.Body).Decode(&updatedNote); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	note, index, err := findNoteByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	note.Title = updatedNote.Title
	note.Content = updatedNote.Content
	notes[index] = *note

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(note)
}

// Получить заметку по ID
func getNoteByIDHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/notes/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	for _, note := range notes {
		if note.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(note)
			return
		}
	}

	http.Error(w, "Note not found", http.StatusNotFound)
}

// Удалить заметку по ID
func deleteNoteHandler(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, "/notes/") {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	idStr := strings.TrimPrefix(r.URL.Path, "/notes/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	fmt.Println("Request URL:", r.URL.Path)
	fmt.Println("Notes before:", notes)

	_, index, err := findNoteByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	deleteNoteByIndex(index)
}
