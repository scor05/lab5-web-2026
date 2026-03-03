package main

import (
	"database/sql"
	"fmt"
	"log"
)

func handleHome() string {
	db, _ := sql.Open("sqlite", "file:series.db")
	defer db.Close()
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
		<td>` + fmt.Sprintf("%d", s.id) + `</td>
		<td>` + s.name + `</td>
		<td>` + fmt.Sprintf("%d", s.currentEp) + `</td>
		<td>` + fmt.Sprintf("%d", s.episodes) + `</td>
		<td> <progress id='` + fmt.Sprintf("%d", s.id) + `' value='` + fmt.Sprintf("%d", s.currentEp) + `' max='` + fmt.Sprintf("%d", s.episodes) + `'></progress>
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
	body {
		font-family: "Arial Narrow";
		font-size: 24px;
	}
    .container{
        display: flex;
        flex-direction: column;
        justify-content: center;
        align-items: center;
        width: 100%;
        height: 100vh;
    }
	table {
		border-radius: 5px;
	}
	table, th, td {
		border: 1px solid black;
		text-align: center;
	}
	th, td {
		padding: 10px;
	}
	tr:nth-of-type(even){
		background-color: aqua;
	}
	tr:nth-of-type(odd){
		background-color: cadetBlue;
	}
	</style>
	<head>
	<script src="./home.js"></script>
	<meta charset = "UTF-8"/>
	<title>Mi Tracker de Series</title>
	</head>
	<body>
    <div class="container">
	<h1>Mi Tracker de Series</h1>
	<table>
	<tr>
	<th>ID</th>
	<th>Nombre</th>
	<th>Episodio Actual</th>
	<th>Episodios Totales</th>
	<th>Progreso</th>
	` + tableRowsString +

		`
	</table>
	<p><a href="./create">Agregar una Serie</a></p>
    </div>
	</body>
	</html>
	`

	return response + headers
}
