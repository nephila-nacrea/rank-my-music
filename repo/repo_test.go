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

func TestSaveTracks(t *testing.T) {
	db := test_utils.DBSetup()

	input := []track.Track{
		track.New(
			"Title 1",
			"Album 1",
			"Artist 1",
			[]string{},
		),
		// Complete duplicate
		track.New(
			"Title 1",
			"Album 1",
			"Artist 1",
			[]string{},
		),
		// Duplicate with secondary artists included
		// TODO
	}

	SaveTracks(db, input)

	expected := []trackResult{
		{
			id:    1,
			title: "Title 1",
			album: albumResult{
				id:    1,
				title: "Album 1",
			},
			primaryArtist: artistResult{
				id:   1,
				name: "Artist 1",
			},
			otherArtists: []artistResult{},
		},
	}

	got := readDB(t, db)
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("\nExpected:\n%#v\ngot:\n%#v", expected, got)
	}

	//////////////////////////////////////////////////////////////////////////

	// New track title
	input = []track.Track{
		track.New(
			"Title 2",
			"Album 1",
			"Artist 1",
			[]string{},
		),
	}

	SaveTracks(db, input)

	expected = append(expected,
		trackResult{
			id:    2,
			title: "Title 2",
			album: albumResult{
				id:    1,
				title: "Album 1",
			},
			primaryArtist: artistResult{
				id:   1,
				name: "Artist 1",
			},
			otherArtists: []artistResult{},
		},
	)

	got = readDB(t, db)
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("\nExpected:\n%#v\ngot:\n%#v", expected, got)
	}
}

func readDB(t *testing.T, db *sql.DB) []trackResult {
	rows, err := db.Query(
		`SELECT t.id,
		        t.title,
		        al.id,
		        al.title,
		        ar.id,
		        ar.name
		  FROM tracks       t
		  JOIN track_album  tal ON tal.track_id = t.id
		  JOIN albums       al  ON al.id = tal.album_id
		  JOIN track_artist tar ON tar.track_id = t.id
		  JOIN artists      ar  ON ar.id = tar.artist_id
		 WHERE tar.is_primary_artist = 1
		 ORDER BY t.id`,
	)
	if err != nil {
		t.Fatal(err)
	}

	tracks := []trackResult{}

	for rows.Next() {
		var track trackResult
		var album albumResult
		var artist artistResult

		if err = rows.Scan(
			&track.id,
			&track.title,
			&album.id,
			&album.title,
			&artist.id,
			&artist.name,
		); err != nil {
			t.Fatal(err)
		}

		track.album = album
		track.primaryArtist = artist
		tracks = append(tracks, track)
	}

	// Get secondary artists
	for i := range tracks {
		rowsArtist, err := db.Query(
			`SELECT ar.id,
			        ar.name
			   FROM artists ar
			   JOIN track_artist tar ON tar.artist_id = ar.id
			  WHERE tar.track_id = ?
			    AND tar.is_primary_artist = 0`,
			tracks[i].id,
		)
		if err != nil {
			t.Fatal(err)
		}

		artists := []artistResult{}

		for rowsArtist.Next() {
			var artist artistResult

			if err = rowsArtist.Scan(
				&artist.id,
				&artist.name,
			); err != nil {
				t.Fatal(err)
			}

			artists = append(artists, artist)
		}

		tracks[i].otherArtists = artists
	}

	return tracks
}
