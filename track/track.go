package track

type Track struct {
	Album   string
	Artists []string
	Title   string

	Ranking float64 // Default to 1000?
}

func New(title string, album string, artists []string) Track {
	return Track{
		Album:   album,
		Artists: artists,
		Title:   title,

		Ranking: 1000,
	}
}
