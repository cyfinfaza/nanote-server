package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"nanote-server/utils"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var config utils.Config

var cacheLocks map[string]*sync.Mutex

func buildLibraryCache(user string) error {
	cacheLocks[user].Lock()
	defer cacheLocks[user].Unlock()
	output, _ := utils.BuildLibraryRecursive(config.Users[user].MediaRoot, "/")
	os.Mkdir("nanote-library-cache", 0666)
	file, err := os.Create(filepath.Join("./nanote-library-cache", user+".json"))
	defer file.Close()
	if err != nil {
		return err
	} else {
		json.NewEncoder(file).Encode(output)
		return nil
	}
}

func httpHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Server", "nanote")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	if r.Method == "OPTIONS" {
		fmt.Fprintf(w, "OK")
	}
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
	username, password, hasBasicAuth := r.BasicAuth()
	authString, hasAuthKey := r.URL.Query()["auth"]
	if hasAuthKey {
		authString, err := base64.URLEncoding.DecodeString(authString[0])
		if err != nil {
			hasAuthKey = false
		} else {
			auth := strings.Split(string(authString), ":")
			fmt.Println(auth)
			username = auth[0]
			password = auth[1]
		}
	}
	if (!hasBasicAuth && !hasAuthKey) || config.Users[username].Key != password {
		if hasBasicAuth {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
		}
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "Unauthorized")
		return
	}
	mediaRoot := config.Users[username].MediaRoot
	userConfig := config.Users[username]
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
		if userConfig.CacheEnabled {
			cacheLocks[username].Lock()
			defer cacheLocks[username].Unlock()
			file, err := os.Open(filepath.Join("./nanote-library-cache", username+".json"))
			defer file.Close()
			if err == nil {
				data, _ := ioutil.ReadAll(file)
				var library []utils.Media
				err = json.Unmarshal(data, &library)
				if err == nil {
					w.Header().Set("nanote-cache", "hit")
					w.Header().Set("Content-Type", "application/json")
					json.NewEncoder(w).Encode(library)
					return
				}
			}
			w.Header().Set("nanote-cache", "fail")
		} else {
			w.Header().Set("nanote-cache", "disabled")
		}
		output, _ := utils.BuildLibraryRecursive(mediaRoot, "")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(output)
	case "rebuildLibrary":
		if userConfig.CacheEnabled {
			fmt.Println("Rebuilding library cache for user", username)
			err := buildLibraryCache(username)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprint(w, err)
			} else {
				fmt.Fprint(w, "Success")
			}
		} else {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "Cache is not enabled for this user")
		}
	default:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "Not found")
	}
}

func main() {
	configFile := "./config.yml"
	if len(os.Args) > 1 {
		configFile = os.Args[1]
	}
	fmt.Println(configFile)
	readConfig, err := utils.ReadConfig(configFile)
	if err != nil {
		panic(err)
	} else {
		config = readConfig
		fmt.Println("Config loaded")
		// fmt.Println(config)
	}
	fmt.Println("Building library caches")
	cacheLocks = make(map[string]*sync.Mutex)
	for username, user := range config.Users {
		if user.CacheEnabled {
			cacheLocks[username] = &sync.Mutex{}
			fmt.Println("Building cache for user", username)
			err := buildLibraryCache(username)
			if err != nil {
				fmt.Println("Failed to create library cache for ", user, ", Error: ", err)
			}
		}
	}
	fmt.Println("Serving")
	server := http.Server{
		Addr:    ":" + string(config.App.Port),
		Handler: http.HandlerFunc(httpHandler),
	}
	server.ListenAndServe()
}
