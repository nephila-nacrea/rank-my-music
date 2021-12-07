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

	t.Log("Brand new track with MusicBrainz ID plus a complete duplicate")

	input := []track.Track{
		track.New(
			"Title 1",
			"Album 1",
			"Artist 1",
			[]string{},
			"MB1",
		),
		// Complete duplicate
		track.New(
			"Title 1",
			"Album 1",
			"Artist 1",
			[]string{},
			"MB1",
		),
	}

	SaveTracks(db, input)

	expected := map[int]trackResult{
		1: {
			id:            1,
			title:         "Title 1",
			musicBrainzID: "MB1",
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

	t.Log("Track with existing MusicBrainz ID but different album")

	input = []track.Track{
		track.New(
			"Title 1",
			"Album 2",
			"Artist 1",
			[]string{},
			"MB1",
		),
	}

	SaveTracks(db, input)

	expected = map[int]trackResult{
		1: {
			id:            1,
			title:         "Title 1",
			musicBrainzID: "MB1",
			albums: []albumResult{
				{
					id:    1,
					title: "Album 1",
				},
				{
					id:    2,
					title: "Album 2",
				},
			},
			primaryArtist: artistResult{
				id:   1,
				name: "Artist 1",
			},
			otherArtists: []artistResult{},
		},
	}

	got = readDB(t, db)
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("\nExpected:\n%#v\ngot:\n%#v", expected, got)
	}

	//////////////////////////////////////////////////////////////////////////

	t.Log("Track with existing MusicBrainz ID, adding secondary artists")

	input = []track.Track{
		track.New(
			"Title 1",
			"Album 1",
			"Artist 1",
			[]string{"Artist 2", "Artist 3"},
			"MB1",
		),
	}

	SaveTracks(db, input)

	expected = map[int]trackResult{
		1: {
			id:            1,
			title:         "Title 1",
			musicBrainzID: "MB1",
			albums: []albumResult{
				{
					id:    1,
					title: "Album 1",
				},
				{
					id:    2,
					title: "Album 2",
				},
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

	got = readDB(t, db)
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("\nExpected:\n%#v\ngot:\n%#v", expected, got)
	}

	//////////////////////////////////////////////////////////////////////////

	t.Log("Track with existing MusicBrainz ID but new secondary artists")
	t.Log("New artists should be added to existing list")

	input = []track.Track{
		track.New(
			"Title 1",
			"Album 1",
			"Artist 1",
			[]string{"Artist 4"},
			"MB1",
		),
	}

	SaveTracks(db, input)

	expected = map[int]trackResult{
		1: {
			id:            1,
			title:         "Title 1",
			musicBrainzID: "MB1",
			albums: []albumResult{
				{
					id:    1,
					title: "Album 1",
				},
				{
					id:    2,
					title: "Album 2",
				},
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
				{
					id:   4,
					name: "Artist 4",
				},
			},
		},
	}

	got = readDB(t, db)
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("\nExpected:\n%#v\ngot:\n%#v", expected, got)
	}

	//////////////////////////////////////////////////////////////////////////

	t.Log("New MBID, existing title + primary artist combo")
	t.Log("Should be an entirely new track")

	input = []track.Track{
		track.New(
			"Title 1",
			"Album 1",
			"Artist 1",
			[]string{},
			"MB2",
		),
	}

	SaveTracks(db, input)

	expected = map[int]trackResult{
		2: {
			id:            2,
			title:         "Title 1",
			musicBrainzID: "MB2",
			albums: []albumResult{
				{
					id:    3,
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

	got = readDB(t, db)
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("\nExpected:\n%#v\ngot:\n%#v", expected, got)
	}

	//////////////////////////////////////////////////////////////////////////

	t.Log("Track with existing MusicBrainz ID but different primary artist")
	t.Log("    Existing primary artist should be uncoupled from track, replaced by new p. artist")

	t.Log("TODO")

	//////////////////////////////////////////////////////////////////////////

	t.Log("Track with existing MusicBrainz ID but different title")
	t.Log("Title should be overwritten")

	t.Log("TODO")

	//////////////////////////////////////////////////////////////////////////

	// TODO
	// Existing MBID, new title
	// Existing MBID, new primary artist
	// New MBID, duplicate (title + primary artist)
	// No MBID, completely new data
	// No MBID, duplicate (title + primary artist)

	//////////////////////////////////////////////////////////////////////////

	// input = []track.Track{
	// 	track.New(
	// 		"Title 1",
	// 		"Album 1",
	// 		"Artist 1",
	// 		[]string{},
	// 		"",
	// 	),
	// 	// Complete duplicate
	// 	track.New(
	// 		"Title 1",
	// 		"Album 1",
	// 		"Artist 1",
	// 		[]string{},
	// 		"",
	// 	),
	// 	// Duplicate with secondary artists included
	// 	// TODO
	// }

	// SaveTracks(db, input)

	// expected := map[int]trackResult{
	// 	1: {
	// 		id:    1,
	// 		title: "Title 1",
	// 		albums: []albumResult{
	// 			{
	// 				id:    1,
	// 				title: "Album 1",
	// 			},
	// 		},
	// 		primaryArtist: artistResult{
	// 			id:   1,
	// 			name: "Artist 1",
	// 		},
	// 		otherArtists: []artistResult{},
	// 	},
	// }

	// got := readDB(t, db)
	// if !reflect.DeepEqual(expected, got) {
	// 	t.Errorf("\nExpected:\n%#v\ngot:\n%#v", expected, got)
	// }

	// //////////////////////////////////////////////////////////////////////////

	// // New track title for existing album & artist
	// input = []track.Track{
	// 	track.New(
	// 		"Title 2",
	// 		"Album 1",
	// 		"Artist 1",
	// 		[]string{},
	// 		"",
	// 	),
	// }

	// SaveTracks(db, input)

	// expected[2] = trackResult{
	// 	id:    2,
	// 	title: "Title 2",
	// 	albums: []albumResult{
	// 		{
	// 			id:    1,
	// 			title: "Album 1",
	// 		},
	// 	},
	// 	primaryArtist: artistResult{
	// 		id:   1,
	// 		name: "Artist 1",
	// 	},
	// 	otherArtists: []artistResult{},
	// }

	// got = readDB(t, db)
	// if !reflect.DeepEqual(expected, got) {
	// 	t.Errorf("\nExpected:\n%#v\ngot:\n%#v", expected, got)
	// }

	// //////////////////////////////////////////////////////////////////////////

	// // New artist with same album title & track title as existing one
	// input = []track.Track{
	// 	track.New(
	// 		"Title 1",
	// 		"Album 1",
	// 		"Artist 2",
	// 		[]string{},
	// 		"",
	// 	),
	// }

	// SaveTracks(db, input)

	// // Track with duplicate name stored with new ID,
	// // album with duplicate name stored with new ID
	// expected[3] = trackResult{
	// 	id:    3,
	// 	title: "Title 1",
	// 	albums: []albumResult{
	// 		{
	// 			id:    2,
	// 			title: "Album 1",
	// 		},
	// 	},
	// 	primaryArtist: artistResult{
	// 		id:   2,
	// 		name: "Artist 2",
	// 	},
	// 	otherArtists: []artistResult{},
	// }

	// got = readDB(t, db)
	// if !reflect.DeepEqual(expected, got) {
	// 	t.Errorf("\nExpected:\n%#v\ngot:\n%#v", expected, got)
	// }

	// //////////////////////////////////////////////////////////////////////////

	// // New album for existing artist, duplicate track title
	// // (= artist having same track on different albums, e.g. original
	// // album then greatest hits album)
	// input = []track.Track{
	// 	track.New(
	// 		"Title 1",
	// 		"Album 2",
	// 		"Artist 1",
	// 		[]string{},
	// 		"",
	// 	),
	// }

	// SaveTracks(db, input)

	// // New album data stored, but track not added
	// expected[1] = trackResult{
	// 	id:    1,
	// 	title: "Title 1",
	// 	albums: []albumResult{
	// 		{
	// 			id:    1,
	// 			title: "Album 1",
	// 		},
	// 		{
	// 			id:    3,
	// 			title: "Album 2",
	// 		},
	// 	},
	// 	primaryArtist: artistResult{
	// 		id:   1,
	// 		name: "Artist 1",
	// 	},
	// 	otherArtists: []artistResult{},
	// }

	// got = readDB(t, db)
	// if !reflect.DeepEqual(expected, got) {
	// 	t.Errorf("\nExpected:\n%#v\ngot:\n%#v", expected, got)
	// }

	// //////////////////////////////////////////////////////////////////////////

	// // New track title, new artist, existing album title
	// input = []track.Track{
	// 	track.New(
	// 		"Title 3",
	// 		"Album 1",
	// 		"Artist 3",
	// 		[]string{},
	// 		"",
	// 	),
	// }

	// SaveTracks(db, input)

	// // New track & album data stored
	// expected[4] = trackResult{
	// 	id:    4,
	// 	title: "Title 3",
	// 	albums: []albumResult{
	// 		{
	// 			id:    4,
	// 			title: "Album 1",
	// 		},
	// 	},
	// 	primaryArtist: artistResult{
	// 		id:   3,
	// 		name: "Artist 3",
	// 	},
	// 	otherArtists: []artistResult{},
	// }

	// got = readDB(t, db)
	// if !reflect.DeepEqual(expected, got) {
	// 	t.Errorf("\nExpected:\n%#v\ngot:\n%#v", expected, got)
	// }

	// //////////////////////////////////////////////////////////////////////////

	// // New track title, existing artist, new album title
	// input = []track.Track{
	// 	track.New(
	// 		"Title 4",
	// 		"Album 3",
	// 		"Artist 1",
	// 		[]string{},
	// 		"",
	// 	),
	// }

	// SaveTracks(db, input)

	// // New track & album data stored
	// expected[5] = trackResult{
	// 	id:    5,
	// 	title: "Title 4",
	// 	albums: []albumResult{
	// 		{
	// 			id:    5,
	// 			title: "Album 3",
	// 		},
	// 	},
	// 	primaryArtist: artistResult{
	// 		id:   1,
	// 		name: "Artist 1",
	// 	},
	// 	otherArtists: []artistResult{},
	// }

	// got = readDB(t, db)
	// if !reflect.DeepEqual(expected, got) {
	// 	t.Errorf("\nExpected:\n%#v\ngot:\n%#v", expected, got)
	// }

	// //////////////////////////////////////////////////////////////////////////

	// // Existing track title, new artist, new album title
	// input = []track.Track{
	// 	track.New(
	// 		"Title 1",
	// 		"Album 4",
	// 		"Artist 4",
	// 		[]string{},
	// 		"",
	// 	),
	// }

	// SaveTracks(db, input)

	// // New data stored for everything
	// expected[6] = trackResult{
	// 	id:    6,
	// 	title: "Title 1",
	// 	albums: []albumResult{
	// 		{
	// 			id:    6,
	// 			title: "Album 4",
	// 		},
	// 	},
	// 	primaryArtist: artistResult{
	// 		id:   4,
	// 		name: "Artist 4",
	// 	},
	// 	otherArtists: []artistResult{},
	// }

	// got = readDB(t, db)
	// if !reflect.DeepEqual(expected, got) {
	// 	t.Errorf("\nExpected:\n%#v\ngot:\n%#v", expected, got)
	// }

	// //////////////////////////////////////////////////////////////////////////

	// // Add track with secondary artists
	// input = []track.Track{
	// 	track.New(
	// 		"Title 101",
	// 		"Album 101",
	// 		"Artist 101",
	// 		[]string{
	// 			"Artist 201",
	// 			"Artist 301",
	// 		},
	// 		"",
	// 	),
	// }

	// SaveTracks(db, input)

	// // New data stored for everything
	// expected[7] = trackResult{
	// 	id:    7,
	// 	title: "Title 101",
	// 	albums: []albumResult{
	// 		{
	// 			id:    7,
	// 			title: "Album 101",
	// 		},
	// 	},
	// 	primaryArtist: artistResult{
	// 		id:   5,
	// 		name: "Artist 101",
	// 	},
	// 	otherArtists: []artistResult{
	// 		{
	// 			id:   6,
	// 			name: "Artist 201",
	// 		},
	// 		{
	// 			id:   7,
	// 			name: "Artist 301",
	// 		},
	// 	},
	// }

	// got = readDB(t, db)
	// if !reflect.DeepEqual(expected, got) {
	// 	t.Errorf("\nExpected:\n%#v\ngot:\n%#v", expected, got)
	// }

	// //////////////////////////////////////////////////////////////////////////

	// // Add new artists as secondary artists to existing track
	// input = []track.Track{
	// 	track.New(
	// 		"Title 1",
	// 		"Album 1",
	// 		"Artist 1",
	// 		[]string{
	// 			"Artist 401",
	// 			"Artist 501",
	// 		},
	// 		"",
	// 	),
	// }

	// SaveTracks(db, input)

	// // Track updated to have secondary artists
	// tmp := expected[1]
	// tmp.otherArtists = []artistResult{
	// 	{
	// 		id:   8,
	// 		name: "Artist 401",
	// 	},
	// 	{
	// 		id:   9,
	// 		name: "Artist 501",
	// 	},
	// }
	// expected[1] = tmp

	// got = readDB(t, db)
	// if !reflect.DeepEqual(expected, got) {
	// 	t.Errorf("\nExpected:\n%#v\ngot:\n%#v", expected, got)
	// }

	// TODO
	// Save brand new with musicbrainz ID
	// Save duplicate with musicbrainz ID
	// Save different with same musicbrainz ID?
	// Update to have musicbrainz ID?
	// Update track that has a musicbrainz ID
}

func readDB(t *testing.T, db *sql.DB) map[int]trackResult {
	rows, err := db.Query(
		`SELECT t.id,
		        t.title,
		        IFNULL(t.musicbrainz_id, ""),
		        ar.id,
		        ar.name
		  FROM tracks       t
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
		var artist artistResult

		if err = rows.Scan(
			&track.id,
			&track.title,
			&track.musicBrainzID,
			&artist.id,
			&artist.name,
		); err != nil {
			t.Fatal(err)
		}

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
