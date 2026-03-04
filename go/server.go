package main

import (
	"bufio"
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

	switch {
	case path == "/" && method == "GET":
		response = handleHome()
	case path == "/create" && method == "GET":
		response = handleCreate()

	case path == "/create" && method == "POST" && contentLength > 0:
		bytesBody := make([]byte, contentLength)
		_, err := io.ReadFull(reader, bytesBody)
		if err != nil {
			log.Print("Error leyendo body: ", err)
			response = "HTTP/1.1 400 Bad Request\r\nConnection: close\r\n\r\n"
		}
		body := string(bytesBody)

		response = handleCreatePOST(body)

	case path == "/update/" && method == "POST":
		query, err := url.ParseQuery(rawQuery)
		if err != nil {
			log.Print("Error parsing query:", err)
		}

		// el método POST lo puse para para incrementar/decrementar y para actualizar series
		// para diferenciar entre ambas está este query
		temp := query.Get("change")
		if temp == "p" || temp == "m" {
			response = handleUpdatePOST(query)
		} else {
			if contentLength <= 0 {
				response = "HTTP/1.1 400 Bad Request\r\nConnection: close\r\n\r\n"
				break
			}

			bytesBody := make([]byte, contentLength)
			_, err := io.ReadFull(reader, bytesBody)
			if err != nil {
				log.Print("Error leyendo body: ", err)
				response = "HTTP/1.1 400 Bad Request\r\nConnection: close\r\n\r\n"
				break
			}

			body := string(bytesBody)
			formVals, err := url.ParseQuery(body)
			if err != nil {
				log.Print("Error parsing body:", err)
				response = "HTTP/1.1 400 Bad Request\r\nConnection: close\r\n\r\n"
				break
			}
			response = handleUpdateUPDATE(formVals)
		}

	case path == "/delete/" && method == "DELETE":
		query, err := url.ParseQuery(rawQuery)
		if err != nil {
			log.Print("Error parsing query:", err)
		}

		response = handleDelete(query)

	case path == "/update/" && method == "GET":
		query, err := url.ParseQuery(rawQuery)
		if err != nil {
			log.Print("Error parsing query:", err)
		}

		response = handleUpdateGET(query)
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
