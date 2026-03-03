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
	var buttonsUp string
	var buttonsDown string
	for _, s := range seriesList {
		tableRowsString += `<tr>
		<td>` + fmt.Sprintf("%d", s.id) + `</td>
		<td>` + s.name + `</td>
        <td id='ep` + fmt.Sprintf("%d", s.id) + `'>` + fmt.Sprintf("%d", s.currentEp) + `</td>
		<td id='tot` + fmt.Sprintf("%d", s.id) + `'>` + fmt.Sprintf("%d", s.episodes) + `</td>
		<td> <progress id='p` + fmt.Sprintf("%d", s.id) + `' value='` + fmt.Sprintf("%d", s.currentEp) + `' max='` + fmt.Sprintf("%d", s.episodes) + `'></progress>
		</tr>`

		buttonsDown += `<button onclick="prevEpisode(` + fmt.Sprintf("%d", s.id) + `)">-1</button>`
		buttonsUp += `<button onclick="nextEpisode(` + fmt.Sprintf("%d", s.id) + `)">+1</button>`
	}

	script := `<script type="module">
    window.nextEpisode = async function nextEpisode(id) {
        const url = "/update/?id=" + id + "&change=p";
        const response = await fetch(url, { method: "POST" });

        const textElement = document.getElementById("ep" + id);
        const progressElement = document.getElementById("p" + id);
        const totalEp = parseInt(document.getElementById("tot" + id).innerText);
        const current = parseInt(textElement.innerText);
        if ((current + 1) <= totalEp){
            textElement.innerText = String(current + 1);
            progressElement.value = String(current + 1);
        }
    }
    window.prevEpisode = async function prevEpisode(id) {
        const url = "/update/?id=" + id + "&change=m";
        const response = await fetch(url, { method: "POST" });

        const textElement = document.getElementById("ep" + id);
        const progressElement = document.getElementById("p" + id);
        const totalEp = parseInt(document.getElementById("tot" + id).innerText);
        const current = parseInt(textElement.innerText);
        if ((current - 1) >= 0){
            textElement.innerText = String(current - 1);
            progressElement.value = String(current - 1);
        }
    }
    </script>`
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
    tr:nth-of-type(even), button:nth-of-type(odd){
		background-color: aqua;
	}
    tr:nth-of-type(odd), button:nth-of-type(even){
		background-color: cadetBlue;
	}
    button{
        font-size: 20px;
        border-radius: 10px;
        padding: 10px;
        margin: 2px;
    }
    button:first-child {
        margin-top: 50px;
    }
    .contents {
        display: flex;
        flex-direction: row;
        align-items: center;
        justify-content: center;
        width: 100%;
    }
    .buttonDiv {
        display: flex;
        flex-direction: row;
    }
    .upBtnDiv {
        display: flex;
        flex-direction: column;
    }
    .downBtnDiv {
        display: flex;
        flex-direction: column;
    }

	</style>
	<head>
	<meta charset = "UTF-8"/>
	<title>Mi Tracker de Series</title>
    ` + script + `

	</head>
	<body>
    <div class="container">
	<h1>Mi Tracker de Series</h1>
    <div class="contents">
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
    <div class="buttonDiv">
    <div class="upBtnDiv">
    ` + buttonsUp + ` 
    </div>
    <div class="downBtnDiv">
    ` + buttonsDown + `
    </div>
    </div>
    </div>
	<p><a href="./create">Agregar una Serie</a></p>
    </div>
	</body>
	</html>
	`

	return response + headers
}
