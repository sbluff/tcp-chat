package room

import (
	"fmt"
	"net"
	"time"
)

type RoomMessage struct {
	creatorConnection net.Conn
	createdAt         time.Time
	content           string
}

func (roomMessage RoomMessage) LogContents() string {
	return fmt.Sprintf(
		"> [%s] %s: %s \n",
		roomMessage.createdAt.Format(time.RFC822),
		roomMessage.creatorConnection.RemoteAddr().String(),
		roomMessage.content,
	)
}
