package room

import (
	"net"
)

type RoomConnection struct {
	Connection         net.Conn
	RoomMessageChannel chan RoomMessage
}

func (roomConnection RoomConnection) GetRemoteAddress() string {
	return roomConnection.Connection.RemoteAddr().String()
}

func (roomConnection *RoomConnection) Run() {
	for {
		message, ok := <-roomConnection.RoomMessageChannel

		if !ok {
			return
		}

		logContents := []byte(message.LogContents())

		roomConnection.Connection.Write(logContents)
	}
}
