package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	_ "modernc.org/sqlite"
	"net"
	"strings"
)

type Serie struct {
	id        int
	name      string
	currentEp int
	episodes  int
}

func get(conn net.Conn) {
	db, _ := sql.Open("sqlite", "file:series.db")
	defer db.Close()
	defer conn.Close()

	reader := bufio.NewReader(conn)

	requestLine, _ := reader.ReadString('\n')

	parts := strings.Fields(requestLine) // Esto parte un string según cada espacio

	// segunda palabra del line es el path
	path := parts[1]
	log.Print("Path requested: ", path)

	// chuparse los headers
	for {
		headerLine, _ := reader.ReadString('\n')
		log.Print("Header: " + headerLine)
		if headerLine == "\r\n" {
			break
		}
	}

	if path == "/" {
		rows, err := db.Query("SELECT * FROM series")
		if err != nil {
			log.Print("Error querying for series: ", err)
		}

		defer rows.Close()

		var seriesList []Serie
		for rows.Next() {
			var serie Serie
			rows.Scan(&serie.id, &serie.name, &serie.currentEp, &serie.episodes)
			seriesList = append(seriesList, serie)
			log.Print("Serie leída: ", serie.id, serie.name, serie.currentEp, serie.episodes)
		}

		var tableRowsString string

		for _, s := range seriesList {
			tableRowsString += `<tr>
				<td>` + s.name + `</td>
				<td>` + fmt.Sprintf("%d", s.currentEp) + `</td>
				<td>` + fmt.Sprintf("%d", s.episodes) + `</td>
				</tr>
				`
		}

		// LOS RESPONSE VAN: http/ver, responseID, text
		response := "HTTP/1.1 200 OK\r\n"
		headers := "Content-Type: text/html\r\n" +
			"\r\n" +
			`
		<!DOCTYPE html>
		<html lang="en">
		<style>
			table, th, td {
				border: 1px solid black;
				text-align: center;
			}
		</style>
		<head>
		<meta charset = "UTF-8"/>
		<title>Series que estoy viendo</title>
		</head>
		<body>
		<h1>Lista de series que he visto en algún momento o estoy viendo actualmente</h1>
		<table>
		<tr>
		<th>Nombre</th>
		<th>Episodio Actual que Voy/me quedé</th>
		<th>Cantidad de episodios totales</th>
		` + tableRowsString +

			`
		</table>
		</body>
		</html>
		`

		_, writerr := conn.Write([]byte(response + headers))

		if writerr != nil {
			log.Print("Error writing to connection: ", writerr)
		}
	}

	/*
		response := "HTTP/1.1 200 OK\r\n"
		headers := "Content-Type: text/html\r\n" +
			"\r\n" +
			`
		<!DOCTYPE html>
		<html lang="en">
		<head>
		<meta charset = "UTF-8"/>
		<title>A</title>
		</head>
		<body>
		<h1>Prueba</h1>
		</body>
		</html>
		`
		_, err := conn.Write([]byte(response + headers))

		if err != nil {
			log.Print("Error writing in conn")
		}
	*/
}

// LAS REQUESTS DE HTTP SIGUEN EL FORMATO:
// GET /path http:1.3.2

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
