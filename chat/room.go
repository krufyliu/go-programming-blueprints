package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/krufyliu/go-programming-blueprints/trace"
)

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

type room struct {
	forward chan []byte
	join    chan *client
	leave   chan *client
	clients map[*client]bool
	tracer  trace.Tracer
}

func newRoom() *room {
	return &room{
		forward: make(chan []byte),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
		tracer:  trace.Off(),
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
			r.tracer.Trace("received message: ", string(msg))
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
	client := &client{
		socket: socket,
		send:   make(chan []byte, messageBufferSize),
		room:   r,
	}
	r.join <- client
	defer func() { r.leave <- client }()
	go client.Write()
	client.Read()
}
