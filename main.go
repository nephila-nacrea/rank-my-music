package main

import (
	"fmt"
	"math/rand"

	"example.com/rank-my-music/elo"
	"example.com/rank-my-music/track"
)

func main() {
	tracks := [...]track.Track{
		track.New("Crucify"),
		track.New("Girl"),
		track.New("Silent All These Years"),
		track.New("Precious Things"),
		track.New("Winter"),
		track.New("Happy Phantom"),
		track.New("China"),
		track.New("Leather"),
		track.New("Mother"),
		track.New("Tear in Your Hand"),
		track.New("Me and a Gun"),
		track.New("Little Earthquakes"),
	}

	fmt.Printf("%v\n", tracks)

	for i := 1; i <= 1000; i++ {
		trackA := &tracks[rand.Intn(len(tracks))]
		trackB := &tracks[rand.Intn(len(tracks))]

		// Make B always win. Just want to see how Elo works for now.
		trackA.Ranking, trackB.Ranking = elo.CalculateNewRankings(
			elo.Elo{
				CurrentRanking: trackA.Ranking,
				Score:          0,
			},
			elo.Elo{
				CurrentRanking: trackB.Ranking,
				Score:          1,
			},
		)

		fmt.Println(trackA)
		fmt.Println(trackB)
		fmt.Println()
	}

	fmt.Printf("%v\n", tracks)
}
