package repo

import (
	"testing"

	"github.com/nephila-nacrea/rank-my-music/test_utils"
	"github.com/nephila-nacrea/rank-my-music/track"
)

func TestSaveTracks(t *testing.T) {
	tracks := []track.Track{
		track.New(
			"Title 1",
			"Album 1",
			[]string{"Artist 1", "Artist 2", "Artist 3"},
		),
	}

	t.Logf("%v", tracks)

	db := test_utils.DBSetup()

	SaveTracks(db, tracks)
}
