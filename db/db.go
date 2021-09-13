package db

import (
	"database/sql"
	"log"
	"strconv"

	"example.com/rank-my-music/track"
)

const DefaultDSN = "file:ranked_music.sqlt"

type DB struct {
	DSN    string
	Handle *sql.DB
}

func New(dsn string) DB {
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Database connected")

	return DB{dsn, db}
}

func (db *DB) PopulateDB(tracks []track.Track) {
	for _, tr := range tracks {
		// TODO Transaction
		// TODO Tests
		// TODO Handle empty track names, album names etc.
		// TODO Prevent duplicate tracks (across artist & album)

		dbh := db.Handle

		res, err := dbh.Exec(
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
			rows, err := dbh.Query(
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

				res, err = dbh.Exec(
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

			_, err = dbh.Exec(
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
		rows, err := dbh.Query(
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
			res, err = dbh.Exec(
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

		_, err = dbh.Exec(
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