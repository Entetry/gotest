package postgre

import (
	"context"
	"entetry/gotest/internal/model"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type User struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *Company {
	return &Company{
		db: db,
	}
}

func (u *User) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	var user *model.User
	err := u.db.QueryRow(ctx, "SELECT id, name FROM user WHERE username = $1", username).
		Scan(&user.ID, &user.Username, &user.PasswordHash)
	return user, err
}

func (u *User) Create(ctx context.Context, user *model.User) (uuid.UUID, error) {
	user.ID = uuid.New()
	err := u.db.QueryRow(ctx, "INSERT INTO user(id, username, passwordHash) VALUES ($1, $2) RETURNING id, username, passwordHash;",
		user.ID, user.Username, user.PasswordHash).Scan(&user.ID, &user.Username, &user.PasswordHash)
	if err != nil {
		return user.ID, fmt.Errorf("cannot create User: %v", err)
	}
	return user.ID, err
}

func (u *User) Update(ctx context.Context, user *model.User) error {
	_, err := u.db.Exec(ctx, "UPDATE user SET username = $2"+
		", passwordHash = $3 WHERE id=$1 RETURNING id, username, passwordHash;",
		user.ID, user.Username, user.PasswordHash)
	if err != nil {
		return fmt.Errorf("cannot update User: %v", err)
	}
	return err
}

func (u *User) Delete(ctx context.Context, uuid uuid.UUID) error {
	_, err := u.db.Exec(ctx, "DELETE FROM user WHERE id = $1", uuid)
	if err != nil {
		return fmt.Errorf("cannot delete User: %v", err)
	}
	return nil
}
