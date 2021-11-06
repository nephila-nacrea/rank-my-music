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

	expected := map[int]trackResult{
		1: {
			id:    1,
			title: "Title 1",
			albums: []albumResult{
				{
					id:    1,
					title: "Album 1",
				},
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

	// New track title for existing album & artist
	input = []track.Track{
		track.New(
			"Title 2",
			"Album 1",
			"Artist 1",
			[]string{},
		),
	}

	SaveTracks(db, input)

	expected[2] = trackResult{
		id:    2,
		title: "Title 2",
		albums: []albumResult{
			{
				id:    1,
				title: "Album 1",
			},
		},
		primaryArtist: artistResult{
			id:   1,
			name: "Artist 1",
		},
		otherArtists: []artistResult{},
	}

	got = readDB(t, db)
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("\nExpected:\n%#v\ngot:\n%#v", expected, got)
	}

	//////////////////////////////////////////////////////////////////////////

	// New artist with same album title & track title as existing one
	input = []track.Track{
		track.New(
			"Title 1",
			"Album 1",
			"Artist 2",
			[]string{},
		),
	}

	SaveTracks(db, input)

	// Track with duplicate name stored with new ID,
	// album with duplicate name stored with new ID
	expected[3] = trackResult{
		id:    3,
		title: "Title 1",
		albums: []albumResult{
			{
				id:    2,
				title: "Album 1",
			},
		},
		primaryArtist: artistResult{
			id:   2,
			name: "Artist 2",
		},
		otherArtists: []artistResult{},
	}

	got = readDB(t, db)
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("\nExpected:\n%#v\ngot:\n%#v", expected, got)
	}

	//////////////////////////////////////////////////////////////////////////

	// New album for existing artist, duplicate track title
	// (= artist having same track on different albums, e.g. original
	// album then greatest hits album)
	input = []track.Track{
		track.New(
			"Title 1",
			"Album 2",
			"Artist 1",
			[]string{},
		),
	}

	SaveTracks(db, input)

	// New album data stored, but track not added
	expected[1] = trackResult{
		id:    1,
		title: "Title 1",
		albums: []albumResult{
			{
				id:    1,
				title: "Album 1",
			},
			{
				id:    3,
				title: "Album 2",
			},
		},
		primaryArtist: artistResult{
			id:   1,
			name: "Artist 1",
		},
		otherArtists: []artistResult{},
	}

	got = readDB(t, db)
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("\nExpected:\n%#v\ngot:\n%#v", expected, got)
	}

	//////////////////////////////////////////////////////////////////////////

	// New track title, new artist, existing album title
	input = []track.Track{
		track.New(
			"Title 3",
			"Album 1",
			"Artist 3",
			[]string{},
		),
	}

	SaveTracks(db, input)

	// New track & album data stored
	expected[4] = trackResult{
		id:    4,
		title: "Title 3",
		albums: []albumResult{
			{
				id:    4,
				title: "Album 1",
			},
		},
		primaryArtist: artistResult{
			id:   3,
			name: "Artist 3",
		},
		otherArtists: []artistResult{},
	}

	got = readDB(t, db)
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("\nExpected:\n%#v\ngot:\n%#v", expected, got)
	}

	//////////////////////////////////////////////////////////////////////////

	// New track title, existing artist, new album title
	input = []track.Track{
		track.New(
			"Title 4",
			"Album 3",
			"Artist 1",
			[]string{},
		),
	}

	SaveTracks(db, input)

	// New track & album data stored
	expected[5] = trackResult{
		id:    5,
		title: "Title 4",
		albums: []albumResult{
			{
				id:    5,
				title: "Album 3",
			},
		},
		primaryArtist: artistResult{
			id:   1,
			name: "Artist 1",
		},
		otherArtists: []artistResult{},
	}

	got = readDB(t, db)
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("\nExpected:\n%#v\ngot:\n%#v", expected, got)
	}

	//////////////////////////////////////////////////////////////////////////

	// Existing track title, new artist, new album title
	input = []track.Track{
		track.New(
			"Title 1",
			"Album 4",
			"Artist 4",
			[]string{},
		),
	}

	SaveTracks(db, input)

	// New data stored for everything
	expected[6] = trackResult{
		id:    6,
		title: "Title 1",
		albums: []albumResult{
			{
				id:    6,
				title: "Album 4",
			},
		},
		primaryArtist: artistResult{
			id:   4,
			name: "Artist 4",
		},
		otherArtists: []artistResult{},
	}

	got = readDB(t, db)
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("\nExpected:\n%#v\ngot:\n%#v", expected, got)
	}

	//////////////////////////////////////////////////////////////////////////

	// Add track with secondary artists
	input = []track.Track{
		track.New(
			"Title 101",
			"Album 101",
			"Artist 101",
			[]string{
				"Artist 201",
				"Artist 301",
			},
		),
	}

	SaveTracks(db, input)

	// New data stored for everything
	expected[7] = trackResult{
		id:    7,
		title: "Title 101",
		albums: []albumResult{
			{
				id:    7,
				title: "Album 101",
			},
		},
		primaryArtist: artistResult{
			id:   5,
			name: "Artist 101",
		},
		otherArtists: []artistResult{
			{
				id:   6,
				name: "Artist 201",
			},
			{
				id:   7,
				name: "Artist 301",
			},
		},
	}

	got = readDB(t, db)
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("\nExpected:\n%#v\ngot:\n%#v", expected, got)
	}

	//////////////////////////////////////////////////////////////////////////

	// Add new artist as secondary artist to existing track
	input = []track.Track{
		track.New(
			"Title 1",
			"Album 1",
			"Artist 1",
			[]string{
				"Artist 401",
			},
		),
	}

	SaveTracks(db, input)

	// Track updated to have secondary artist
	tmp := expected[1]
	tmp.otherArtists = []artistResult{
		{
			id:   8,
			name: "Artist 401",
		},
	}
	expected[1] = tmp

	got = readDB(t, db)
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("\nExpected:\n%#v\ngot:\n%#v", expected, got)
	}
}

