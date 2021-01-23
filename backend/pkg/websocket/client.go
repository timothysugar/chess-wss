package websocket

import (
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	Conn *websocket.Conn
}

type Message struct {
	Type int    `json:"type"`
	Body string `json:"body"`
}

func (c *Client) Read() chan Message {
	// defer func() {
	// 	c.Conn.Close()
	// }()

	ch := make(chan Message)
	go func() {
		for {
			messageType, p, err := c.Conn.ReadMessage()
			if err != nil {
				log.Println(err)
			}

			message := Message{Type: messageType, Body: string(p)}
			fmt.Printf("Message Received: %+v\n", message)
			ch <- message
			fmt.Printf("Message sent to channel: %+v\n", message)
		}
	}()

	return ch
}

func (c *Client) Write(body string) error {
	message := Message{
		Type: 0,
		Body: body,
	}
	fmt.Printf("Sending message: %+v %+v\n", *c, message)

	return c.Conn.WriteJSON(message)
}
