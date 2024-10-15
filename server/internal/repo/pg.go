package repo

import (
	"context"
	"fmt"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
	mu sync.Mutex
}

func NewPostgresRepo(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) InsertClientData(ctx context.Context, data string, socketId string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil || err != nil {
			if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
				return
			}
		}
	}()

	query := `
		INSERT INTO client_data (data, socket_id)
		VALUES ($1, $2)
	`

	_, err = tx.Exec(ctx, query, data, socketId)
	if err != nil {
		return fmt.Errorf("failed to add client data: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
