package repo

import (
	"context"
	"database/sql"

	"lnkshrt/internal/domain/models"
)

type UserRepoInterface interface {
	CreateUser(c context.Context, user *models.User) error
	GetUser(c context.Context, userID int) (*models.User, error)
	GetUserByUsername(c context.Context, username string) (*models.User, error)
}

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) CreateUser(c context.Context, user *models.User) error {

	query := `INSERT INTO users (username, password, email, created_at) VALUES ($1, $2, $3, $4) RETURNING id`

	return r.db.QueryRowContext(c, query, user.Username, user.Password, user.Email, user.CreatedAt).Scan(&user.ID)

}

func (r *UserRepo) GetUser(c context.Context, userID int) (*models.User, error) {

	user := &models.User{ID: userID}

	query := `SELECT username, email, password, created_at FROM users WHERE id = $1`

	err := r.db.QueryRowContext(c, query, userID).Scan(&user.Username, &user.Email, &user.Password, &user.CreatedAt)

	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepo) GetUserByUsername(c context.Context, username string) (*models.User, error) {

	user := &models.User{}

	query := `SELECT id, username, email, password, created_at FROM users WHERE username = $1`

	err := r.db.QueryRowContext(c, query, username).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt)

	if err != nil {
		return nil, err
	}
	return user, nil

}
