package main

import (
	"database/sql"
	"log"
	"net/url"
	"strconv"
)

func handleDelete(query url.Values) string {
	db, err := sql.Open("sqlite", "file:../series.db")
	if err != nil {
		log.Print("Error opening DB:", err)
	}
	defer db.Close()

	idStr := query.Get("id")
	id, _ := strconv.Atoi(idStr)

	_, err2 := db.Exec("DELETE FROM series WHERE id=?", id)
	if err2 != nil {
		log.Print("Error deleting from db:", err2)
	}

	response := "HTTP/1.1 200 OK\r\nConnection: close\r\n\r\n"
	return response
}
