package repo

import (
	"database/sql"
	"log"
	"strconv"

	"github.com/nephila-nacrea/rank-my-music/track"
)

type albumResult struct {
	id    int
	title string
}

type artistResult struct {
	id   int
	name string
}

type trackResult struct {
	id    int
	title string

	album   albumResult
	artists []artistResult
}

func SaveTracks(db *sql.DB, inputTracks []track.Track) {
	for _, inputTrack := range inputTracks {
		trackIsDupe := checkIfDuplicateTrack(db, inputTrack)

		if !trackIsDupe {
			// TODO Transaction
			// TODO Tests
			// TODO Handle empty track names, album names etc.
			// TODO Prevent duplicate tracks (across artist & album)

			res, err := db.Exec(
				`INSERT INTO tracks
			        (title, ranking)
			 VALUES (?,?)`,
				inputTrack.Title,
				inputTrack.Ranking,
			)
			if err != nil {
				log.Fatalln(err)
			}

			trackID, err := res.LastInsertId()
			if err != nil {
				log.Fatalln(err)
			}

			// Insert artist if not a duplicate
			for _, artist := range inputTrack.Artists {
				// Does artist already exist?
				row := db.QueryRow(
					"SELECT id FROM artists WHERE name = ?",
					artist,
				)

				var artistID int64

				if err = row.Scan(&artistID); err != nil && err != sql.ErrNoRows {
					log.Fatalln(err)
				}
				if err == sql.ErrNoRows {
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
				inputTrack.Album,
			)

			var albumID int64

			if err = row.Scan(&albumID); err != nil && err != sql.ErrNoRows {
				log.Fatalln(err)
			}
			if err == sql.ErrNoRows {
				res, err = db.Exec(
					`INSERT INTO albums
				            (title)
				     VALUES (?)`,
					inputTrack.Album,
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

func checkIfDuplicateTrack(db *sql.DB, inputTrack track.Track) bool {
	var trackIsDupe bool

	// Get any existing data for given track title. We will do duplication
	// checks against these.
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
		inputTrack.Title,
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
		// Get artists
		rows, err = db.Query(
			`SELECT ar.name
			   FROM artists ar
			   JOIN track_artist tar ON tar.artist_id = ar.id
			  WHERE tar.track_id = ?`,
			et.id,
		)
		if err != nil {
			log.Fatal(err)
		}

		var existingArtists = map[string]bool{}

		for rows.Next() {
			var name string
			if err = rows.Scan(
				&name,
			); err != nil {
				log.Fatal(err)
			}

			existingArtists[name] = true
		}

		var hasSharedArtist bool
		for _, inputArtist := range inputTrack.Artists {
			if _, exists := existingArtists[inputArtist]; exists {
				hasSharedArtist = true
			}
		}

		if inputTrack.Album == et.album.title &&
			hasSharedArtist {
			trackIsDupe = true
		}
	}

	return trackIsDupe
}
