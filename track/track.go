package track

const StartingRanking = 1000

type Track struct {
	Album         string
	MusicBrainzID string
	OtherArtists  []string
	PrimaryArtist string
	Title         string

	Ranking float64 // Default to 1000?
}

func New(
	title string,
	album string,
	primaryArtist string,
	otherArtists []string,
	musicBrainzID string,
) Track {
	return Track{
		Album:         album,
		MusicBrainzID: musicBrainzID,
		PrimaryArtist: primaryArtist,
		OtherArtists:  otherArtists,
		Title:         title,

		Ranking: StartingRanking,
	}
}
