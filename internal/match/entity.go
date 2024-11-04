package match

type Match struct {
	ID      int
	ClubID  int   `db:"club_id"`
	GameID  int   `db:"game_id"`
	TeamIDs []int `db:"team_ids"`
	Sets    []string
	Result  rune
}

type Team struct {
	ID         int
	ClubID     int   `db:"club_id"`
	MembersIds []int `db:"members_ids"`
}
