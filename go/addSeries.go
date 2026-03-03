package main

import (
	"database/sql"
	"log"
)

func handleCreate() string {
	db, err := sql.Open("sqlite", "file:../series.db")
	if err != nil {
		log.Print("Error opening db: ", err)
	}
	defer db.Close()

	response := "HTTP/1.1 200 OK\r\n" +
		"Content-Type: text/html\r\n" +
		"\r\n" +
		`
	<!DOCTYPE html>
	<html>
		<head>
			<meta charset = "UTF-8"/>
			<title>Agregar serie nueva</title>
		</head>

		<body>
			<style>
				.container {
					margin: 0px;
					width: 100%;
					height: 100vh;
					display: flex;
					flex-direction: column;
					justify-content: center;
					align-items: center;
					font-family: Tahoma;
				}
				form > p:not(:first-child){
					margin-top: 30px;
				}
				p {
					font-size: 24px;
					margin: 0px;
					padding: 10px;
				}
				.back {
					padding-top: 100px;
				}
				input {
					font-size: 18px;
					border-radius: 5px;
					border-color: cadetBlue;
					padding: 5px;
				}
				button {
					background-color: cadetBlue;
					border-radius: 10px;
					font-size: 24px;
					padding: 10px;
                    margin-top: 20px;
				}
			</style>
			<div class="container">
				<h1>Agregar una Serie</h1>
				<form id="crear" method="POST" action="/create">
					<p>Nombre de la serie</p>
					<input type="text" name="series_name" required>
					<p>Episodio actual</p>
					<input type="number" name="current_episode" min="1" value="1" required>
					<p>Total de Episodios</p>
					<input type="number" name="total_episodes" min="1" required>
				</form>
				<button type="submit" form="crear">Crear</button>
				<p class="back"><a href="../">Regresar</a></p>
			</div>
		</body>
	</html>`
	return response
}
