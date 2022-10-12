package postgres

import (
	dbInterface "GoNews/pkg/storage"
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
)

// Storage Хранилище данных.
type Storage struct {
	Db *pgxpool.Pool
}

// New Конструктор, принимает строку подключения к БД.
func New(ctx context.Context, constr string) (*Storage, error) {
	db, err := pgxpool.Connect(ctx, constr)
	if err != nil {
		return nil, err
	}
	s := Storage{
		Db: db,
	}
	return &s, nil
}

func (pg *Storage) Posts() ([]dbInterface.Post, error) {
	query := `SELECT id, title, content, author_id, created_at FROM posts;`
	rows, err := pg.Db.Query(context.Background(), query)

	if err != nil {
		return nil, err
	}

	var posts []dbInterface.Post

	for rows.Next() {
		var p dbInterface.Post
		err = rows.Scan(
			&p.ID,
			&p.Title,
			&p.Content,
			&p.AuthorID,
			&p.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		posts = append(posts, p)
	}

	return posts, rows.Err()
}

func (pg *Storage) AddPost(p dbInterface.Post) error {
	query := `INSERT INTO posts(author_id, title, content, created_at) VALUES($1, $2, $3, $4);`
	_, err := pg.Db.Exec(context.Background(), query, p.AuthorID, p.Title, p.Content, p.CreatedAt)
	return err
}

func (pg *Storage) UpdatePost(p dbInterface.Post) error {
	query := `UPDATE posts set title = $3, content = $4 WHERE author_id = $2 AND id=$1;`
	_, err := pg.Db.Exec(context.Background(), query, p.ID, p.AuthorID, p.Title, p.Content)
	return err
}

func (pg *Storage) DeletePost(p dbInterface.Post) error {
	query := `DELETE FROM posts WHERE id = $1;`
	_, err := pg.Db.Exec(context.Background(), query, p.ID)
	return err
}
