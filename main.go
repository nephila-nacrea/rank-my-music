package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"

	"github.com/dhowden/tag"
	"github.com/nephila-nacrea/rank-my-music/elo"
	"github.com/nephila-nacrea/rank-my-music/track"
)

func init() {
	log.SetFlags(log.Llongfile)
}

func main() {
	// TODO
	// Handle duplicates
	// What if no metadata? (e.g. wma)

	var tracks []track.Track

	var files []string
	err := filepath.Walk(
		"/media/vmihell-hale/Arch/home/vmihell-hale/Music/Tori Amos",
		func(path string, info os.FileInfo, err error) error {
			if !info.IsDir() {
				files = append(files, path)
			}
			return nil
		},
	)
	if err != nil {
		log.Fatalln(err)
	}

	for _, f := range files {
		file, err := os.Open(f)
		if err != nil {
			log.Println(f+": ", err)
		}

		meta, err := tag.ReadFrom(file)
		if err != nil {
			log.Println(f+": ", err)
		} else {
			// for _, v := range []string{
			// 	meta.Album(),
			// 	meta.AlbumArtist(),
			// 	meta.Artist(),
			// 	meta.Composer(),
			// 	meta.Genre(),
			// 	meta.Title(),
			// } {
			// 	fmt.Println(v)
			// }

			tracks = append(
				tracks,
				track.New(
					meta.Title(),
					meta.Album(),
					[]string{
						meta.Artist(), meta.AlbumArtist(), meta.Composer(),
					},
				),
			)
		}
	}

	for i := 1; i <= 1000; i++ {
		// Use pointers so ranking is updated for originals
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
