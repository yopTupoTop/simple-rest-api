package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetNotesHandler(t *testing.T) {
	notes = []Note{
		{ID: 1, Title: "Test Note 1", Content: "This is the first test note."},
		{ID: 2, Title: "Test Note 2", Content: "This is the second test note."},
	}

	req := httptest.NewRequest(http.MethodGet, "/notes", nil)
	rr := httptest.NewRecorder()

	getNotesHandler(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, status)
		t.Logf("response body: %s", rr.Body.String())
	}

	contentType := rr.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", contentType)
	}

	var result []Note
	if err := json.Unmarshal(rr.Body.Bytes(), &result); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if len(result) != len(notes) {
		t.Errorf("expected %d notes, got %d", len(notes), len(result))
	}
}

func TestGetNotesHandlerWithFilter(t *testing.T) {
	notes = []Note{
		{ID: 1, Title: "Go Basics", Content: "Learning Go."},
		{ID: 2, Title: "Advanced", Content: "Concurrency in Go."},
	}

	req := httptest.NewRequest(http.MethodGet, "/notes?title=go", nil)
	rr := httptest.NewRecorder()

	getNotesHandler(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, status)
	}

	var result []Note
	if err := json.Unmarshal(rr.Body.Bytes(), &result); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if len(result) != 1 {
		t.Errorf("expected 2 notes, got %d", len(result))
	}
}

func TestCreateNoteHandler(t *testing.T) {
	notes = []Note{}
	nextID = 1

	newNote := `{"title":"New Note","content":"This is a new note"}`
	req := httptest.NewRequest(http.MethodPost, "/notes", strings.NewReader(newNote))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	createNoteHandler(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("expected status code %d, got %d", http.StatusCreated, status)
		t.Logf("Response body: %s", rr.Body.String())
	}

	contentType := rr.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", contentType)
	}

	var result Note
	if err := json.Unmarshal(rr.Body.Bytes(), &result); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if result.ID != 1 || result.Title != "New Note" || result.Content != "This is a new note" {
		t.Errorf("unexpected note data: %+v", result)
	}

	if len(notes) != 1 {
		t.Errorf("expected 1 note, got %d", len(notes))
	}

	t.Logf("Notes after request: %+v", notes)
	t.Logf("Next ID after request: %d", nextID)
}

func TestUpdateNoteHandler(t *testing.T) {
	notes = []Note{
		{ID: 1, Title: "Original Title", Content: "Original Content"},
	}
	reqBody := `{"title":"Updated Title","content":"Updated Content"}`
	req := httptest.NewRequest(http.MethodPut, "/notes/1", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	updateNoteHandler(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, status)
	}

	var result Note
	if err := json.Unmarshal(rr.Body.Bytes(), &result); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if result.Title != "Updated Title" || result.Content != "Updated Content" {
		t.Errorf("unexpected updated note data: %+v", result)
	}

	if notes[0].Title != "Updated Title" || notes[0].Content != "Updated Content" {
		t.Errorf("note not updated in storage: %+v", notes[0])
	}
}

func TestGetNoteByIDHandler(t *testing.T) {
	notes = []Note{
		{ID: 1, Title: "Note 1", Content: "First note content"},
	}

	req := httptest.NewRequest(http.MethodGet, "/notes/1", nil)
	rr := httptest.NewRecorder()

	getNoteByIDHandler(rr, req)

	// Проверяем статус ответа
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, status)
	}

	// Проверяем содержимое ответа
	var result Note
	if err := json.Unmarshal(rr.Body.Bytes(), &result); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if result.ID != 1 || result.Title != "Note 1" || result.Content != "First note content" {
		t.Errorf("unexpected note data: %+v", result)
	}
}

func TestDeleteNoteHandler(t *testing.T) {
	notes = []Note{
		{ID: 1, Title: "First Note", Content: "Content of the first note"},
		{ID: 2, Title: "Second Note", Content: "Content of the second note"},
	}

	req := httptest.NewRequest(http.MethodDelete, "/notes/1", nil)
	rr := httptest.NewRecorder()

	deleteNoteHandler(rr, req)

	if len(notes) != 1 {
		t.Errorf("expected 1 note after deletion, got %d", len(notes))
	}

	if notes[0].ID != 2 {
		t.Errorf("expected remaining note ID to be 2, got %d", notes[0].ID)
	}
}
