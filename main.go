package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

func handleConnection(connection net.Conn) {
	defer connection.Close()

	message := fmt.Sprintf(
		"Connection for local address %s is being handled",
		connection.LocalAddr(),
	)

	fmt.Println(message)

	scanner := bufio.NewScanner(connection)

	for {
		isReadSuccesful := scanner.Scan()

		if !isReadSuccesful {
			fmt.Println("Read attempt failed")
		}

		scannedText := []byte(scanner.Text())
		logText := fmt.Sprintf(
			"Connection %s: %s",
			connection.LocalAddr().String(),
			scannedText,
		)

		connection.Write([]byte(logText))
		log.Println(logText)
	}
}

func main() {
	fmt.Println("Start listening connections in port 4000")
	connectionListener, error := net.Listen("tcp", "localhost:4000")

	if error != nil {
		fmt.Println("connection listener could not be created")
		return
	}

	for {
		connection, error := connectionListener.Accept()

		if error != nil {
			fmt.Println("connection restarted")
			continue
		}

		handleConnection(connection)
	}
}
