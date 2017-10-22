package main

import (
	"time"

	"github.com/gorilla/websocket"
)

type client struct {
	socket   *websocket.Conn
	send     chan *message
	room     *room
	userData map[string]interface{}
}

func (c *client) Read() {
	defer c.socket.Close()

	for {
		msg := new(message)
		err := c.socket.ReadJSON(msg)
		if err != nil {
			return
		}
		msg.From = c.userData["Name"].(string)
		msg.When = time.Now().Unix()
		if avatarURL, ok := c.userData["AvatarURL"]; ok {
			msg.AvatarURL = avatarURL.(string)
		}
		c.room.forward <- msg
	}
}

func (c *client) Write() {
	defer c.socket.Close()

	for msg := range c.send {
		err := c.socket.WriteJSON(msg)
		if err != nil {
			return
		}
	}
}
