package main

import (
	"bufio"
	"log"
	_ "modernc.org/sqlite"
	"net"
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

	// chuparse los headers
	for {
		headerLine, _ := reader.ReadString('\n')
		if headerLine == "\r\n" {
			break
		}
	}

	response := ""
	if path == "/" {
		response = handleHome()
	}
	if path == "/create" && method == "GET" {
		response = handleCreate()
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
