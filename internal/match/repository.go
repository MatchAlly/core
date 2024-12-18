package match

import (
	"context"
	"core/internal/member"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type Repository interface {
	CreateMatch(ctx context.Context, m *Match) (int, error)
	GetMatches(ctx context.Context, clubID int) ([]Match, error)
	GetMatchesByGame(ctx context.Context, clubID, gameID int) ([]Match, error)
	GetTeam(ctx context.Context, teamID int) (*Team, error)
	CreateTeam(ctx context.Context, clubID int, memberIDs []int) (int, error)
	TeamOfMembersExists(ctx context.Context, clubID int, memberIDs []int) (bool, int, error)
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) CreateMatch(ctx context.Context, m *Match) (int, error) {
	var matchID int

	const query = `INSERT INTO matches (club_id, game_id) VALUES ($1, $2) RETURNING id`
	if err := r.db.GetContext(ctx, &matchID, query, m.ClubID, m.GameID); err != nil {
		return 0, err
	}

	return matchID, nil
}

func (r *repository) GetMatches(ctx context.Context, clubID int) ([]Match, error) {
	matchesMap := make(map[int]*Match)

	rows, err := r.db.QueryxContext(ctx, `
			SELECT
					m.id AS match_id,
					m.club_id AS match_club_id,
					m.game_id AS match_game_id,
					m.mode AS match_mode,
					m.ranked AS match_ranked,
					m.sets AS match_sets,
					m.created_at AS match_created_at,
					t.id AS team_id,
					t.club_id AS team_club_id,
					mem.id AS member_id,
					mem.club_id AS member_club_id,
					mem.user_id AS member_user_id,
					mem.display_name AS member_display_name,
					mem.role AS member_role
			FROM matches m
			LEFT JOIN match_teams mt ON m.id = mt.match_id
			LEFT JOIN teams t ON mt.team_id = t.id
			LEFT JOIN team_members tm ON t.id = tm.team_id
			LEFT JOIN members mem ON tm.member_id = mem.id
			WHERE m.club_id = $1
			ORDER BY m.id, t.id, mem.id; -- Crucial for correct grouping
	`, clubID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var m Match
		var t Team
		var mem member.Member

		err = rows.Scan(
			&m.ID, &m.ClubID, &m.GameID, &m.Gamemode, &m.Ranked, &m.Sets, &m.CreatedAt,
			&t.ID, &t.ClubID,
			&mem.ID, &mem.ClubID, &mem.UserID, &mem.DisplayName, &mem.Role,
		)
		if err != nil {
			return nil, err
		}

		match, ok := matchesMap[m.ID]
		if !ok {
			match = &m
			match.Teams = make([]Team, 0)
			matchesMap[m.ID] = match
		}

		if t.ID != 0 {
			teamExists := false
			for i := range match.Teams {
				if match.Teams[i].ID == t.ID {
					match.Teams[i].Members = append(match.Teams[i].Members, mem)
					teamExists = true
					break
				}
			}
			if !teamExists {
				t.Members = []member.Member{mem}
				match.Teams = append(match.Teams, t)
			}
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	matches := make([]Match, 0, len(matchesMap))
	for _, matchPtr := range matchesMap {
		matches = append(matches, *matchPtr)
	}

	return matches, nil
}

func (r *repository) GetMatchesByGame(ctx context.Context, clubID int, gameID int) ([]Match, error) {
	matchesMap := make(map[int]*Match)

	rows, err := r.db.QueryxContext(ctx, `
        SELECT
            m.id AS match_id,
            m.club_id AS match_club_id,
            m.game_id AS match_game_id,
            m.mode AS match_mode,
            m.ranked AS match_ranked,
            m.sets AS match_sets,
            m.created_at AS match_created_at,
            t.id AS team_id,
            t.club_id AS team_club_id,
            mem.id AS member_id,
            mem.club_id AS member_club_id,
            mem.user_id AS member_user_id,
            mem.display_name AS member_display_name,
            mem.role AS member_role
        FROM matches m
        LEFT JOIN match_teams mt ON m.id = mt.match_id
        LEFT JOIN teams t ON mt.team_id = t.id
        LEFT JOIN team_members tm ON t.id = tm.team_id
        LEFT JOIN members mem ON tm.member_id = mem.id
        WHERE m.club_id = $1 AND m.game_id = $2
        ORDER BY m.id, t.id, mem.id; -- Crucial for correct grouping
    `, clubID, gameID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var m Match
		var t Team
		var mem member.Member

		err = rows.Scan(
			&m.ID, &m.ClubID, &m.GameID, &m.Gamemode, &m.Ranked, &m.Sets, &m.CreatedAt,
			&t.ID, &t.ClubID,
			&mem.ID, &mem.ClubID, &mem.UserID, &mem.DisplayName, &mem.Role,
		)
		if err != nil {
			return nil, err
		}

		match, ok := matchesMap[m.ID]
		if !ok {
			match = &m
			match.Teams = make([]Team, 0)
			matchesMap[m.ID] = match
		}

		if t.ID != 0 {
			teamExists := false
			for i := range match.Teams {
				if match.Teams[i].ID == t.ID {
					match.Teams[i].Members = append(match.Teams[i].Members, mem)
					teamExists = true
					break
				}
			}
			if !teamExists {
				t.Members = []member.Member{mem}
				match.Teams = append(match.Teams, t)
			}
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	matches := make([]Match, 0, len(matchesMap))
	for _, matchPtr := range matchesMap {
		matches = append(matches, *matchPtr)
	}

	return matches, nil
}

func (r *repository) GetTeam(ctx context.Context, teamID int) (*Team, error) {
	const teamQuery = `SELECT * FROM teams WHERE id = $1`

	var team Team
	err := r.db.GetContext(ctx, &team, teamQuery, teamID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	const membersQuery = `
        SELECT *
        FROM members m
        JOIN team_members tm ON tm.member_id = m.id
        WHERE tm.team_id = $1
        ORDER BY m.display_name`

	var members []member.Member
	err = r.db.SelectContext(ctx, &members, membersQuery, teamID)
	if err != nil {
		return nil, err
	}

	team.Members = members

	return &team, nil
}

func (r *repository) CreateTeam(ctx context.Context, clubID int, memberIDs []int) (int, error) {
	var teamID int

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return 0, err
	}

	const queryTeam = `INSERT INTO teams (club_id) VALUES ($1) RETURNING id`
	if err := tx.GetContext(ctx, &teamID, queryTeam, clubID); err != nil {
		tx.Rollback()
		return 0, err
	}

	type teamMember struct {
		TeamID   int `db:"team_id"`
		MemberID int `db:"member_id"`
	}

	teamMembers := make([]teamMember, len(memberIDs))
	for i, memberID := range memberIDs {
		teamMembers[i] = teamMember{
			TeamID:   teamID,
			MemberID: memberID,
		}
	}

	const queryTeamMembers = `INSERT INTO team_members (team_id, member_id) VALUES (:team_id, :member_id)`
	if _, err = tx.NamedExec(queryTeamMembers, teamMembers); err != nil {
		tx.Rollback()
		return 0, err
	}

	return teamID, nil
}

func (r *repository) TeamOfMembersExists(ctx context.Context, clubID int, memberIDs []int) (bool, int, error) {
	const query = `
		WITH team_candidates AS (
			SELECT t.id as team_id
			FROM teams t
			WHERE t.club_id = $1
		),
		member_counts AS (
			SELECT 
				tm.team_id,
				COUNT(*) as member_count,
				-- Count how many of the specified members are in this team
				COUNT(CASE WHEN tm.member_id = ANY($2) THEN 1 END) as matching_members
			FROM team_candidates tc 
			JOIN team_members tm ON tm.team_id = tc.team_id
			GROUP BY tm.team_id
		)
		SELECT team_id
		FROM member_counts
		WHERE member_count = array_length($2, 1)
		AND matching_members = array_length($2, 1)`

	var teamID int
	err := r.db.GetContext(ctx, &teamID, query, clubID, pq.Array(memberIDs))
	if err != nil {
		if err == sql.ErrNoRows {
			return false, 0, nil
		}
		return false, 0, err
	}

	return true, teamID, nil
}
