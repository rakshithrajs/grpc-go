package storage

import (
	"github.com/rakshithrajs/cloud/services/account/internal/models"
	"github.com/rakshithrajs/cloud/services/account/internal/utils"
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/lib/pq"
	"github.com/lib/pq/pqerror"
)

type userStore struct {
	db *sql.DB
}

func NewUserStore(db *sql.DB) UserService {
	return &userStore{db: db}
}

func (u *userStore) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	query := `INSERT INTO users (name, email, password, phone) VALUES ($1, $2, $3, $4) RETURNING "ID", name, email, phone`

	stmt, err := u.db.PrepareContext(ctx, query)
	if err != nil {
		slog.Error("[CreateUser]: prepare statement", slog.Any("error", err))
		return nil, ErrFailedToCreateUser
	}
	defer stmt.Close()

	var newUser models.User
	if err := stmt.QueryRowContext(ctx, *user.Name, *user.Email, *user.Password, *user.Phone).Scan(&newUser.ID, &newUser.Name, &newUser.Email, &newUser.Phone); err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == pqerror.UniqueViolation {
			switch pqErr.Constraint {
			case "users_email_key":
				return nil, ErrUserEmailAlreadyExists
			case "users_phone_key":
				return nil, ErrPhoneNumberAlreadyExists
			default:
				slog.Error("[CreateUser]: unique constraint violation", slog.Any("error", err))
				return nil, ErrFailedToCreateUser
			}
		}
		slog.Error("[CreateUser]: query row", slog.Any("error", err))
		return nil, ErrFailedToCreateUser
	}

	return &newUser, nil
}

func (u *userStore) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	query := `SELECT "ID", name, email, password, phone FROM users WHERE "ID" = $1`

	stmt, err := u.db.PrepareContext(ctx, query)
	if err != nil {
		slog.Error("[GetUserByID]: prepare query", slog.Any("error", err))
		return nil, ErrFailedToGetUserByID
	}
	defer stmt.Close()

	var user models.User
	if err := stmt.QueryRowContext(ctx, id).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Phone); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		slog.Error("[GetUserByID]: query", slog.Any("error", err))
		return nil, ErrFailedToGetUserByID
	}

	return &user, nil
}

func (u *userStore) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `SELECT "ID", name, email, password, phone FROM users WHERE email = $1`

	stmt, err := u.db.PrepareContext(ctx, query)
	if err != nil {
		slog.Error("[GetUserByEmail]: prepare query", slog.Any("error", err))
		return nil, ErrFailedToGetUserByEmail
	}
	defer stmt.Close()

	var user models.User
	if err := stmt.QueryRowContext(ctx, email).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Phone); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrEmailNotFound
		}
		slog.Error("[GetUserByEmail]: query", slog.Any("error", err))
		return nil, ErrFailedToGetUserByEmail
	}

	return &user, nil
}

func (u *userStore) UpdateUser(ctx context.Context, id string, req models.UpdateUserRequest) error {
	fields := make([]utils.UpdateField, 0, 4)
	if req.Name != nil {
		fields = append(fields, utils.UpdateField{Column: "name", Value: *req.Name})
	}
	if req.Password != nil {
		fields = append(fields, utils.UpdateField{Column: "password", Value: *req.Password})
	}
	if req.Phone != nil {
		fields = append(fields, utils.UpdateField{Column: "phone", Value: *req.Phone})
	}
	if req.Email != nil {
		fields = append(fields, utils.UpdateField{Column: "email", Value: *req.Email})
	}

	query, args := utils.BuildUpdateSQL("users", fields, []string{"ID"})
	args[0] = id

	stmt, err := u.db.PrepareContext(ctx, query)
	if err != nil {
		slog.Error("[UpdateUser]: prepare statement", slog.Any("error", err))
		return ErrFailedToUpdateUser
	}
	defer stmt.Close()

	if _, err := stmt.ExecContext(ctx, args...); err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == pqerror.UniqueViolation {
			switch pqErr.Constraint {
			case "users_email_key":
				return ErrUserEmailAlreadyExists
			case "users_phone_key":
				return ErrPhoneNumberAlreadyExists
			default:
				slog.Error("[UpdateUser]: unique constraint violation", slog.Any("error", err))
				return ErrFailedToUpdateUser
			}
		}
		slog.Error("[UpdateUser]: update user", slog.Any("error", err))
		return ErrFailedToUpdateUser
	}

	return nil
}
