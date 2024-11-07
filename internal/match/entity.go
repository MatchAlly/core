package match

import (
	"core/internal/member"
	"time"
)

type Result rune

const (
	ResultWin  Result = 'W'
	ResultLoss Result = 'L'
	ResultDraw Result = 'D'
)

type Match struct {
	ID        int       `json:"id" db:"id"`
	ClubID    int       `json:"club_id" db:"club_id"`
	GameID    int       `json:"game_id" db:"game_id"`
	Result    rune      `json:"result" db:"result"`
	Sets      []string  `json:"sets" db:"sets"`
	Teams     []Team    `json:"teams,omitempty"` // Loaded via joins
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type Team struct {
	ID      int             `json:"id" db:"id"`
	ClubID  int             `json:"club_id" db:"club_id"`
	Members []member.Member `json:"members,omitempty"` // Loaded via joins
}
