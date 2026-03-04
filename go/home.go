package main

import (
	"database/sql"
	"fmt"
	"log"
)

func handleHome() string {
	db, err := sql.Open("sqlite", "file:../series.db")
	if err != nil {
		log.Print("Error opening db:", err)
	}

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
	var deleteButtons string
	for _, s := range seriesList {
		status := ""
		if s.currentEp == s.episodes {
			status = "✅"
		}
		tableRowsString += `<tr id="row` + fmt.Sprintf("%d", s.id) + `">
		<td>` + fmt.Sprintf("%d", s.id) + `</td>
		<td id='name` + fmt.Sprintf("%d", s.id) + `'>` + s.name + status + `</td>
        <td id='ep` + fmt.Sprintf("%d", s.id) + `'>` + fmt.Sprintf("%d", s.currentEp) + `</td>
		<td id='tot` + fmt.Sprintf("%d", s.id) + `'>` + fmt.Sprintf("%d", s.episodes) + `</td>
		<td> <progress id='p` + fmt.Sprintf("%d", s.id) + `' value='` + fmt.Sprintf("%d", s.currentEp) + `' max='` + fmt.Sprintf("%d", s.episodes) + `'></progress>
		</tr>`

		buttonsDown += `<button id="dec` + fmt.Sprintf("%d", s.id) + `" onclick="prevEpisode(` + fmt.Sprintf("%d", s.id) + `)">-1</button>`
		buttonsUp += `<button id="inc` + fmt.Sprintf("%d", s.id) + `" onclick="nextEpisode(` + fmt.Sprintf("%d", s.id) + `)">+1</button>`
		deleteButtons += `<button id="del` + fmt.Sprintf("%d", s.id) + `" onclick="deleteSeries(` + fmt.Sprintf("%d", s.id) + `)">Eliminar</button>`
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

            const nameElement = document.getElementById("name" + id);
            if ((current + 1) == totalEp && !nameElement.innerText.includes("✅")){
                nameElement.innerText = nameElement.innerText + "✅";
            }
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

            const nameElement = document.getElementById("name" + id);
            if ((current - 1) < totalEp && nameElement.innerText.includes("✅")){
                nameElement.innerText = nameElement.innerText.replaceAll("✅", "").trim();
            }
        }
    }
    window.deleteSeries = async function deleteEpisode(id) {
        const deleteConfirm = confirm("Estás seguro de que quieres eliminar esa serie?")

        if (deleteConfirm){
            const url = "/delete/?id=" + id;
            const response = await fetch(url, { method: "DELETE" } );

            const seriesRowElement = document.getElementById("row" + id);
            const incrementElement = document.getElementById("inc" + id);
            const decrementElement = document.getElementById("dec" + id);
            const deleteElement = document.getElementById("del" + id);

            seriesRowElement.remove();
            incrementElement.remove();
            decrementElement.remove();
            deleteElement.remove();
        }
    }
    </script>`

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
    .deleteBtnDiv {
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
    <div class="deleteBtnDiv">
    ` + deleteButtons + `
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
