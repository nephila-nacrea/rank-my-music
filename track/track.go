package track

const StartingRanking = 1000

type Album struct {
	InternalID    int
	MusicBrainzID string
	Title         string
}

type Artist struct {
	InternalID    int
	MusicBrainzID string
	Name          string
}

type Track struct {
	InternalID    int
	MusicBrainzID string
	Title         string

	Albums []Album

	PrimaryArtist Artist
	OtherArtists  []Artist

	Ranking float64
}

func New(track Track) Track {
	// Provide defaults

	return Track{
		MusicBrainzID: track.MusicBrainzID,
		Title:         track.Title,

		Albums: track.Albums,

		PrimaryArtist: track.PrimaryArtist,
		OtherArtists:  track.OtherArtists,

		Ranking: StartingRanking,
	}
}
