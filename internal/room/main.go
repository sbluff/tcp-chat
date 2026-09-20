package room

import (
	"errors"
	"fmt"
	"net"
	"time"
)

type Room struct {
	Name                    string
	connections             []RoomConnection
	messages                []RoomMessage
	addConnectionChannel    chan net.Conn
	removeConnectionChannel chan net.Conn
}

func (room *Room) AddConnection(connection net.Conn) {
	room.addConnectionChannel <- connection
}

func (room *Room) RemoveConnection(connection net.Conn) {
	room.removeConnectionChannel <- connection
}

func (room *Room) addConnection(connection net.Conn) error {
	if room.isConnectionInRoom(connection) {
		return errors.New("Cant add a connection that is already in the room")
	}

	connectionToAdd := RoomConnection{
		Connection:  connection,
		RoomMessage: make(chan RoomMessage),
	}

	go connectionToAdd.Run()

	updatedConnections := append(room.connections, connectionToAdd)
	room.connections = updatedConnections

	return nil
}

func (room *Room) removeConnection(connection net.Conn) error {
	if !room.isConnectionInRoom(connection) {
		return errors.New("Cant remove a connection that is not in the room")
	}

	var updatedConnections []RoomConnection

	for i := range room.connections {
		roomConnection := room.connections[i]

		if roomConnection.GetRemoteAddress() != connection.RemoteAddr().String() {
			updatedConnections = append(updatedConnections, roomConnection)
		}
	}

	room.connections = updatedConnections

	return nil
}

func (room *Room) Run() {
	room.addConnectionChannel = make(chan net.Conn)
	room.removeConnectionChannel = make(chan net.Conn)

	go func() {
		for {
			select {
			case channelAddConnectionData := <-room.addConnectionChannel:
				room.addConnection(channelAddConnectionData)
			case channelRemoveConnectionData := <-room.removeConnectionChannel:
				room.removeConnection(channelRemoveConnectionData)
			}
		}
	}()
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

	fmt.Println(roomMessage.LogContents())

	room.messages = append(room.messages, roomMessage)
	room.updateListeners(roomMessage)

	return nil
}

func (room *Room) isConnectionInRoom(connection net.Conn) bool {
	for i := range room.connections {
		localConnection := room.connections[i]

		if localConnection.GetRemoteAddress() == connection.RemoteAddr().String() {
			return true
		}
	}

	return false
}

func (room *Room) updateListeners(roomMessage RoomMessage) {
	for i := range room.connections {
		roomConnection := room.connections[i]

		if roomConnection.GetRemoteAddress() == roomMessage.creatorConnection.RemoteAddr().String() {
			continue
		}

		roomConnection.RoomMessage <- roomMessage
	}
}
