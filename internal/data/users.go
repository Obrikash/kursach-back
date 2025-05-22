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
	ErrDuplicateEmail         = errors.New("duplicate email")
	ErrTrainerAlreadyAttached = errors.New("Trainer is already attached to the pool.")
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

type PoolWithTrainers struct {
	Pool     Pool    `json:"pool"`
	Trainers []*User `json:"trainers"`
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
	query := `SELECT t.user_id, full_name, u.image, u.email FROM users u JOIN trainers t ON t.user_id = u.id`

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

		err := rows.Scan(&trainer.ID, &trainer.FullName, &trainer.Image, &trainer.Email)
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

func (um UserModel) GetTrainersForPools() ([]PoolWithTrainers, error) {
	query := `
        SELECT p.id, p.name, p.address, p.type,
               u.id, u.full_name, u.email, u.image
        FROM trainers t
        JOIN pools p ON t.pool_id = p.id
        JOIN users u ON t.user_id = u.id
        ORDER BY p.name, u.id
    `

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := um.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []PoolWithTrainers
	var currentPool Pool
	var trainers []*User

	for rows.Next() {
		var pool Pool
		var trainer User

		err := rows.Scan(
			&pool.ID, &pool.Name, &pool.Address, &pool.PoolType,
			&trainer.ID, &trainer.FullName, &trainer.Email, &trainer.Image,
		)
		if err != nil {
			return nil, err
		}

		if currentPool.ID == 0 || pool.ID != currentPool.ID {
			// New pool detected
			if currentPool.ID != 0 {
				result = append(result, PoolWithTrainers{
					Pool:     currentPool,
					Trainers: trainers,
				})
			}
			currentPool = pool
			trainers = nil // Reset trainers for the new pool
		}

		trainers = append(trainers, &trainer)
	}

	// Add the last pool
	if currentPool.ID != 0 {
		result = append(result, PoolWithTrainers{
			Pool:     currentPool,
			Trainers: trainers,
		})
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
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

type ProfitTrainersPools struct {
	ID       int64   `json:"id"`
	FullName string  `json:"full_name"`
	PoolId   int64   `json:"pool_id"`
	PoolName string  `json:"pool_name"`
	Profit   float64 `json:"profit"`
}

func (um UserModel) ProfitForEachTrainerInEachPool() ([]*ProfitTrainersPools, error) {
	query := `SELECT tr.id AS trainer_id, u_trainer.full_name AS trainer_name, p.id AS pool_id, p.name AS pool_name,
	SUM(sub.price) AS total_profit FROM user_subscriptions us JOIN user_groups ug ON us.user_id = ug.user_id 
	JOIN training_groups tg ON ug.group_id = tg.id JOIN trainers tr ON tg.trainer_id = tr.id JOIN users u_trainer ON 
	tr.user_id = u_trainer.id JOIN pools p ON tr.pool_id = p.id JOIN subscriptions sub ON us.subscription_id = sub.id 
	GROUP BY tr.id, u_trainer.full_name, p.id, p.name ORDER BY p.name, u_trainer.full_name;`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := um.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	profitsOfTrainers := []*ProfitTrainersPools{}

	for rows.Next() {
		var trainerProfit ProfitTrainersPools

		err := rows.Scan(&trainerProfit.ID, &trainerProfit.FullName, &trainerProfit.PoolId, &trainerProfit.PoolName, &trainerProfit.Profit)
		if err != nil {
			return nil, err
		}

		profitsOfTrainers = append(profitsOfTrainers, &trainerProfit)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return profitsOfTrainers, nil
}

func (um UserModel) AttachTrainerToPool(userID int64, poolID int64) error {
	query := "INSERT INTO trainers (user_id, pool_id) VALUES ($1, $2) RETURNING id"

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := um.DB.QueryRowContext(ctx, query, userID, poolID).Scan(&userID)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "trainers_user_id_key"`:
			return ErrTrainerAlreadyAttached
		default:
			return err
		}
	}

	return nil
}
