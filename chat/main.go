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
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/facebook"
	"github.com/stretchr/gomniauth/providers/github"
	"github.com/stretchr/gomniauth/providers/google"
	"github.com/stretchr/objx"
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
	data := map[string]interface{}{
		"Host": req.Host,
	}
	if authCookie, err := req.Cookie("auth"); err == nil {
		data["UserData"] = objx.MustFromBase64(authCookie.Value)
	}
	h.tpl.Execute(w, data)
}

func main() {
	gomniauth.SetSecurityKey("EDGE")
	gomniauth.WithProviders(
		github.New("a547d64960af094dd75e", "008164004cb3be53377f84f8000a624af74073fa",
			"http://localhost:8888/auth/callback/github"),
		google.New("", "", "http://localhost:8888/auth/callback/google"),
		facebook.New("", "", "http://localhost:8888/auth/callback/facebook"),
	)
	var addr = flag.String("addr", ":8888", "The addr of the application")
	flag.Parse()
	var r = newRoom()
	// r.tracer = trace.New(os.Stdout)
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/room", r)
	http.HandleFunc("/auth/", loginHandler)
	go r.run()
	log.Println("Starting web server on ", *addr)
	if err := http.ListenAndServe(*addr, handlers.LoggingHandler(os.Stdout, http.DefaultServeMux)); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
