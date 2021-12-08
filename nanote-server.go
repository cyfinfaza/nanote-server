package main

import (
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"nanote-server/utils"
	"net/http"
	"os"
	"strings"
)

var library []utils.Media
const mediaRoot = "./Unlike Pluto/"

func httpHandler(w http.ResponseWriter, r *http.Request) {
	operation := strings.Split(r.URL.Path, "/")[1]
	path := strings.Join(strings.Split(r.URL.Path, "/")[2:], "/")
	fmt.Println(operation)
	switch operation {
	case "coverImage":
		data, err := utils.ReadMetadata(mediaRoot + path)
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
		// open the specified file and stream it to the client
		file, err := os.OpenFile(mediaRoot + path, os.O_RDONLY, 0666)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, err)
			return
		} else {
			defer file.Close()
			w.Header().Set("Content-Type", mime.TypeByExtension("."+strings.Split(path, ".")[len(strings.Split(path, "."))-1]))
			io.Copy(w, file)
		}
	case "library":
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(library)
	default:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "Not found")
	}
}

func main() {
	fmt.Println("Building library...")
	output, _ := utils.BuildLibraryRecursive(mediaRoot, "")
	library = output
	fmt.Println("Library built.")
	fmt.Println("Serving")
	server := http.Server{
		Addr: ":8080",
		Handler: http.HandlerFunc(httpHandler),
	}
	server.ListenAndServe()
	// data, _ := utils.ReadMetadata("audio.flac")
	// gob.NewEncoder(os.Stdout).Encode(data)
	// fmt.Println(output)
	// file, _ := os.OpenFile("cache.gob", os.O_RDWR|os.O_CREATE, 0666)
	// gob.NewEncoder(file).Encode(output)
}