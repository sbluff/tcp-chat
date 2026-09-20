package main

import (
	"bufio"
	"fmt"
	"net"
	"tcpchat/internal/room"
)

func handleConnectionClose(connection net.Conn, room *room.Room) {
	connection.Close()
	room.RemoveConnection(connection)
}

func handleConnection(connection net.Conn, room *room.Room) {
	defer handleConnectionClose(connection, room)

	room.AddConnection(connection)

	scanner := bufio.NewScanner(connection)

	for {
		isReadSuccesful := scanner.Scan()

		if !isReadSuccesful {
			fmt.Println("Read attempt failed")
			return
		}

		room.LogMessage(connection, scanner.Text())
	}
}

func main() {
	room := room.Room{
		Name: "localRoom",
	}

	room.Run()

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

		go handleConnection(connection, &room)
	}
}
