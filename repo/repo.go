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
	id            int
	title         string
	musicBrainzID string

	albums        []albumResult
	primaryArtist artistResult
	otherArtists  []artistResult
}

func SaveTracks(db *sql.DB, inputTracks []track.Track) {
	for _, inputTrack := range inputTracks {
		err := saveTrack(db, inputTrack)
		if err != nil {
			log.Printf("%s\n\n", err)
		}
	}
}

// Save individual track, wrapped in a transaction
func saveTrack(db *sql.DB, inputTrack track.Track) error {
	existingTrack := getExistingDataForTrackMBID(db, inputTrack.MusicBrainzID)

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if existingTrack.InternalID > 0 {
		log.Printf("Track '%s' exists for MBID '%s'",
			existingTrack.Title,
			existingTrack.MusicBrainzID,
		)

		// TODO:
		// It is possible for a different primary artist to
		// be provided. In which case, we need to uncouple the old primary
		// artist from the track and attach the new.

		// Album from inputTrack may not already be associated with track in
		// DB
		var existingAlbumsForTrack map[string]bool
		for _, album := range existingTrack.Albums {
			existingAlbumsForTrack[album.Title] = true
		}

		// We assume the input track only ever has one album
		// TODO What if album list empty?
		inputAlbum := inputTrack.Albums[0]
		if _, exists := existingAlbumsForTrack[inputAlbum.Title]; !exists {
			// Does album exist already?
			var albumInternalID int64

			row := tx.QueryRow(
				`SELECT id
				   FROM albums
				  WHERE musicbrainz_id = ?`,
				inputAlbum.MusicBrainzID,
			)

			err := row.Scan(&albumInternalID)
			if err != nil && err != sql.ErrNoRows {
				return err
			}
			if err == sql.ErrNoRows {
				log.Printf(
					"    Inserting album: %s",
					inputAlbum.Title,
				)

				// Insert album
				res, err := tx.Exec(
					`INSERT INTO albums
					             (musicbrainz_id, title)
					      VALUES (?,?)`,
					inputAlbum.MusicBrainzID,
					inputAlbum.Title,
				)
				if err != nil {
					return err
				}

				albumInternalID, err = res.LastInsertId()
				if err != nil {
					return err
				}
			}

			log.Printf(
				"    Associating album '%s' with track",
				inputAlbum.Title,
			)

			// Associate track with album
			_, err = tx.Exec(
				`INSERT INTO track_album
				             (track_id, album_id)
				      VALUES (?,?)`,
				existingTrack.InternalID,
				albumInternalID,
			)
			if err != nil {
				return err
			}
		}

		// Secondary artists from input track may not already be associated
		// with track in DB
		for _, inputOtherArtist := range inputTrack.OtherArtists {
			var existingOtherArtistsForTrack map[string]bool
			if _, exists := existingOtherArtistsForTrack[inputOtherArtist.Name]; !exists {
				// Check if artist exists in DB.
				// Secondary artists do not have a MusicBrainz ID so we have
				// to go by name.
				// TODO Is there a better way of handling secondary artist
				// data?

				var otherArtistInternalID int64
				row := tx.QueryRow(
					`SELECT ar.id
					   FROM artists ar
					  WHERE ar.name = ?`,
					inputOtherArtist.Name,
				)

				err := row.Scan(&otherArtistInternalID)
				if err != nil && err != sql.ErrNoRows {
					return err
				}
				if err == sql.ErrNoRows {
					log.Printf(
						"    Inserting secondary artist %s",
						inputOtherArtist.Name,
					)

					// Add artist
					res, err := tx.Exec(
						`INSERT INTO artists
						             (name)
						      VALUES (?)`,
						inputOtherArtist,
					)
					if err != nil {
						return err
					}

					otherArtistInternalID, err = res.LastInsertId()
					if err != nil {
						return err
					}

					log.Printf("    Artist ID: %d", otherArtistInternalID)
				}

				log.Printf(
					"    Associating artist '%s' with track",
					inputOtherArtist.Name,
				)

				// Associate track with artist
				_, err = tx.Exec(
					`INSERT INTO track_artist
					             (track_id, artist_id, is_primary_artist)
					      VALUES (?, ?, 0)`,
					existingTrack.InternalID,
					otherArtistInternalID,
				)
				if err != nil {
					return err
				}
			}
		}
	} else {
		// Brand new track

		// TODO REFACTOR FOR MBID

		// TODO Handle empty track names, album names etc.

		log.Println("Inserting track: " + inputTrack.Title)
		log.Println("MusicBrainz ID: " + inputTrack.MusicBrainzID)

		res, err := tx.Exec(
			`INSERT INTO tracks
			             (title, musicbrainz_id, ranking)
			      VALUES (?,?,?)`,
			inputTrack.Title,
			inputTrack.MusicBrainzID,
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
			row := tx.QueryRow(
				"SELECT id FROM artists WHERE name = ?",
				artist,
			)

			var artistID int64
			if err = row.Scan(&artistID); err != nil && err != sql.ErrNoRows {
				log.Fatalln(err)
			}
			if err == sql.ErrNoRows {
				log.Println("    Inserting artist: " + artist)

				res, err = tx.Exec(
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

			_, err = tx.Exec(
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
		row := tx.QueryRow(
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

			res, err = tx.Exec(
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
			_, err = tx.Exec(
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
		_, err = tx.Exec(
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

	log.Print("    Committing transaction\n\n")
	return tx.Commit()
}

func getExistingDataForTrackMBID(db *sql.DB, trackMBID string) (
	existingTrack track.Track,
) {
	row := db.QueryRow(
		`SELECT t.id,
		        t.musicbrainz_id,
		        t.title,
		        ar.id,
		        ar.musicbrainz_id,
		        ar.name
		   FROM tracks       t
		   JOIN track_artist tar ON tar.track_id = t.id
		   JOIN artists      ar  ON ar.id = tar.artist_id
		  WHERE t.musicbrainz_id      = ?
		    AND tar.is_primary_artist = 1`,
		trackMBID,
	)

	err := row.Scan(
		existingTrack.InternalID,
		existingTrack.MusicBrainzID,
		existingTrack.Title,
		existingTrack.PrimaryArtist.InternalID,
		existingTrack.PrimaryArtist.MusicBrainzID,
		existingTrack.PrimaryArtist.Name,
	)
	if err != nil && err != sql.ErrNoRows {
		log.Fatalln(err)
	} else if err == sql.ErrNoRows {
		return track.Track{}
	}

	// Get secondary artists
	rows, err := db.Query(
		`SELECT ar.id,
		        ar.musicbrainz_id,
		        ar.name
		   FROM tracks       t
		   JOIN track_artist tar ON tar.track_id = t.id
		   JOIN artists      ar  ON ar.id = tar.artist_id
		  WHERE t.musicbrainz_id      = ?
		    AND tar.is_primary_artist = 0`,
		trackMBID,
	)
	if err != nil {
		log.Fatalln(err)
	}

	var artists []track.Artist
	for rows.Next() {
		var artist track.Artist

		if err = rows.Scan(
			artist.InternalID,
			artist.MusicBrainzID,
			artist.Name,
		); err != nil {
			log.Fatalln(err)
		}

		artists = append(artists, artist)
	}

	existingTrack.OtherArtists = artists

	// Get albums
	rows, err = db.Query(
		`SELECT al.id,
		        al.musicbrainz_id,
		        al.title
		   FROM tracks      t
		   JOIN track_album tal ON tal.track_id = t.id
		   JOIN albums      al  ON al.id        = tal.album_id
		  WHERE t.musicbrainz_id = ?`,
		trackMBID,
	)
	if err != nil {
		log.Fatalln(err)
	}

	var albums []track.Album
	for rows.Next() {
		var album track.Album

		if err = rows.Scan(
			&album.InternalID,
			&album.MusicBrainzID,
			&album.Title,
		); err != nil {
			log.Fatal(err)
		}

		albums = append(albums, album)
	}

	existingTrack.Albums = albums

	return existingTrack
}
