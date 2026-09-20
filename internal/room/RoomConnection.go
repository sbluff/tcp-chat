package room

import (
	"net"
)

type RoomConnection struct {
	Connection  net.Conn
	RoomMessage chan RoomMessage
}

func (roomConnection RoomConnection) GetRemoteAddress() string {
	return roomConnection.Connection.RemoteAddr().String()
}

func (roomConnection *RoomConnection) Run() {
	for {
		message := <-roomConnection.RoomMessage

		logContents := []byte(message.LogContents())

		roomConnection.Connection.Write(logContents)
	}
}
