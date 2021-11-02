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

	album         albumResult
	primaryArtist artistResult
	otherArtists  []artistResult
}

func SaveTracks(db *sql.DB, inputTracks []track.Track) {
	for _, inputTrack := range inputTracks {
		trackIsDupe := checkIfDuplicateTrack(db, inputTrack)

		if !trackIsDupe {
			// TODO Transaction
			// TODO Tests
			// TODO Handle empty track names, album names etc.
			// TODO Prevent duplicate tracks (across artist & album)

			log.Println("Inserting track: " + inputTrack.Title)

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

			var artistIDs []int64

			// Insert artist if not a duplicate
			for idx, artist := range append(
				inputTrack.OtherArtists,
				inputTrack.PrimaryArtist,
			) {
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
					log.Println("    Inserting artist: " + artist)

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

				}

				log.Println("    Artist ID: " + strconv.Itoa(int(artistID)))

				_, err = db.Exec(
					`INSERT INTO track_artist
				            (track_id, artist_id, is_primary_artist)
				     VALUES (?,?,?)`,
					trackID,
					artistID,
					idx == 0, // Assume primary artist is first one in list
				)
				if err != nil {
					log.Fatalln(err)
				}

				artistIDs = append(artistIDs, artistID)
			}

			// Insert album if not a duplicate. It is possible for different
			// albums to have the same name, but we assume album names to be
			// unique per primary artist.
			row := db.QueryRow(
				`SELECT al.id
				   FROM albums al
				   JOIN track_album tal ON tal.album_id = al.id
				   JOIN tracks t ON t.id = tal.track_id
				   JOIN track_artist tar ON tar.track_id = t.id
				  WHERE al.title = ?
				    AND tar.artist_id == ?`,
				inputTrack.Album,
				artistIDs[0], // Primary artist
			)

			var albumID int64

			if err = row.Scan(&albumID); err != nil && err != sql.ErrNoRows {
				log.Fatalln(err)
			}
			if err == sql.ErrNoRows {
				log.Println("    Inserting album: " + inputTrack.Album)

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

			log.Println("    Album ID: " + strconv.Itoa(int(albumID)))

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
	row := db.QueryRow(
		`SELECT 1
		   FROM tracks t
		   JOIN track_artist tar ON tar.track_id = t.id
		   JOIN artists ar ON ar.id = tar.artist_id
		  WHERE t.title = ?
		  	AND ar.name = ?
		    AND tar.is_primary_artist = 1`,
		inputTrack.Title,
		inputTrack.PrimaryArtist,
	)

	var isDuplicate bool
	if err := row.Scan(&isDuplicate); err != nil && err != sql.ErrNoRows {
		log.Fatal(err)
	} else if err == sql.ErrNoRows {
		return false
	}

	return true
}
