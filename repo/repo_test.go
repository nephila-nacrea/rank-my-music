package repo

import (
	"database/sql"
	"reflect"
	"testing"

	"github.com/nephila-nacrea/rank-my-music/test_utils"
	"github.com/nephila-nacrea/rank-my-music/track"
)

type albumRes struct {
	id    int
	title string
}

type artistRes struct {
	id   int
	name string
}

type trackRes struct {
	id    int
	title string

	album   albumRes
	artists []artistRes
}

func TestSaveTracks(t *testing.T) {
	input := []track.Track{
		track.New(
			"Title 1",
			"Album 1",
			[]string{"Artist 1", "Artist 2", "Artist 3"},
		),
		// // Complete duplicate
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

	expected := []trackRes{
		{
			id:    1,
			title: "Title 1",
			album: albumRes{
				id:    1,
				title: "Album 1",
			},
			artists: []artistRes{
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

	got := ReadDB(t, db)

	if !reflect.DeepEqual(expected, got) {
		t.Errorf("Expected %v, got %v", expected, got)
	}
}

func ReadDB(t *testing.T, db *sql.DB) []trackRes {
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

	tracks := []trackRes{}

	for rows.Next() {
		var track trackRes
		var album albumRes

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

		artists := []artistRes{}

		for rowsArtist.Next() {
			var artist artistRes

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
