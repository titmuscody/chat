package chat

import (
	"fmt"
	"github.com/gorilla/websocket"
)


type ChatRoom struct {
	Users []*Client
	SessionName string
	receive chan string
}

func (cr *ChatRoom)Start() {
	cr.receive = make(chan string, 1)
	for {
		fmt.Println("pulling string from rec: ")
		mess := <- cr.receive
		fmt.Println("size of list=", len(cr.Users))
		fmt.Println("done pulling", mess)
		//fmt.Println("pulling string from rec: ", mess)
		for i := range cr.Users {
			err := cr.Users[i].Conn.WriteMessage(1, []byte(mess))
			if err != nil {
				fmt.Println("error writing to socket", err)
			}
		}
	}
}
func (cr *ChatRoom)StartClient(client *Client) {
	cr.Users = append(cr.Users, client)
	
	for {
	messType, mess, err := client.Conn.ReadMessage()
	if err != nil {
		fmt.Println(err)
	}
	if messType == 1 {
	message := client.Username + ":" + string(mess)
	fmt.Println("sending to rec", message)
	cr.receive <-message
	fmt.Println("done sending")
	} else if messType == -1 {
		// todo: remove from list
		client.Conn.Close()
	}
	}
	
}

type Client struct {
	Username string
	Conn *websocket.Conn
}

