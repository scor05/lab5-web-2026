package main

import (
	"database/sql"
	"log"
	"net/url"
	"strconv"
)

func handleUpdateUPDATE(query url.Values) string {
	db, err := sql.Open("sqlite", "file:../series.db")
	if err != nil {
		log.Print("Error opening DB:", err)
	}
	defer db.Close()

	idStr := query.Get("id")
	id, _ := strconv.Atoi(idStr)
	originalId, _ := strconv.Atoi(query.Get("original_id"))
	name := query.Get("series_name")
	currentEp := query.Get("current_episode")
	totalEp := query.Get("total_episodes")

	var exists int
	err2 := db.QueryRow(`SELECT EXISTS(SELECT 1 FROM series WHERE id = ?)`, id).Scan(&exists)
	if err2 != nil {
		log.Print("Error querying exists:", err2)
	}

	body := ""
	response := ""
	if exists == 1 && id != originalId {
		body = `<script>
        alert("Ya existe una serie con ese ID");
        history.back();
        </script>`

		response = "HTTP/1.1 422 Unprocessable Entity\r\n" +
			"Content-Type: text/html\r\n" +
			"Content-Length: " + strconv.Itoa(len(body)) + "\r\n" +
			"Connection: close\r\n" +
			"\r\n"
	} else {
		_, err := db.Exec("UPDATE series SET id=?, name=?, current_episode=?, total_episodes=? WHERE id=?", id, name, currentEp, totalEp, originalId)
		if err != nil {
			log.Print("Error in sql query:", err)
		}

		body = `
            <script>
                alert("Serie modificada exitosamente");
                window.location = "/";
            </script>`

		response = "HTTP/1.1 200 OK\r\n" +
			"Content-Type: text/html\r\n" +
			"Content-Length: " + strconv.Itoa(len(body)) + "\r\n" +
			"Connection: close\r\n" +
			"\r\n"
	}
	return response + body
}
