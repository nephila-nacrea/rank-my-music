package track

type Track struct {
	Title   string
	Artists []string
	Album   string

	Ranking float64 // Default to 1000?
}

func New(title string) Track {
	return Track{
		Title: title,

		Ranking: 1000,
	}
}
