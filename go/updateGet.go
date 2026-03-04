package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"strconv"
)

func handleUpdateGET(query url.Values) string {
	idStr := query.Get("id")
	id, _ := strconv.Atoi(idStr)

	var serie Serie

	db, err := sql.Open("sqlite", "file:../series.db")
	if err != nil {
		log.Print("Error opening DB:", err)
	}
	defer db.Close()

	row := db.QueryRow("SELECT * FROM series WHERE id=?", id)
	row.Scan(&serie.id, &serie.name, &serie.currentEp, &serie.episodes)

	response := "HTTP/1.1 200 OK\r\n" +
		"Content-Type: text/html\r\n" +
		"\r\n" +
		`
        <!DOCTYPE html>
        <html lang="en">
            <head>
                <meta charset="UTF-8"/>
                <title>Editar una Serie</title>
            </head>

            <body>
                <style>
                    body > * {
                        margin: 0px;
                        padding: 0px;
                        font-family: tahoma;
                        font-size: 24px;
                    }
                    .container{
                        width: 100vw;
                        height: 100vh;
                        display: flex;
                        flex-direction: column;
                        justify-content: center;
                        align-items: center;
                    }
                    h1 {
                        font-size: 32px;
                    }
                    form > p:not(:first-child){
                        margin-top: 30px;
                    }
                    p {
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
                    <h1>Editando serie: ` + serie.name + `</h1>
                    <p>Información actual: ID: ` + fmt.Sprintf("%d", serie.id) + ` | Episodio Actual: ` + fmt.Sprintf("%d", serie.currentEp) + ` | Episodios Totales: ` + fmt.Sprintf("%d", serie.episodes) + `</p>
                    <form id="editar" method="POST" action="/update/">
                    <p>ID de la serie</p>
                    <input type="number" name="id" min="0" value="` + fmt.Sprintf("%d", serie.id) + `" required>
                    <p>Nombre nuevo</p>
                    <input type="text" name="series_name" value="` + serie.name + `" required>
                    <p>Episodio actual</p>
                    <input type="number" name="current_episode" min="1" value="` + fmt.Sprintf("%d", serie.currentEp) + `" required>
                    <p>Total de Episodios</p>
                    <input type="number" name="total_episodes" min="1" value="` + fmt.Sprintf("%d", serie.episodes) + `" required>
                    <input type="hidden" name="original_id" value="` + fmt.Sprintf("%d", serie.id) + `">
                    </form>
                    <button type="submit" form="editar">Guardar Cambios</button>
                    <p class="back"><a href="../">Regresar</a></p> 
                </div>
            </body>
        </html>
        `
	return response
}
