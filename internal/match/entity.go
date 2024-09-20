package match

type Result rune

const (
	TeamAWins Result = 'A'
	TeamBWins Result = 'B'
	Draw      Result = 'D'
)

type Match struct {
	ID      uint
	ClubID  uint   `db:"club_id"`
	GameID  uint   `db:"game_id"`
	TeamIDs []uint `db:"team_ids"`
	Sets    []string
	Result  Result
}

type Team struct {
	ID         uint
	ClubID     uint   `db:"club_id"`
	MembersIds []uint `db:"members_ids"`
}
