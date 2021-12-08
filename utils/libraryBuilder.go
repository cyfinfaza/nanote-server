package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

type Media struct {
	MediaUrl string
	Title, Album, Artist, Genre string
	Year int
	CoverUrl string
}

var acceptedFileTypes = []string{
	".mp3",
	".ogg",
	".flac",
	".m4a",
	".wav",
	".wma",
	".aac",
}

func BuildLibraryRecursive(root string, path string) ([]Media, error) {
	directory := root + path
	output := []Media{}
	dir, err := os.ReadDir(directory)
	if err != nil {
		fmt.Println(err)
		return nil, err
	} else {
		// fmt.Println(dir)
		var defaultCover string;
		for _, file := range dir {
			if file.Name() == "cover.jpg" {
				defaultCover = "/content" + path + "/" + file.Name()
			}
		}
		for _, file := range dir {
			// fmt.Println(file.Name())
			// ensure the file is acceptable
			if file.IsDir() {
				items, err := BuildLibraryRecursive(root, path + "/" + file.Name())
				if err != nil {
					fmt.Println(err)
				} else {
					output = append(output, items...)
				}
			} else {
				if ! ContainsString(acceptedFileTypes, filepath.Ext(file.Name())) {
					continue
				}
				data, err := ReadMetadata(directory + "/" + file.Name())
				new := Media{ MediaUrl: "/content" + path + "/" + file.Name() }
				// if err is not null set the title to the path + file name, otherwise check the metadata
				if err != nil {
					new.Title = path + "/" + file.Name()
				} else {
					new.Title = data.Title
					new.Album = data.Album
					new.Artist = data.Artist
					new.Genre = data.Genre
					new.Year = data.Year
					if data.Picture != nil {
						new.CoverUrl = "/coverImage" + path + "/" + file.Name()
					} else {
						new.CoverUrl = defaultCover
					}
				}
				output = append(output, new)
			}
		}
		return output, nil
	}
}