package repo

import (
	"database/sql"
	"log"
	"reflect"
	"testing"

	"github.com/nephila-nacrea/rank-my-music/test_utils"
	"github.com/nephila-nacrea/rank-my-music/track"
)

func init() {
	log.SetFlags(log.Llongfile)
}

func TestCheckIfDuplicateTrack(t *testing.T) {
	db := test_utils.DBSetup()

	res, err := db.Exec(
		`INSERT INTO tracks (title)
		      VALUES ("Title 1")`,
	)
	if err != nil {
		t.Fatal(err)
	}

	trackID, err := res.LastInsertId()
	if err != nil {
		t.Fatal(err)
	}

	res, err = db.Exec(
		`INSERT INTO albums (title)
		      VALUES ("Album 1")`,
	)
	if err != nil {
		t.Fatal(err)
	}

	albumID, err := res.LastInsertId()
	if err != nil {
		t.Fatal(err)
	}

	res, err = db.Exec(
		`INSERT INTO artists (name)
		      VALUES ("Artist 1")`,
	)
	if err != nil {
		t.Fatal(err)
	}

	artistID1, err := res.LastInsertId()
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec(
		`INSERT INTO track_album (track_id, album_id)
		      VALUES (?, ?)`,
		trackID,
		albumID,
	)
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec(
		`INSERT INTO track_artist (track_id, artist_id)
		      VALUES (?, ?)`,
		trackID,
		artistID1,
	)
	if err != nil {
		t.Fatal(err)
	}

	inputTrack := track.Track{
		Album:   "Album 1",
		Artists: []string{"Artist 1"},
		Title:   "Title 1",
	}

	expected := true
	got := checkIfDuplicateTrack(db, inputTrack)

	if expected != got {
		t.Errorf("Expected %t, got %t", expected, got)
	}
}

func TestSaveTracks(t *testing.T) {
	input := []track.Track{
		track.New(
			"Title 1",
			"Album 1",
			[]string{"Artist 1", "Artist 2", "Artist 3"},
		),
		// Complete duplicate
		track.New(
			"Title 1",
			"Album 1",
			[]string{"Artist 1", "Artist 2", "Artist 3"},
		),
		// // New title, duplicate album & artists
		// track.New(
		// 	"Title 2",
		// 	"Album 1",
		// 	[]string{"Artist 1", "Artist 2", "Artist 3"},
		// ),
		// // Duplicate title, new album, new artist
		// track.New(
		// 	"Title 1",
		// 	"Album 2",
		// 	[]string{"Artist 4"},
		// ),
		// // Duplicate title, duplicate album, new artist
		// track.New(
		// 	"Title 1",
		// 	"Album 2",
		// 	[]string{"Artist 5"},
		// ),
		// // Duplicate title, duplicate album, new *primary* artist
		// track.New(
		// 	"Title 1",
		// 	"Album 2",
		// 	[]string{"Artist 2", "Artist 1", "Artist 3"},
		// ),

		// TODO
		// Empty strings
		// Long names

	}

	db := test_utils.DBSetup()

	SaveTracks(db, input)

	expected := []trackResult{
		{
			id:    1,
			title: "Title 1",
			album: albumResult{
				id:    1,
				title: "Album 1",
			},
			artists: []artistResult{
				{
					id:   1,
					name: "Artist 1",
				},
				{
					id:   2,
					name: "Artist 2",
				},
				{
					id:   3,
					name: "Artist 3",
				},
			},
		},
	}

	got := readDB(t, db)

	if !reflect.DeepEqual(expected, got) {
		t.Errorf("Expected %v, got %v", expected, got)
	}
}

func readDB(t *testing.T, db *sql.DB) []trackResult {
	rows, err := db.Query(
		`SELECT t.id,
		        t.title,
		        al.id,
		        al.title
		  FROM tracks t
		  JOIN track_album tal ON tal.track_id = t.id
		  JOIN albums al ON al.id = tal.album_id`,
	)
	if err != nil {
		t.Fatal(err)
	}

	tracks := []trackResult{}

	for rows.Next() {
		var track trackResult
		var album albumResult

		if err = rows.Scan(
			&track.id,
			&track.title,
			&album.id,
			&album.title,
		); err != nil {
			t.Fatal(err)
		}

		track.album = album
		tracks = append(tracks, track)
	}

	for i := range tracks {
		rowsArtist, err := db.Query(
			`SELECT ar.id,
			        ar.name
			   FROM artists ar
			   JOIN track_artist tar ON tar.artist_id = ar.id
			  WHERE tar.track_id = ?`,
			tracks[i].id,
		)
		if err != nil {
			t.Fatal(err)
		}

		artists := []artistResult{}

		for rowsArtist.Next() {
			var artist artistResult

			if rowsArtist.Scan(
				&artist.id,
				&artist.name,
			); err != nil {
				t.Fatal(err)
			}

			artists = append(artists, artist)
		}

		tracks[i].artists = artists
	}

	return tracks
}
