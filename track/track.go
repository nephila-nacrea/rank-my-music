package track

const StartingRanking = 1000

type Track struct {
	Album         string
	PrimaryArtist string
	OtherArtists  []string
	Title         string

	Ranking float64 // Default to 1000?
}

func New(
	title string,
	album string,
	primaryArtist string,
	otherArtists []string,
) Track {
	return Track{
		Album:         album,
		PrimaryArtist: primaryArtist,
		OtherArtists:  otherArtists,
		Title:         title,

		Ranking: StartingRanking,
	}
}
