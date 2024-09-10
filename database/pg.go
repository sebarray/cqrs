package database

import (
	"context"
	model "cqrs/models"
	"database/sql"

	_ "github.com/lib/pq"
)

type PgRepository struct {
	db *sql.DB
}

func NewPgRepository(url string) (*PgRepository, error) {
	db, err := sql.Open("postgres", url)

	if err != nil {
		return nil, err
	}

	return &PgRepository{db: db}, nil

}

func (repo *PgRepository) Close() {
	repo.db.Close()
}

func (repo *PgRepository) InsertFeed(ctx context.Context, feed model.Feed) error {

	_, err := repo.db.ExecContext(ctx, "INSERT INTO feeds(id, title, description) VALUES($1, $2, $3)", feed.ID, feed.Title, feed.Description)

	if err != nil {
		return err
	}

	return nil
}

func (repo *PgRepository) ListFeed(ctx context.Context) ([]model.Feed, error) {

	row, err := repo.db.QueryContext(ctx, "SELECT id, title, description, created_at FROM feeds")
	if err != nil {
		return nil, err
	}
	defer row.Close()
	var feeds []model.Feed
	for row.Next() {
		var feed model.Feed
		err := row.Scan(&feed.ID,
			&feed.Title,
			&feed.Description,
			&feed.CreatedAt)
		if err != nil {
			return nil, err
		}
		feeds = append(feeds, feed)
	}
	return feeds, nil
}
