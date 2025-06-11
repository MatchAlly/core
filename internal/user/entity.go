package user

import "time"

type User struct {
	ID        int       `db:"id"`
	Email     string    `db:"email"`
	Name      string    `db:"name"`
	Hash      string    `db:"hash"`
	CreatedAt time.Time `db:"created_at"`
	LastLogin time.Time `db:"last_login"`
	UpdatedAt time.Time `db:"updated_at"`
}
