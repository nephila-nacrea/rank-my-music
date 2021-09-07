package elo

import "math"

// See https://en.wikipedia.org/wiki/Elo_rating_system#Mathematical_details
// and https://www.geeksforgeeks.org/elo-rating-algorithm/

const K = 32

type Elo struct {
	CurrentRanking float64
	Score          float64 // 0 = loss, 0.5 = draw, 1 = win
}

func CalculateNewRankings(eloA Elo, eloB Elo) (float64, float64) {
	newRankA := eloA.CurrentRanking +
		K*(eloA.Score-ExpectedScore(
			eloA.CurrentRanking, eloB.CurrentRanking,
		))

	newRankB := eloB.CurrentRanking +
		K*(eloB.Score-ExpectedScore(
			eloB.CurrentRanking, eloA.CurrentRanking,
		))

	return newRankA, newRankB
}

func ExpectedScore(rankSelf float64, rankOpponent float64) float64 {
	return 1 / (1 + math.Pow(
		10,
		(rankOpponent-rankSelf)/400,
	))
}
