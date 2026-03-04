package main

import (
	"database/sql"
	"log"
	"net/url"
)

func handleUpdateUPDATE(query url.Values) string {
	db, err := sql.Open("sqlite", "file:../series.db")
	if err != nil {
		log.Print("Error opening DB:", err)
	}
	defer db.Close()

	idStr := query.Get("id")
	name := query.Get("name")
	currentEp := query.Get("current_episode")
	totalEp := query.Get("total_episodes")

	log.Print("Updating with ", idStr, name, currentEp, totalEp)

	response := "HTTP/1.1 303 See Other\r\n" + "Location: /\r\n"
	return response
}
