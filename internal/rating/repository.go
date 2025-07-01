package rating

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Repository interface {
	GetRatingByMemberId(ctx context.Context, memberId uuid.UUID) (*Rating, error)
	GetRatingsByMemberIds(ctx context.Context, memberIds []uuid.UUID) ([]Rating, error)
	CreateRating(ctx context.Context, rating *Rating) (uuid.UUID, error)
	UpdateRating(ctx context.Context, ratings *Rating) error
	UpdateRatings(ctx context.Context, ratings []Rating) error
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetRatingByMemberId(ctx context.Context, memberId uuid.UUID) (*Rating, error) {
	var rating *Rating

	err := r.db.GetContext(ctx, rating, "SELECT * FROM ratings WHERE member_id = $1", memberId)
	if err != nil {
		return nil, err
	}

	return rating, nil
}

func (r *repository) GetRatingsByMemberIds(ctx context.Context, memberIds []uuid.UUID) ([]Rating, error) {
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

func (r *repository) CreateRating(ctx context.Context, rating *Rating) (uuid.UUID, error) {
	var id uuid.UUID

	err := r.db.QueryRowContext(ctx,
		"INSERT INTO ratings (member_id, game_id, mu, sigma) VALUES ($1, $2, $3, $4) RETURNING id",
		rating.MemberID, rating.GameID, rating.Mu, rating.Sigma).Scan(&id)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func (r *repository) UpdateRatings(ctx context.Context, ratings []Rating) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	for _, rating := range ratings {
		_, err = tx.ExecContext(ctx,
			"UPDATE ratings SET mu = $1, sigma = $2 WHERE member_id = $3 AND game_id = $4",
			rating.Mu, rating.Sigma, rating.MemberID, rating.GameID,
		)
		if err != nil {
			if errRollback := tx.Rollback(); errRollback != nil {
				return errRollback
			}

			return err
		}
	}

	return tx.Commit()
}

func (r *repository) UpdateRating(ctx context.Context, rating *Rating) error {
	_, err := r.db.ExecContext(ctx,
		"UPDATE ratings SET mu = $1, sigma = $2 WHERE member_id = $3 AND game_id = $4",
		rating.Mu, rating.Sigma, rating.MemberID, rating.GameID,
	)
	if err != nil {
		return err
	}

	return nil
}
