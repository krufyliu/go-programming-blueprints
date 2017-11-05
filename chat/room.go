package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/krufyliu/go-programming-blueprints/trace"
	"github.com/stretchr/objx"
)

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

type room struct {
	forward chan *message
	join    chan *client
	leave   chan *client
	clients map[*client]bool
	avatar  Avatar
	tracer  trace.Tracer
}

func newRoom(avatar Avatar) *room {
	return &room{
		forward: make(chan *message),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
		tracer:  trace.Off(),
		avatar:  avatar,
	}
}

func (r *room) run() {
	for {
		select {
		case c := <-r.join:
			r.clients[c] = true
			r.tracer.Trace("client joined")
		case c := <-r.leave:
			delete(r.clients, c)
			close(c.send)
			r.tracer.Trace("client left")
		case msg := <-r.forward:
			r.tracer.Trace("received message: ", msg.Message)
			for c := range r.clients {
				c.send <- msg
				r.tracer.Trace(" -- send message")
			}
		}
	}
}

var upgrader = &websocket.Upgrader{
	ReadBufferSize:  socketBufferSize,
	WriteBufferSize: socketBufferSize,
}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServeHTTP: ", err)
		return
	}
	authCookie, err := req.Cookie("auth")
	if err != nil {
		log.Fatal("can not get auth cookie")
		return
	}
	client := &client{
		socket:   socket,
		send:     make(chan *message, messageBufferSize),
		room:     r,
		userData: objx.MustFromBase64(authCookie.Value),
	}
	r.join <- client
	defer func() { r.leave <- client }()
	go client.Write()
	client.Read()
}
