package utils

import (
	"os"

	"github.com/dhowden/tag"
)

type MetadataObject struct {
	Title, Album, Artist, Genre string
	Year                        int
	Picture                     *tag.Picture
}

func ReadMetadata(filename string) (MetadataObject, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY, 0666)
	defer file.Close()
	if err != nil {
		return MetadataObject{}, err
	} else {
		m, err := tag.ReadFrom(file)
		if err != nil {
			return MetadataObject{}, err
		} else {
			return MetadataObject{
				Title:   m.Title(),
				Album:   m.Album(),
				Artist:  m.Artist(),
				Genre:   m.Genre(),
				Year:    m.Year(),
				Picture: m.Picture(),
			}, nil
		}
	}
}
