package rating

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type Repository interface {
	GetRatingByMemberId(ctx context.Context, memberId int) (*Rating, error)
	GetRatingsByMemberIds(ctx context.Context, memberIds []int) ([]Rating, error)
	CreateRating(ctx context.Context, rating *Rating) (int, error)
	UpdateRating(ctx context.Context, ratings *Rating) error
	UpdateRatings(ctx context.Context, ratings []Rating) error
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetRatingByMemberId(ctx context.Context, memberId int) (*Rating, error) {
	var rating *Rating

	err := r.db.GetContext(ctx, rating, "SELECT * FROM ratings WHERE member_id = $1", memberId)
	if err != nil {
		return nil, err
	}

	return rating, nil
}

func (r *repository) GetRatingsByMemberIds(ctx context.Context, memberIds []int) ([]Rating, error) {
	var ratings []Rating

	query, args, err := sqlx.In("SELECT * FROM ratings WHERE member_id IN (?)", memberIds)
	if err != nil {
		return nil, err
	}

	query = r.db.Rebind(query)
	err = r.db.SelectContext(ctx, &ratings, query, args...)
	if err != nil {
		return nil, err
	}

	return ratings, nil
}

func (r *repository) CreateRating(ctx context.Context, rating *Rating) (int, error) {
	result, err := r.db.ExecContext(ctx,
		"INSERT INTO ratings (member_id, game_id, value) VALUES ($1, $2)",
		rating.MemberID, rating.GameID, rating.Value,
	)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (r *repository) UpdateRatings(ctx context.Context, ratings []Rating) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	for _, rating := range ratings {
		_, err = tx.ExecContext(ctx,
			"UPDATE ratings SET value = $1 WHERE member_id = $2 AND game_id = $3",
			rating.Value, rating.MemberID, rating.GameID,
		)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (r *repository) UpdateRating(ctx context.Context, rating *Rating) error {
	_, err := r.db.ExecContext(ctx,
		"UPDATE ratings SET value = $1 WHERE member_id = $2 AND game_id = $3",
		rating.Value, rating.MemberID, rating.GameID,
	)
	if err != nil {
		return err
	}

	return nil
}
