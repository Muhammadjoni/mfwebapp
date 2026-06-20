package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/muhammadjoni/mfwebapp/internal/domain/entity"
	"github.com/muhammadjoni/mfwebapp/internal/domain/repository"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

const userCols = `id, email, password_hash, first_name, last_name, phone, role, status, language, avatar_url, created_at, updated_at`

func scanUser(row pgx.Row) (*entity.User, error) {
	u := &entity.User{}
	var id [16]byte
	var role, status, language string
	err := row.Scan(
		&id, &u.Email, &u.PasswordHash, &u.FirstName, &u.LastName,
		&u.Phone, &role, &status, &language, &u.AvatarURL,
		&u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	u.ID = uuid.UUID(id)
	u.Role = entity.Role(role)
	u.Status = entity.UserStatus(status)
	u.Language = entity.Language(language)
	return u, nil
}

func (r *UserRepository) Create(ctx context.Context, u *entity.User) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO users (id, email, password_hash, first_name, last_name, phone, role, status, language, avatar_url, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`,
		uuidBytes(u.ID), u.Email, u.PasswordHash, u.FirstName, u.LastName, u.Phone,
		string(u.Role), string(u.Status), string(u.Language), u.AvatarURL,
		u.CreatedAt, u.UpdatedAt,
	)
	return err
}

func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	row := r.pool.QueryRow(ctx, `SELECT `+userCols+` FROM users WHERE id = $1`, uuidBytes(id))
	u, err := scanUser(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}
	return u, err
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	row := r.pool.QueryRow(ctx, `SELECT `+userCols+` FROM users WHERE email = $1`, email)
	u, err := scanUser(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}
	return u, err
}

func (r *UserRepository) Update(ctx context.Context, u *entity.User) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE users SET email=$1, first_name=$2, last_name=$3, phone=$4, language=$5, avatar_url=$6, updated_at=NOW()
		WHERE id=$7`,
		u.Email, u.FirstName, u.LastName, u.Phone, string(u.Language), u.AvatarURL, uuidBytes(u.ID),
	)
	return err
}

func (r *UserRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status entity.UserStatus) error {
	_, err := r.pool.Exec(ctx, `UPDATE users SET status=$1, updated_at=NOW() WHERE id=$2`, string(status), uuidBytes(id))
	return err
}

func (r *UserRepository) UpdateRole(ctx context.Context, id uuid.UUID, role entity.Role) error {
	_, err := r.pool.Exec(ctx, `UPDATE users SET role=$1, updated_at=NOW() WHERE id=$2`, string(role), uuidBytes(id))
	return err
}

func (r *UserRepository) List(ctx context.Context, f repository.UserFilter) ([]*entity.User, int64, error) {
	where := []string{"1=1"}
	args := []any{}
	n := 1

	if f.Role != "" {
		where = append(where, fmt.Sprintf("role = $%d", n))
		args = append(args, string(f.Role))
		n++
	}
	if f.Status != "" {
		where = append(where, fmt.Sprintf("status = $%d", n))
		args = append(args, string(f.Status))
		n++
	}
	if f.Search != "" {
		where = append(where, fmt.Sprintf("(first_name ILIKE $%d OR last_name ILIKE $%d OR email ILIKE $%d)", n, n, n))
		args = append(args, "%"+f.Search+"%")
		n++
	}

	limit := f.Limit
	if limit <= 0 {
		limit = 20
	}
	page := f.Page
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit

	query := fmt.Sprintf(`
		SELECT `+userCols+`, COUNT(*) OVER() AS total
		FROM users WHERE %s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d`,
		strings.Join(where, " AND "), n, n+1,
	)
	args = append(args, limit, offset)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var users []*entity.User
	var total int64
	for rows.Next() {
		u := &entity.User{}
		var id [16]byte
		var role, status, language string
		if err := rows.Scan(
			&id, &u.Email, &u.PasswordHash, &u.FirstName, &u.LastName,
			&u.Phone, &role, &status, &language, &u.AvatarURL,
			&u.CreatedAt, &u.UpdatedAt, &total,
		); err != nil {
			return nil, 0, err
		}
		u.ID = uuid.UUID(id)
		u.Role = entity.Role(role)
		u.Status = entity.UserStatus(status)
		u.Language = entity.Language(language)
		users = append(users, u)
	}
	return users, total, rows.Err()
}

func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM users WHERE id=$1`, uuidBytes(id))
	return err
}

func (r *UserRepository) SaveRefreshToken(ctx context.Context, t *entity.RefreshToken) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO refresh_tokens (id, user_id, token, expires_at, created_at)
		VALUES ($1,$2,$3,$4,$5)`,
		uuidBytes(t.ID), uuidBytes(t.UserID), t.Token, t.ExpiresAt, t.CreatedAt,
	)
	return err
}

func (r *UserRepository) GetRefreshToken(ctx context.Context, token string) (*entity.RefreshToken, error) {
	rt := &entity.RefreshToken{}
	var id, userID [16]byte
	err := r.pool.QueryRow(ctx, `
		SELECT id, user_id, token, expires_at, created_at FROM refresh_tokens WHERE token=$1`, token).
		Scan(&id, &userID, &rt.Token, &rt.ExpiresAt, &rt.CreatedAt)
	if err != nil {
		return nil, err
	}
	rt.ID = uuid.UUID(id)
	rt.UserID = uuid.UUID(userID)
	return rt, nil
}

func (r *UserRepository) DeleteRefreshToken(ctx context.Context, token string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM refresh_tokens WHERE token=$1`, token)
	return err
}

func (r *UserRepository) DeleteUserRefreshTokens(ctx context.Context, userID uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM refresh_tokens WHERE user_id=$1`, uuidBytes(userID))
	return err
}
