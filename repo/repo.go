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

	albums        []albumResult
	primaryArtist artistResult
	otherArtists  []artistResult
}

func SaveTracks(db *sql.DB, inputTracks []track.Track) {
	for _, inputTrack := range inputTracks {
		trackID,
			pArtistID,
			existingAlbumsForTrack,
			existingOtherArtistsForTrack :=
			getExistingDataForTrack(db, inputTrack)

		if trackID > 0 {
			log.Printf(
				"Track '%s' exists for primary artist '%s'",
				inputTrack.Title,
				inputTrack.PrimaryArtist,
			)

			// Track (title + primary artist combo) already exists, but album
			// from inputTrack may not already be associated with track.
			// Same with secondary artists.
			if _, exists := existingAlbumsForTrack[inputTrack.Album]; !exists {
				// Check if album already exists for primary artist
				var albumID int64
				row := db.QueryRow(
					`SELECT al.id
					   FROM albums al
					   JOIN album_artist aa ON aa.album_id = al.id
					  WHERE al.title = ?
					    AND aa.artist_id  = ?`,
					inputTrack.Album,
					pArtistID,
				)

				err := row.Scan(&albumID)
				if err != nil && err != sql.ErrNoRows {
					log.Fatal(err)
				}
				if err == sql.ErrNoRows {
					log.Printf(
						"    Inserting album: %s",
						inputTrack.Album,
					)

					// Insert album
					res, err := db.Exec(
						`INSERT INTO albums
						             (title)
						      VALUES (?)`,
						inputTrack.Album,
					)
					if err != nil {
						log.Fatal(err)
					}

					albumID, err = res.LastInsertId()
					if err != nil {
						log.Fatal(err)
					}

					log.Printf("    Album ID: %d", albumID)

					// Associate album with artist
					_, err = db.Exec(
						`INSERT INTO album_artist
						             (album_id, artist_id)
						      VALUES (?,?)`,
						albumID,
						pArtistID,
					)
					if err != nil {
						log.Fatal(err)
					}
				}

				log.Printf(
					"    Associating album '%s' with track",
					inputTrack.Album,
				)

				// Associate track with album
				_, err = db.Exec(
					`INSERT INTO track_album
					             (track_id, album_id)
					      VALUES (?,?)`,
					trackID,
					albumID,
				)
				if err != nil {
					log.Fatal(err)
				}
			}
			for _, inputOtherArtist := range inputTrack.OtherArtists {
				if _, exists := existingOtherArtistsForTrack[inputOtherArtist]; !exists {
					// Check if artist exists in DB
					var otherArtistID int64
					row := db.QueryRow(
						`SELECT ar.id
						   FROM artists ar
						  WHERE ar.name = ?`,
						inputOtherArtist,
					)

					err := row.Scan(&otherArtistID)
					if err != nil && err != sql.ErrNoRows {
						log.Fatal(err)
					}
					if err == sql.ErrNoRows {
						log.Printf(
							"    Inserting artist %s",
							inputOtherArtist,
						)

						// Add artist
						res, err := db.Exec(
							`INSERT INTO artists
							             (name)
							      VALUES (?)`,
							inputOtherArtist,
						)
						if err != nil {
							log.Fatal(err)
						}

						otherArtistID, err = res.LastInsertId()
						if err != nil {
							log.Fatal(err)
						}

						log.Printf("    Artist ID: %d", otherArtistID)
					}

					log.Printf(
						"    Associating artist '%s' with track",
						inputOtherArtist,
					)

					// Associate track with artist
					_, err = db.Exec(
						`INSERT INTO track_artist
						             (track_id, artist_id, is_primary_artist)
						      VALUES ( ?, ?, 0 )`,
						trackID,
						otherArtistID,
					)
					if err != nil {
						log.Fatal(err)
					}
				}
			}
		} else {
			// Brand new track

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
				[]string{inputTrack.PrimaryArtist},
				inputTrack.OtherArtists...,
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

			// Insert album if not a duplicate. It is possible for
			// different albums to have the same name, but we assume
			// album names to be unique per primary artist.
			row := db.QueryRow(
				`SELECT al.id
				   FROM albums al
				   JOIN album_artist aa ON aa.album_id = al.id
				  WHERE al.title = ?
				    AND aa.artist_id  = ?`,
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

				// Populate album_artist
				_, err = db.Exec(
					`INSERT INTO album_artist
					        (album_id, artist_id)
					 VALUES (?,?)`,
					albumID,
					artistIDs[0], // Primary artist
				)
				if err != nil {
					log.Fatalln(err)
				}
			}

			log.Println("    Album ID: " + strconv.Itoa(int(albumID)))

			// Populate track_album
			_, err = db.Exec(
				`INSERT INTO track_album
				        (track_id, album_id)
				 VALUES (?,?)`,
				trackID,
				albumID,
			)
			if err != nil {
				log.Fatalln(err)
			}
		}
		log.Println()
	}
}

func getExistingDataForTrack(db *sql.DB, inputTrack track.Track) (
	trackID int64,
	primaryArtistID int64,
	albums map[string]bool,
	artists map[string]bool,
) {
	// Get albums for given track title and primary artist
	rows, err := db.Query(
		`SELECT t.id,
		        ar.id,
		        al.title
		   FROM albums al
		   JOIN track_album  tal ON tal.album_id = al.id
		   JOIN tracks       t   ON t.id = tal.track_id
		   JOIN track_artist tar ON tar.track_id = t.id
		   JOIN artists      ar  ON ar.id = tar.artist_id
		  WHERE t.title = ?
		    AND ar.name = ?
		    AND tar.is_primary_artist = 1`,
		inputTrack.Title,
		inputTrack.PrimaryArtist,
	)
	if err != nil {
		log.Fatal(err)
	}

	albums = map[string]bool{}
	for rows.Next() {
		var album string
		if err = rows.Scan(
			&trackID,
			&primaryArtistID,
			&album,
		); err != nil {
			log.Fatal(err)
		}
		albums[album] = true
	}

	// Get other artists for track
	rows, err = db.Query(
		`SELECT ar.name
		   FROM artists ar
		   JOIN track_artist tar ON tar.artist_id = ar.id
		  WHERE tar.track_id = ?
		    AND tar.is_primary_artist = 0`,
		trackID,
	)
	if err != nil {
		log.Fatal(err)
	}

	otherArtists := map[string]bool{}
	for rows.Next() {
		var artist string
		if err = rows.Scan(
			&artist,
		); err != nil {
			log.Fatal(err)
		}
		otherArtists[artist] = true
	}

	return trackID, primaryArtistID, albums, otherArtists
}
