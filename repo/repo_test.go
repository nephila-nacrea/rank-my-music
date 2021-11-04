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
			[]string{"Artist 1", "Artist 2", "Artist 3"},
		),
		// Complete duplicate
		track.New(
			"Title 1",
			"Album 1",
			[]string{"Artist 1", "Artist 2", "Artist 3"},
		),
		// Complete duplicate with non-primary artists switched
		track.New(
			"Title 1",
			"Album 1",
			[]string{"Artist 1", "Artist 3", "Artist 2"},
		),
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
			otherArtists: []artistResult{
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
		t.Errorf("\nExpected:\n%#v\ngot:\n%#v", expected, got)
	}

	//////////////////////////////////////////////////////////////////////////

	// New track title, duplicate album title & artists (multiple)
	input = []track.Track{
		track.New(
			"Title 2",
			"Album 1",
			[]string{"Artist 1", "Artist 2", "Artist 3"},
		),
	}

	SaveTracks(db, input)

	expected = append(
		expected,
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
			otherArtists: []artistResult{
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
	)

	got = readDB(t, db)
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("\nExpected:\n%#v\ngot:\n%#v", expected, got)
	}

	//////////////////////////////////////////////////////////////////////////

	// New track title, duplicate album title & artist (single)
	input = []track.Track{
		track.New(
			"Title 3",
			"Album 1",
			[]string{"Artist 1"},
		),
	}

	SaveTracks(db, input)

	expected = append(
		expected,
		trackResult{
			id:    3,
			title: "Title 3",
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

	//////////////////////////////////////////////////////////////////////////

	// Duplicate track title, new album title, new artist
	input = []track.Track{
		track.New(
			"Title 1",
			"Album 2",
			[]string{"Artist 4"},
		),
	}

	SaveTracks(db, input)

	expected = append(
		expected,
		trackResult{
			id:    4,
			title: "Title 1",
			album: albumResult{
				id:    2,
				title: "Album 2",
			},
			primaryArtist: artistResult{
				id:   4,
				name: "Artist 4",
			},
			otherArtists: []artistResult{},
		},
	)

	got = readDB(t, db)
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("\nExpected:\n%#v\ngot:\n%#v", expected, got)
	}

	//////////////////////////////////////////////////////////////////////////

	// Duplicate track title, duplicate album title, new primary artist
	// => different album with same name as another
	input = []track.Track{
		track.New(
			"Title 1",
			"Album 2",
			[]string{"Artist 5"},
		),
	}

	SaveTracks(db, input)

	expected = append(
		expected,
		trackResult{
			id:    5,
			title: "Title 1",
			album: albumResult{
				id:    3,
				title: "Album 2",
			},
			primaryArtist: artistResult{
				id:   5,
				name: "Artist 5",
			},
			otherArtists: []artistResult{},
		},
	)

	got = readDB(t, db)
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("\nExpected:\n%#v\ngot:\n%#v", expected, got)
	}

	//////////////////////////////////////////////////////////////////////////

	// Duplicate track title, duplicate album title, duplicate artists but
	// with primary artist new for given track
	// => different album with same name as another
	input = []track.Track{
		track.New(
			"Title 1",
			"Album 1",
			[]string{"Artist 2", "Artist 1", "Artist 3"},
		),
	}

	SaveTracks(db, input)

	expected = append(
		expected,
		trackResult{
			id:    6,
			title: "Title 1",
			album: albumResult{
				id:    4,
				title: "Album 1",
			},
			primaryArtist: artistResult{
				id:   2,
				name: "Artist 2",
			},
			otherArtists: []artistResult{},
		},
	)

	got = readDB(t, db)
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("\nExpected:\n%#v\ngot:\n%#v", expected, got)
	}

	// input := []track.Track{

	// 	// Duplicate track title, duplicate album title, new *primary* artist
	// 	// =>
	// 	track.New(
	// 		"Title 1",
	// 		"Album 1",
	// 		[]string{"Artist 2", "Artist 1", "Artist 3"},
	// 	),

	// 	// TODO
	// 	// Empty strings
	// 	// Long names

	// }

	// expected := []trackResult{
	// 	{
	// 		id:    1,
	// 		title: "Title 1",
	// 		album: albumResult{
	// 			id:    1,
	// 			title: "Album 1",
	// 		},
	// 		artists: []artistResult{
	// 			{
	// 				id:   1,
	// 				name: "Artist 1",
	// 			},
	// 			{
	// 				id:   2,
	// 				name: "Artist 2",
	// 			},
	// 			{
	// 				id:   3,
	// 				name: "Artist 3",
	// 			},
	// 		},
	// 	},
	// 	{
	// 		id:    2,
	// 		title: "Title 2",
	// 		album: albumResult{
	// 			id:    1,
	// 			title: "Album 1",
	// 		},
	// 		artists: []artistResult{
	// 			{
	// 				id:   1,
	// 				name: "Artist 1",
	// 			},
	// 			{
	// 				id:   2,
	// 				name: "Artist 2",
	// 			},
	// 			{
	// 				id:   3,
	// 				name: "Artist 3",
	// 			},
	// 		},
	// 	},
	// 	{
	// 		id:    3,
	// 		title: "Title 3",
	// 		album: albumResult{
	// 			id:    1,
	// 			title: "Album 1",
	// 		},
	// 		artists: []artistResult{
	// 			{
	// 				id:   1,
	// 				name: "Artist 1",
	// 			},
	// 		},
	// 	},
	// 	{
	// 		id:    4,
	// 		title: "Title 1",
	// 		album: albumResult{
	// 			id:    2,
	// 			title: "Album 2",
	// 		},
	// 		artists: []artistResult{
	// 			{
	// 				id:   4,
	// 				name: "Artist 4",
	// 			},
	// 		},
	// 	},
	// 	{
	// 		id:    5,
	// 		title: "Title 1",
	// 		album: albumResult{
	// 			id:    3,
	// 			title: "Album 2",
	// 		},
	// 		artists: []artistResult{
	// 			{
	// 				id:   5,
	// 				name: "Artist 5",
	// 			},
	// 		},
	// 	},
	// }

	// got := readDB(t, db)
	// if !reflect.DeepEqual(expected, got) {
	// 	t.Errorf("\nExpected:\n%#v\ngot:\n%#v", expected, got)
	// }
}

func readDB(t *testing.T, db *sql.DB) []trackResult {
	rows, err := db.Query(
		`SELECT t.id,
		        t.title,
		        al.id,
		        al.title
		  FROM tracks t
		  JOIN track_album tal ON tal.track_id = t.id
		  JOIN albums al ON al.id = tal.album_id
	      ORDER BY t.id`,
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
