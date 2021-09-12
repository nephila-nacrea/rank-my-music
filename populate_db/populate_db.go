// Program to populate sqlite DB with track data, given a music folder

package main

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"example.com/rank-my-music/track"
	"github.com/dhowden/tag"

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

	log.Println("Database connected.")

	rows, err := db.Query(
		"SELECT name FROM sqlite_master WHERE type = 'table'")
	if err != nil {
		log.Fatalln(err)
	}

	var name string
	for rows.Next() {
		if err = rows.Scan(&name); err != nil {
			log.Fatalln(err)
		}
		log.Println(name)
	}

	for _, tr := range tracks {
		// TODO Transaction

		res, err := db.Exec(
			`INSERT INTO tracks
			        (title, ranking)
			 VALUES (?,?)`,
			tr.Title,
			tr.Ranking,
		)
		if err != nil {
			log.Fatalln(err)
		}

		trackID, err := res.LastInsertId()
		if err != nil {
			log.Fatalln(err)
		}

		for _, artist := range tr.Artists {
			// Does artist already exist?
			rows, err = db.Query(
				"SELECT id FROM artists WHERE name = ?",
				artist,
			)
			if err != nil {
				log.Fatalln(err)
			}

			var artistID int64

			for rows.Next() {
				if err = rows.Scan(&artistID); err != nil {
					log.Fatalln(err)
				}
			}

			if artistID == 0 {
				log.Println("Inserting artist " + artist)

				res, err = db.Exec(
					`INSERT INTO artists
					        (name)
					 VALUES (?)`,
					artist,
				)
				if err != nil {
					log.Fatalln(err)
				}

				artistID, err = res.LastInsertId()
				if err != nil {
					log.Fatalln(err)
				}

				log.Println("Artist ID: " + strconv.Itoa(int(artistID)))
			}

			_, err = db.Exec(
				`INSERT INTO track_artist
				            (track_id, artist_id)
				     VALUES (?,?)`,
				trackID, artistID,
			)
			if err != nil {
				log.Fatalln(err)
			}
		}

		// Does album already exist?
		// TODO Albums with the same name may exist
		rows, err = db.Query(
			"SELECT id FROM albums WHERE title = ?",
			tr.Album,
		)
		if err != nil {
			log.Fatalln(err)
		}

		var albumID int64

		for rows.Next() {
			if err = rows.Scan(&albumID); err != nil {
				log.Fatalln(err)
			}
		}

		if albumID == 0 {
			res, err = db.Exec(
				`INSERT INTO albums
				            (title)
				     VALUES (?)`,
				tr.Album,
			)
			if err != nil {
				log.Fatalln(err)
			}

			albumID, err = res.LastInsertId()
			if err != nil {
				log.Fatalln(err)
			}
		}

		_, err = db.Exec(
			`INSERT INTO track_album
			            (track_id, album_id)
			     VALUES (?,?)`,
			trackID, albumID,
		)
		if err != nil {
			log.Fatalln(err)
		}
	}
}
