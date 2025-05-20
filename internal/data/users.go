package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/obrikash/swimming_pool/internal/validator"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrDuplicateEmail = errors.New("duplicate email")
)

type password struct {
	plaintext *string
	hash      []byte
}

var AnonymousUser = &User{}

type User struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	FullName  string    `json:"full_name"`
	Email     string    `json:"email"`
	Password  password  `json:"-"`
	RoleID    uint8     `json:"role_id"`
	Image     string    `json:"image_url"`
}

type UserModel struct {
	DB *sql.DB
}

func (u *User) IsAnonymous() bool {
	return u == AnonymousUser
}

func (p *password) Set(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return err
	}

	p.plaintext = &plaintextPassword
	p.hash = hash

	return nil
}

func (p *password) Matches(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plaintextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}
	return true, nil
}

func ValidatePasswordPlaintext(v *validator.Validator, password string) {
	v.Check(password != "", "password", "must be provided")
	v.Check(len(password) >= 8, "password", "must be at 8 bytes long")
	v.Check(len(password) <= 72, "password", "must not be more than 72 bytes long")
}

func ValidateEmail(v *validator.Validator, email string) {
	v.Check(email != "", "email", "must be provided")
	v.Check(validator.Matches(email, validator.EmailRX), "email", "must be a valid email address")
}

func ValidateUser(v *validator.Validator, user *User) {
	v.Check(user.FullName != "", "name", "must be provided")
	v.Check(len(user.FullName) <= 500, "name", "must not be more than 500 bytes long")

	ValidateEmail(v, user.Email)
	if user.Password.plaintext != nil {
		ValidatePasswordPlaintext(v, *user.Password.plaintext)
	}

	if user.Password.hash == nil {
		panic("missing password hash for user")
	}
}

func (um UserModel) GetTrainers() ([]*User, error) {
	query := `SELECT t.user_id, full_name, u.image FROM users u JOIN trainers t ON t.user_id = u.id`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := um.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	trainers := []*User{}

	for rows.Next() {
		var trainer User

		err := rows.Scan(&trainer.ID, &trainer.FullName, &trainer.Image)
		if err != nil {
			return nil, err
		}

		trainers = append(trainers, &trainer)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return trainers, nil
}

func (um UserModel) GetTrainersForPools() (map[Pool][]*User, error) {
	query := `SELECT p.id AS "pool_id", p.name AS "pool_name", p.address, p.type AS "category",
	 u.id AS "trainer_id", u.full_name, u.image FROM trainers t JOIN pools p ON t.pool_id = p.id JOIN users u ON t.user_id = u.id ORDER BY p.name`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := um.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	trainers := map[Pool][]*User{}

	for rows.Next() {
		var pool Pool
		var trainer User

		err := rows.Scan(&pool.ID, &pool.Name, &pool.Address, &pool.PoolType, &trainer.ID, &trainer.FullName, &trainer.Image)
		if err != nil {
			return nil, err
		}

		trainers[pool] = append(trainers[pool], &trainer)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return trainers, nil
}

func (um UserModel) Insert(user *User) error {

	query := "INSERT INTO users (full_name, email, hashed_password, role_id) VALUES ($1, $2, $3, $4) RETURNING id, created_at"

	args := []any{user.FullName, user.Email, user.Password.hash, user.RoleID}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := um.DB.QueryRowContext(ctx, query, args...).Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		default:
			return err
		}
	}

	return nil
}

func (um UserModel) GetByEmail(email string) (*User, error) {
	query := `SELECT id, created_at, full_name, email, hashed_password, role_id, image FROM users WHERE email = $1`
	var user User

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := um.DB.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.CreatedAt, &user.FullName, &user.Email, &user.Password.hash, &user.RoleID, &user.Image)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}

func (um UserModel) Get(id int64) (*User, error) {
	query := `SELECT id, created_at, full_name, email, hashed_password, role_id, image FROM users WHERE id = $1`

	var user User

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := um.DB.QueryRowContext(ctx, query, id).Scan(
		&user.ID, &user.CreatedAt, &user.FullName, &user.Email, &user.Password.hash, &user.RoleID, &user.Image,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}
