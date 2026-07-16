package storage

import (
	"cloud/internal/models"
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/lib/pq"
	"github.com/lib/pq/pqerror"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrFailedToCreateUser        = errors.New("failed to create user")
	ErrFailedToGetUserByID       = errors.New("failed to get user by ID")
	ErrUserIDDoesntExist         = errors.New("user id doesn't exist")
	ErrUserEmailAlreadyExists    = errors.New("user email already exists")
	ErrPhoneNumberAlreadyExists  = errors.New("phone number already exists")
	ErrFailedToGetUserByEmail    = errors.New("failed to get user by email")
	ErrEmailDoesntExist          = errors.New("user email doesnt exist")
	ErrPasswordSameAsOldPassword = errors.New("new password is same as old password")
	ErrFailedToUpdateUser        = errors.New("failed to update user")
)

type UserSvc struct {
	db *sql.DB
}

func NewUserService(db *sql.DB) UserService {
	return &UserSvc{db: db}
}

func (u *UserSvc) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	query := `INSERT INTO users (name, email, password, phone) VALUES ($1, $2, $3, $4) RETURNING "ID", name, email, phone`
	stmt, err := u.db.PrepareContext(ctx, query)
	if err != nil {
		slog.Error("[CreateUser]: prepare statement", slog.Any("error", err))
		return nil, ErrFailedToCreateUser
	}
	defer stmt.Close()

	var newUser models.User
	if err := stmt.QueryRowContext(ctx, user.Name, user.Email, user.Password, user.Phone).Scan(&newUser.ID, &newUser.Name, &newUser.Email, &newUser.Phone); err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			switch pqErr.Constraint {
			case "users_email_key":
				slog.Error("[CreateUser]: email already exists", slog.Any("error", err))
				return nil, ErrUserEmailAlreadyExists
			case "users_phone_key":
				slog.Error("[CreateUser]: phone number already exists", slog.Any("error", err))
				return nil, ErrPhoneNumberAlreadyExists
			default:
				slog.Error("[CreateUser]: unique constraint violation", slog.Any("error", err))
				return nil, ErrFailedToCreateUser
			}
		}
	}

	return &newUser, nil
}

func (u *UserSvc) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	query := `SELECT "ID", name, email, password, phone FROM users WHERE "ID" = $1`

	stmt, err := u.db.PrepareContext(ctx, query)
	if err != nil {
		slog.Error("[GetUserByID]: prepare query ", slog.Any("error", err))
		return nil, ErrFailedToGetUserByID
	}
	defer stmt.Close()

	var user models.User
	if err := stmt.QueryRowContext(ctx, id).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Phone); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			slog.Error("[GetUserByID]: No user with id found", slog.Any("error", err))
			return nil, ErrUserIDDoesntExist
		}
		slog.Error("[GetUserByID]: query ", slog.Any("error", err))
		return nil, ErrFailedToGetUserByID
	}

	return &user, nil
}

func (u *UserSvc) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `SELECT "ID", name, email, password, phone FROM users WHERE "email" = $1`

	stmt, err := u.db.PrepareContext(ctx, query)
	if err != nil {
		slog.Error("[GetUserByEmail]: prepare query ", slog.Any("error", err))
		return nil, ErrFailedToGetUserByEmail
	}
	defer stmt.Close()

	var user models.User
	if err := stmt.QueryRowContext(ctx, email).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Phone); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			slog.Error("[GetUserByEmail]: No user with email found", slog.Any("error", err))
			return nil, ErrEmailDoesntExist
		}
		slog.Error("[GetUserByEmail]: query ", slog.Any("error", err))
		return nil, ErrFailedToGetUserByEmail
	}

	return &user, nil
}

func (u *UserSvc) UpdateUser(ctx context.Context, id string, req models.UpdateUserRequest) error {
	user, err := u.GetUserByID(ctx, id)
	if err != nil {
		return err
	}

	fields := make([]UpdateField, 0, 4)
	if req.Name != nil {
		fields = append(fields, UpdateField{Column: "name", Value: req.Name})
	}

	if req.Password != nil {
		if err := bcrypt.CompareHashAndPassword([]byte(*user.Password), []byte(*req.Password)); err == nil {
			slog.Error("[UpdateUser]: new password is same as old password", slog.Any("error", ErrPasswordSameAsOldPassword))
			return ErrPasswordSameAsOldPassword
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
		if err != nil {
			slog.Error("[UpdateUser]: generate password hash", slog.Any("error", err))
			return ErrFailedToUpdateUser
		}
		fields = append(fields, UpdateField{Column: "password", Value: string(hashedPassword)})
	}

	if req.Phone != nil {
		fields = append(fields, UpdateField{Column: "phone", Value: req.Phone})
	}

	if req.Email != nil {
		fields = append(fields, UpdateField{Column: "email", Value: req.Email})
	}

	query, args := BuildUpdateSQL("users", fields, []string{"ID"})
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
				slog.Error("[UpdateUser]: email already exists", slog.Any("error", err))
				return ErrUserEmailAlreadyExists
			case "users_phone_key":
				slog.Error("[UpdateUser]: phone number already exists", slog.Any("error", err))
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
