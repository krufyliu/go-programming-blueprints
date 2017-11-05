package main

import (
	"io"
	"io/ioutil"
	"net/http"
	"path"
)

const MaxUploadSize = 5 * (1 << 20)

type Uploader struct {
}

func (Uploader) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		th := &templateHandler{filename: "upload.html"}
		th.ServeHTTP(w, r)
		return
	}
	if r.Method == http.MethodPost {
		r.ParseMultipartForm(MaxUploadSize)
		userID := r.FormValue("userid")
		file, header, err := r.FormFile("avatarFile")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data, err := ioutil.ReadAll(file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		filename := path.Join("avatars", userID+path.Ext(header.Filename))
		err = ioutil.WriteFile(filename, data, 0755)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		io.WriteString(w, "Successful")
		return
	}
	http.NotFound(w, r)
}
