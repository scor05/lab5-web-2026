package main

import (
	"database/sql"
	"log"
	"net/url"
	"strconv"
)

func handleUpdatePOST(query url.Values) string {
	idStr := query.Get("id")
	id, _ := strconv.Atoi(idStr)
	change := query.Get("change")
	mult := 1
	if change == "m" {
		mult *= -1
	}
	log.Print("ID changed: ", id, "; change: ", change)

	db, err := sql.Open("sqlite", "file:../series.db")
	if err != nil {
		log.Print("Error opening DB: ", err)
	}
	defer db.Close()
	var currentEp, episodes int
	err = db.QueryRow("SELECT current_episode, total_episodes FROM series WHERE id=?", id).Scan(&currentEp, &episodes)
	if err != nil {
		log.Print("Error in queryrow: ", err)
	}

	newEp := currentEp + mult
	if newEp < 0 {
		newEp = 0
	}
	if newEp > episodes {
		newEp = episodes
	}

	_, err2 := db.Exec("UPDATE series SET current_episode=? WHERE id=?", newEp, id)

	if err2 != nil {
		log.Print("Error updating DB", err2)
	}
	response := "HTTP/1.1 200 OK\r\nConnection: close\r\n\r\n"

	return response
}