func readDB(t *testing.T, db *sql.DB) map[int]trackResult {
	rows, err := db.Query(
		`SELECT t.id,
		        t.title,
		--        al.id,
		--        al.title,
		        ar.id,
		        ar.name
		  FROM tracks       t
		--  JOIN track_album  tal ON tal.track_id = t.id
		--  JOIN albums       al  ON al.id = tal.album_id
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
		// var album albumResult
		var artist artistResult

		if err = rows.Scan(
			&track.id,
			&track.title,
			// &album.id,
			// &album.title,
			&artist.id,
			&artist.name,
		); err != nil {
			t.Fatal(err)
		}

		// track.album = album
		track.primaryArtist = artist
		tracks = append(tracks, track)
	}

	trackMap := map[int]trackResult{}

	// Get albums & secondary artists
	for i := range tracks {
		rowsAlbum, err := db.Query(
			`SELECT al.id,
			        al.title
			   FROM albums al
			   JOIN track_album tal ON tal.album_id = al.id
			  WHERE tal.track_id = ?
			  ORDER BY al.id`,
			tracks[i].id,
		)
		if err != nil {
			t.Fatal(err)
		}

		albums := []albumResult{}

		for rowsAlbum.Next() {
			var album albumResult

			if err = rowsAlbum.Scan(
				&album.id,
				&album.title,
			); err != nil {
				t.Fatal(err)
			}

			albums = append(albums, album)
		}

		rowsArtist, err := db.Query(
			`SELECT ar.id,
			        ar.name
			   FROM artists ar
			   JOIN track_artist tar ON tar.artist_id = ar.id
			  WHERE tar.track_id = ?
			    AND tar.is_primary_artist = 0
			  ORDER BY ar.id`,
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

		tracks[i].albums = albums
		tracks[i].otherArtists = artists

		trackMap[tracks[i].id] = tracks[i]
	}

	return trackMap
}
