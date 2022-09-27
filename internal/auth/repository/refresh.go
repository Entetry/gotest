package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	log "github.com/sirupsen/logrus"
)

type Refresh struct {
	db *pgxpool.Pool
}

func NewRefresh(db *pgxpool.Pool) *Refresh {
	return &Refresh{db: db}
}

func (r *Refresh) GetByUserID(ctx context.Context, userID int) (string, error) {
	var hash string

	err := r.db.QueryRow(ctx, "SELECT hash FROM refresh WHERE user_id = $1", userID).Scan(&hash)
	if err != nil {
		log.Errorf("QueryRow in GetByUserID func failed: %v\n", err)
		return "", err
	}

	return hash, nil
}

func (r *Refresh) Create(ctx context.Context, userID int, hash string) error {
	err := r.db.QueryRow(ctx, "INSERT INTO refresh (user_id, hash) VALUES ($1, $2) RETURNING user_id, hash", userID, hash)
	if err != nil {
		log.Errorf("Insert in Create func failed: %v\n", err)
		return fmt.Errorf("cannot create Refresh: %v", err)
	}
	return nil
}

func (r *Refresh) Delete(ctx context.Context, userID int) error {
	_, err := r.db.Exec(ctx, "DELETE FROM refresh WHERE user_id = $1", userID)
	if err != nil {
		log.Errorf("Exec in Delete func failed: %v\n", err)
		return err
	}

	return nil
}
