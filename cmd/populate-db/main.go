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

const problemTracksFilename = "problem-tracks.txt"

func init() {
	log.SetFlags(log.Llongfile)
}

func main() {
	folderPath := os.Args[1]

	// TODO
	// 	Handle duplicates
	// 	What if no metadata? (e.g. wma)

	var tracks []track.Track

	var filenames []string
	err := filepath.Walk(
		folderPath,
		func(path string, info os.FileInfo, err error) error {
			if info != nil && !info.IsDir() {
				filenames = append(filenames, path)
			}
			return nil
		},
	)
	if err != nil {
		log.Fatalln(err)
	}

	problemTracksFile, err := os.OpenFile(
		problemTracksFilename,
		os.O_WRONLY,
		os.ModeAppend,
	)
	if err != nil {
		log.Println(err)

		problemTracksFile, err = os.Create(problemTracksFilename)
		if err != nil {
			panic(err)
		}
	}

	defer problemTracksFile.Close()

	for _, filename := range filenames {
		file, err := os.Open(filename)
		if err != nil {
			log.Println(filename+": ", err)
		}

		meta, err := tag.ReadFrom(file)
		if err != nil {
			log.Println(filename+": ", err)

			_, err = problemTracksFile.WriteString(filename + "\n")
			if err != nil {
				panic(err)
			}
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
					"", // TODO musicbrainz_id
				),
			)

			raw := meta.Raw()
			// log.Println(meta.Title())
			// log.Println(meta.Album())
			// log.Println(meta.Format())

			// For format = VORBIS
			// log.Println(raw["musicbrainz_trackid"])

			// For format = ID3v2.3
			// TODO
			// In format 'http://musicbrainz.org (*)',
			// need to get * only
			// log.Println(raw["UFID"])

			// For format = MP4
			// log.Println(raw["MusicBrainz Track Id"])

			// log.Println("==========")

			if raw["musicbrainz_trackid"] == nil &&
				raw["UFID"] == nil &&
				raw["MusicBrainz Track Id"] == nil {
				problemTracksFile.WriteString(filename + "\n")
				problemTracksFile.WriteString("    Title: " + meta.Title() + "\n")
				problemTracksFile.WriteString("    Album: " + meta.Album() + "\n")
				problemTracksFile.WriteString("    Artist: " + meta.Artist() + "\n")
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
