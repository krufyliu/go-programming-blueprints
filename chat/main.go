package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/gorilla/handlers"
)

type templateHandler struct {
	once     sync.Once
	filename string
	tpl      *template.Template
}

func (h *templateHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	h.once.Do(func() {
		h.tpl = template.Must(template.ParseFiles(filepath.Join("templates", h.filename)))
	})
	h.tpl.Execute(w, req)
}

func main() {
	var addr = flag.String("addr", ":8888", "The addr of the application")
	flag.Parse()
	var r = newRoom()
	// r.tracer = trace.New(os.Stdout)
	http.Handle("/", &templateHandler{filename: "chat.html"})
	http.Handle("/room", r)
	go r.run()
	log.Println("Starting web server on ", *addr)
	if err := http.ListenAndServe(*addr, handlers.LoggingHandler(os.Stdout, http.DefaultServeMux)); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
