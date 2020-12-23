package transport

type MessageListener interface {
	HandleMessage([]byte, *Connection)
}
