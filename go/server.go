package main

import (
	"bufio"
	"database/sql"
	"io"
	"log"
	_ "modernc.org/sqlite"
	"net"
	"net/url"
	"strconv"
	"strings"
)

func get(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	requestLine, _ := reader.ReadString('\n')

	parts := strings.Fields(requestLine) // Esto parte un string según cada espacio

	// segunda palabra del line es el path
	method := parts[0]
	path := parts[1]
	log.Print("Path requested: ", path)
	log.Print("Method: ", method)

	u, err := url.ParseRequestURI(path)
	rawQuery := ""
	if err == nil {
		rawQuery = u.RawQuery
		path = u.Path
	}

	// chuparse los headers
	contentLength := 0
	for {
		headerLine, _ := reader.ReadString('\n')
		if strings.HasPrefix(headerLine, "Content-Length:") {
			lengthStr := strings.TrimSpace(strings.TrimPrefix(headerLine, "Content-Length:"))
			contentLength, _ = strconv.Atoi(lengthStr)
		}
		if headerLine == "\r\n" {
			break
		}
	}

	response := ""
	if path == "/" && method == "GET" {
		response = handleHome()
	}

	if path == "/create" && method == "GET" {
		response = handleCreate()
	}

	if path == "/create" && method == "POST" && contentLength > 0 {
		bytesBody := make([]byte, contentLength)
		_, err := io.ReadFull(reader, bytesBody)
		if err != nil {
			log.Print("Error leyendo body: ", err)
			response = "HTTP/1.1 400 Bad Request\r\nConnection: close\r\n\r\n"
		}
		body := string(bytesBody)

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

		response = "HTTP/1.1 303 See Other\r\n" +
			"Location: /\r\n" +
			"Connection: close \r\n" +
			"\r\n"
	}

	if path == "/update/" && method == "POST" {
		query, _ := url.ParseQuery(rawQuery)
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
		if newEp+mult < 0 {
			newEp = 0
		}
		if newEp+mult > episodes {
			newEp = episodes
		}

		_, err2 := db.Exec("UPDATE series SET current_episode=? WHERE id=?", newEp, id)

		if err2 != nil {
			log.Print("Error updating DB", err)
		}
		response = "HTTP/1.1 200 OK\r\nConnection: close\r\n\r\n"
	}

	_, writer := conn.Write([]byte(response))

	if writer != nil {
		log.Print("Error writing to connection: ", writer)
	}
}

func main() {
	// SIEMPRE HAY QUE CERRAR LA CONEXIÓN CUANDO SE CREA
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Print("Error listening to port.")
	}
	log.Print("listening to port 8080")

	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print("Error accepting connection")
		}
		go get(conn)
	}
}
