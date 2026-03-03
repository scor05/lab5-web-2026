package main

import (
	"database/sql"
	// "fmt"
	// "log"
)

func handleCreate() string {
	db, _ := sql.Open("sqlite", "file:series.db")
	defer db.Close()

	response := `
	HTTP/1.1 200 OK\r\n
	Content-Type: text/html\r\n
	\r\n
	<!DOCTYPE html>
	<html>
		<head>
			<meta charset = "UTF-8"/>
			<title>Agregar serie nueva</title>
			<script src="./create.js"></script>
		</head>

		<body>
			<style>
				body {
					display: flex;
					flex-direction: column;
					justify-content: space-between;
					align-items: center;
					font-family: Tahoma;
				}
				p {
					font-size: 24;
				}
				input {
					font-size: 18;
				button {
					background-color: cadetBlue;
					border-radius: 5px;
					font-size 24;
				}
			</style>
			<form method="POST" action="/create">
					<p>Nombre de la serie</p>
					<input type="text" name="series_name" required>
					<p>Episodio actual</p>
					<input type="number" name="current_episode" min="1" value="1" required>
					<p>Total de Episodios</p>
					<input type="number" name="total_episodes" min="1" required>
				<button type="submit">Crear</button>
			</form>
		</body>
	</html>`
	return response
}
