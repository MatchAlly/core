package leaderboard

type LeaderboardType string

const (
	TypeWins   LeaderboardType = "wins"
	TypeStreak LeaderboardType = "streak"
	TypeRating LeaderboardType = "rating"
)

type Entry struct {
	Value  float64 `json:"value"`
	UserId uint    `json:"user_id"`
	Name   string  `json:"name"`
}

type Leaderboard struct {
	Type    LeaderboardType `json:"type"`
	Entries []Entry         `json:"entries"`
}
