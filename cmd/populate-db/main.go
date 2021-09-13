// Program to populate sqlite DB with track data, given a music folder

package main

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	"github.com/dhowden/tag"
	"github.com/nephila-nacrea/rank-my-music/repo"
	"github.com/nephila-nacrea/rank-my-music/track"

	_ "modernc.org/sqlite"
)

func init() {
	log.SetFlags(log.Llongfile)
}

func main() {
	folderPath := os.Args[1]

	// TODO
	// 	Handle duplicates
	// 	What if no metadata? (e.g. wma)

	var tracks []track.Track

	var files []string
	err := filepath.Walk(
		folderPath,
		func(path string, info os.FileInfo, err error) error {
			if !info.IsDir() {
				files = append(files, path)
			}
			return nil
		},
	)
	if err != nil {
		log.Fatalln(err)
	}

	for _, f := range files {
		file, err := os.Open(f)
		if err != nil {
			log.Println(f+": ", err)
		}

		meta, err := tag.ReadFrom(file)
		if err != nil {
			log.Println(f+": ", err)
		} else {
			// Dedupe artist data
			artists := []string{meta.Artist()}

			if meta.AlbumArtist() != meta.Artist() {
				artists = append(artists, meta.AlbumArtist())
			}

			if meta.Composer() != meta.Artist() &&
				meta.Composer() != meta.AlbumArtist() {
				artists = append(artists, meta.Composer())
			}

			tracks = append(
				tracks,
				track.New(
					meta.Title(),
					meta.Album(),
					artists,
				),
			)
		}
	}

	log.Println("Now for the database!")

	db, err := sql.Open("sqlite", "file:ranked_music.sqlt")
	if err != nil {
		log.Fatalln(err)
	}

	repo.SaveTracks(db, tracks)
}
