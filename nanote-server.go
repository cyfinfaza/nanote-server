package main

import (
	"encoding/json"
	"fmt"
	"nanote-server/utils"
	"net/http"
	"path/filepath"
	"strings"
)

var library []utils.Media
var config utils.Config

func httpHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Server", "nanote")
	operation := strings.Split(r.URL.Path, "/")[1]
	path := strings.Join(strings.Split(r.URL.Path, "/")[2:], "/")
	fmt.Println(operation)
	switch operation {
	case "users":
		type user struct {
			User string `json:"user"`
			Name string `json:"name"`
		}
		var users []user
		for key, data := range config.Users {
			users = append(users, user{User: key, Name: data.Name})
		}
		json.NewEncoder(w).Encode(users)
		return
	case "userImg":
		if config.Users[path].Picture != "" {
			http.ServeFile(w, r, config.Users[path].Picture)
		} else {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "Not found")
		}
		return
	}
	username, password, ok := r.BasicAuth()
	if !ok || config.Users[username].Key != password {
		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "Unauthorized")
		return
	}
	mediaRoot := config.Users[username].MediaRoot
	switch operation {
	case "test":
		fmt.Fprint(w, "OK")
	case "coverImage":
		data, err := utils.ReadMetadata(filepath.Join(mediaRoot, path))
		if err != nil || data.Picture == nil {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, err)
			return
		} else {
			w.Header().Set("Content-Type", data.Picture.MIMEType)
			w.Write(data.Picture.Data)
			return
		}
	case "content":
		http.ServeFile(w, r, filepath.Join(mediaRoot, path))
	case "library":
		output, _ := utils.BuildLibraryRecursive(mediaRoot, "")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(output)
	default:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "Not found")
	}
}

func main() {
	readConfig, err := utils.ReadConfig("./config.yml")
	if err != nil {
		panic(err)
	} else {
		config = readConfig
		fmt.Println("Config loaded")
	}
	fmt.Println("Serving")
	server := http.Server{
		Addr:    ":" + string(config.App.Port),
		Handler: http.HandlerFunc(httpHandler),
	}
	server.ListenAndServe()
	// data, _ := utils.ReadMetadata("audio.flac")
	// gob.NewEncoder(os.Stdout).Encode(data)
	// fmt.Println(output)
	// file, _ := os.OpenFile("cache.gob", os.O_RDWR|os.O_CREATE, 0666)
	// gob.NewEncoder(file).Encode(output)
}
