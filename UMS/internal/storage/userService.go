package storage

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/rakshithrajs/cloud/UMS/internal/config"
	handlerErrors "github.com/rakshithrajs/cloud/UMS/internal/handlers/errors"
	"github.com/rakshithrajs/cloud/UMS/internal/models"
	storageUtils "github.com/rakshithrajs/cloud/UMS/internal/storage/utils"

	"github.com/lib/pq"
	"github.com/lib/pq/pqerror"
)

const (
	// function name for Create User
	fnCreateUser = "CreateUser"

	// function name for Get User by ID
	fnGetUserByID = "GetUserByID"

	// function name for Get User by Email
	fnGetUserByEmail = "GetUserByEmail"

	// function name for Update User
	fnUpdateUser = "UpdateUser"

	// unique constraint names for users email key
	uniqueConstraintUsersEmailKey = "users_email_key"

	// unique constraint names for users phone key
	uniqueConstraintUsersPhoneKey = "users_phone_key"
)

func logPrefix(fn string) string { return "[" + fn + "]: " }

type userStore struct {
	db *sql.DB
}

func NewUserStore(db *sql.DB) UserService {
	return &userStore{db: db}
}

func (u *userStore) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	query := `INSERT INTO users (name, email, password, phone) VALUES ($1, $2, $3, $4) RETURNING "ID", name, email, phone, "createdAtUTC", "updatedAtUTC"`

	stmt, err := u.db.PrepareContext(ctx, query)
	if err != nil {
		slog.Error(logPrefix(fnCreateUser)+"prepare statement", slog.Any(config.ErrorKey, err))
		return nil, ErrFailedToCreateUser
	}
	defer stmt.Close()

	var newUser models.User
	if err := stmt.QueryRowContext(ctx, user.Name, user.Email, user.Password, user.Phone).Scan(&newUser.ID, &newUser.Name, &newUser.Email, &newUser.Phone, &newUser.CreatedAtUTC, &newUser.UpdatedAtUTC); err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == pqerror.UniqueViolation {
			switch pqErr.Constraint {
			case uniqueConstraintUsersEmailKey:
				return nil, handlerErrors.ErrUserEmailAlreadyExists
			case uniqueConstraintUsersPhoneKey:
				return nil, handlerErrors.ErrPhoneNumberAlreadyExists
			default:
				slog.Error(logPrefix(fnCreateUser)+"unique constraint violation", slog.Any(config.ErrorKey, err))
				return nil, ErrFailedToCreateUser
			}
		}
		slog.Error(logPrefix(fnCreateUser)+"query row", slog.Any(config.ErrorKey, err))
		return nil, ErrFailedToCreateUser
	}

	return &newUser, nil
}

func (u *userStore) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	query := `SELECT "ID", name, email, password, phone, "createdAtUTC", "updatedAtUTC" FROM users WHERE "ID" = $1`

	stmt, err := u.db.PrepareContext(ctx, query)
	if err != nil {
		slog.Error(logPrefix(fnGetUserByID)+"prepare query", slog.Any(config.ErrorKey, err))
		return nil, ErrFailedToGetUserByID
	}
	defer stmt.Close()

	var user models.User
	if err := stmt.QueryRowContext(ctx, id).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Phone, &user.CreatedAtUTC, &user.UpdatedAtUTC); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, handlerErrors.ErrUserNotFound
		}
		slog.Error(logPrefix(fnGetUserByID)+"query", slog.Any(config.ErrorKey, err))
		return nil, ErrFailedToGetUserByID
	}

	return &user, nil
}

func (u *userStore) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `SELECT "ID", name, email, password, phone FROM users WHERE email = $1`

	stmt, err := u.db.PrepareContext(ctx, query)
	if err != nil {
		slog.Error(logPrefix(fnGetUserByEmail)+"prepare query", slog.Any(config.ErrorKey, err))
		return nil, ErrFailedToGetUserByEmail
	}
	defer stmt.Close()

	var user models.User
	if err := stmt.QueryRowContext(ctx, email).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Phone); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, handlerErrors.ErrEmailNotFound
		}
		slog.Error(logPrefix(fnGetUserByEmail)+"query", slog.Any(config.ErrorKey, err))
		return nil, ErrFailedToGetUserByEmail
	}

	return &user, nil
}

func (u *userStore) UpdateUser(ctx context.Context, id string, req models.UpdateUserRequest) error {
	fields := make([]storageUtils.UpdateField, 0, 4)
	if req.Name != config.NullString {
		fields = append(fields, storageUtils.UpdateField{Column: "name", Value: req.Name})
	}
	if req.Password != config.NullString {
		fields = append(fields, storageUtils.UpdateField{Column: "password", Value: req.Password})
	}
	if req.Phone != config.NullString {
		fields = append(fields, storageUtils.UpdateField{Column: "phone", Value: req.Phone})
	}
	if req.Email != config.NullString {
		fields = append(fields, storageUtils.UpdateField{Column: "email", Value: req.Email})
	}

	query, args := storageUtils.BuildUpdateSQL("users", fields, []string{"ID"})
	args[0] = id

	stmt, err := u.db.PrepareContext(ctx, query)
	if err != nil {
		slog.Error(logPrefix(fnUpdateUser)+"prepare statement", slog.Any(config.ErrorKey, err))
		return ErrFailedToUpdateUser
	}
	defer stmt.Close()

	if _, err := stmt.ExecContext(ctx, args...); err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == pqerror.UniqueViolation {
			switch pqErr.Constraint {
			case uniqueConstraintUsersEmailKey:
				return handlerErrors.ErrUserEmailAlreadyExists
			case uniqueConstraintUsersPhoneKey:
				return handlerErrors.ErrPhoneNumberAlreadyExists
			default:
				slog.Error(logPrefix(fnUpdateUser)+"unique constraint violation", slog.Any(config.ErrorKey, err))
				return ErrFailedToUpdateUser
			}
		}
		slog.Error(logPrefix(fnUpdateUser)+"update user", slog.Any(config.ErrorKey, err))
		return ErrFailedToUpdateUser
	}

	return nil
}
