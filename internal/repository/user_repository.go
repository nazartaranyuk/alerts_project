package repository

import (
	"database/sql"
	"errors"
	"nazartaraniuk/alertsProject/internal/app/db"
	"nazartaraniuk/alertsProject/internal/domain"

	"golang.org/x/crypto/bcrypt"
)

var ErrCannotGenerateHash = errors.New("cannot generate password hash")
var ErrEmailAlreadyExists = errors.New("email already exists")
var ErrCannotCreateUser = errors.New("cannot create user")

type UserRepository struct {
	db *db.Database
}

func NewUserRepository(database *db.Database) *UserRepository {
	return &UserRepository{
		db: database,
	}
}

func (r *UserRepository) GetUserByEmail(email string) (*domain.User, error) {
	row := r.db.DB.QueryRow(`SELECT id, email, username, password_hash, created_at FROM users WHERE email=$1`, email)
	var user domain.User
	err := row.Scan(&user.ID, &user.Email, &user.Username, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) CreateUser(req domain.RegisterReq) (int64, error) {
	h, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, ErrCannotGenerateHash
	}
	var id int64
	err = r.db.DB.QueryRow(`INSERT INTO users(email, username, password_hash) VALUES ($1,$2,$3) RETURNING id`, req.Email, req.Username, string(h)).Scan(&id)
	if err != nil {
		if isUniqueViolation(err) {
			return 0, ErrEmailAlreadyExists
		}
		return 0, ErrCannotCreateUser
	}
	return id, nil
}

func isUniqueViolation(err error) bool {
	type state interface{ SQLState() string }
	if e, ok := err.(state); ok && e.SQLState() == "23505" {
		return true
	}
	return false
}
