// Program to populate sqlite DB with track data, given a music folder

package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/dhowden/tag"
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
			if info != nil && !info.IsDir() {
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
			artists := []string{}

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
					meta.Artist(),
					artists,
					"", // TODO
				),
			)

			raw := meta.Raw()
			log.Println(meta.Title())
			log.Println(meta.Album())
			log.Println(meta.Format())

			// // For format = VORBIS
			log.Println(raw["musicbrainz_trackid"])

			// For format = ID3v2.3
			// TODO
			// In format 'http://musicbrainz.org (*)',
			// need to get * only
			log.Println(raw["UFID"])

			// For format = MP4
			log.Println(raw["MusicBrainz Track Id"])

			log.Println("==========")

			if raw["musicbrainz_trackid"] == nil &&
				raw["UFID"] == nil &&
				raw["MusicBrainz Track Id"] == nil {
				// log.Println(f)
				// log.Printf("%#v", raw)
				// log.Println(raw["UFID"])
				// log.Println(meta.Title())
				// log.Println(meta.Album())
				// log.Println(meta.Format())
				// TODO Store in a file
			}
		}
	}

	log.Println("Now for the database!")

	// db, err := sql.Open("sqlite", "file:ranked_music.sqlt")
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// repo.SaveTracks(db, tracks)
}
