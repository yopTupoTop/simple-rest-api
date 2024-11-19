package main

import (
	"errors"
	"fmt"
)

type Note struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

var (
	notes  []Note
	nextID = 1
)

func findNoteByID(id int) (*Note, int, error) {
	for i, note := range notes {
		if note.ID == id {
			return &note, i, nil
		}
	}
	return nil, -1, errors.New("note not found")
}

func deleteNoteByIndex(index int) {
	notes = append(notes[:index], notes[index+1:]...)
	fmt.Println("notes:", notes)
}
