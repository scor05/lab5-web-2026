package main

import (
	"database/sql"
	"log"
	"net/url"
)

func handleCreatePOST(body string) string {
	data, _ := url.ParseQuery(body)
	name := data.Get("series_name")
	currentEp := data.Get("current_episode")
	episodes := data.Get("total_episodes")

	db, err := sql.Open("sqlite", "file:../series.db")
	if err != nil {
		log.Print("Error opening database: ", err)
	}

	defer db.Close()
	_, err1 := db.Exec("INSERT INTO series VALUES (NULL, ?, ?, ?)", name, currentEp, episodes)
	if err1 != nil {
		log.Print("Error inserting into DB: ", err)
	}

	response := "HTTP/1.1 303 See Other\r\n" +
		"Location: /\r\n" +
		"Connection: close \r\n" +
		"\r\n"

	return response
}
