package repo

import (
	"database/sql"
	"log"
	"strconv"

	"github.com/nephila-nacrea/rank-my-music/track"
)

func SaveTracks(db *sql.DB, tracks []track.Track) {
	for _, track := range tracks {
		var trackIsDupe bool

		// Get any existing data for given track title. We will do duplication
		// checks against these.
		type trackResult struct {
			id    int
			title string

			album struct {
				id    int
				title string
			}
			artists []struct {
				id   int
				name string
			}
		}
		existingTracks := []trackResult{}

		rows, err := db.Query(
			`SELECT t.id,
			        t.title,
			        al.id,
			        al.title
			   FROM tracks t
			   JOIN track_album tal ON tal.track_id = t.id
			   JOIN albums al ON al.id = tal.album_id
			  WHERE t.title = ?`,
			track.Title,
		)
		if err != nil {
			log.Fatal(err)
		}

		for rows.Next() {
			var tr trackResult

			if err = rows.Scan(
				&tr.id,
				&tr.title,
				&tr.album.id,
				&tr.album.title,
			); err != nil {
				log.Fatal(err)
			}

			existingTracks = append(existingTracks, tr)
		}

		for _, et := range existingTracks {
			// TODO
		}

		if !trackIsDupe {

			// TODO Transaction
			// TODO Tests
			// TODO Handle empty track names, album names etc.
			// TODO Prevent duplicate tracks (across artist & album)

			res, err := db.Exec(
				`INSERT INTO tracks
			        (title, ranking)
			 VALUES (?,?)`,
				track.Title,
				track.Ranking,
			)
			if err != nil {
				log.Fatalln(err)
			}

			trackID, err := res.LastInsertId()
			if err != nil {
				log.Fatalln(err)
			}

			// Insert artist if not a duplicate
			for _, artist := range track.Artists {
				// Does artist already exist?
				row := db.QueryRow(
					"SELECT id FROM artists WHERE name = ?",
					artist,
				)

				var artistID int64

				if err = row.Scan(&artistID); err != nil && err != sql.ErrNoRows {
					log.Fatalln(err)
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

			// Insert album if not a duplicate. It is possible for different
			// albums to have the same name, but we assume album names to be
			// unique for an artist.
			row := db.QueryRow(
				"SELECT id FROM albums WHERE title = ?",
				track.Album,
			)

			var albumID int64

			if err = row.Scan(&albumID); err != nil {
				log.Fatalln(err)
			}

			// Check against artists

			if albumID == 0 {
				res, err = db.Exec(
					`INSERT INTO albums
				            (title)
				     VALUES (?)`,
					track.Album,
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
}
