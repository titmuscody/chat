package chat

import (
	//"fmt"
	"github.com/gorilla/websocket"
)


type ChatRoom struct {
	Users []Client
	SessionName string
}

type Client struct {
	Username string
	Conn *websocket.Conn
}