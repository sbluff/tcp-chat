package room

import (
	"errors"
	"fmt"
	"net"
	"time"
)

type Room struct {
	Name        string
	connections []net.Conn
	messages    []RoomMessage
}

func (room *Room) AddConnection(connection net.Conn) error {
	if room.isConnectionInRoom(connection) {
		return errors.New("Cant add a connection that is already in the room")
	}

	updatedConnections := append(room.connections, connection)
	room.connections = updatedConnections

	return nil
}

func (room *Room) RemoveConnection(connection net.Conn) error {
	if !room.isConnectionInRoom(connection) {
		return errors.New("Cant remove a connection that is not in the room")
	}

	updatedConnections := append(room.connections, connection)
	room.connections = updatedConnections

	return nil
}

func (room *Room) isConnectionInRoom(connection net.Conn) bool {
	for i := range room.connections {
		localConnection := room.connections[i]

		if localConnection.LocalAddr().String() != connection.LocalAddr().String() {
			continue
		}

		return true
	}

	return false
}

type RoomMessage struct {
	creatorConnection net.Conn
	createdAt         time.Time
	content           string
}

func (room *Room) LogMessage(connection net.Conn, logMessage string) error {
	if !room.isConnectionInRoom(connection) {
		return errors.New("Cant log a message for a connection not in the room")
	}

	timeStamp := time.Now()

	roomMessage := RoomMessage{
		creatorConnection: connection,
		createdAt:         timeStamp,
		content:           logMessage,
	}

	logText := fmt.Sprintf(
		"[%s] %s: %s",
		timeStamp.Format(time.RFC822),
		connection.LocalAddr().String(),
		logMessage,
	)

	fmt.Println(logText)

	room.messages = append(room.messages, roomMessage)
	return nil
}

func main() {
}
