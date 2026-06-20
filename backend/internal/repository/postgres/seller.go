package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/muhammadjoni/mfwebapp/internal/domain/entity"
	"github.com/muhammadjoni/mfwebapp/internal/domain/repository"
)

type SellerRepository struct{ pool *pgxpool.Pool }

func NewSellerRepository(pool *pgxpool.Pool) *SellerRepository {
	return &SellerRepository{pool: pool}
}

const sellerCols = `
	id, user_id, business_name, business_email, business_phone, country,
	status, commission_rate, description, logo_url, verified_at, created_at, updated_at`

func collectSeller(row pgx.CollectableRow) (*entity.Seller, error) {
	s := &entity.Seller{}
	var id, userID [16]byte
	var status string
	var verifiedAt pgtype.Timestamptz
	if err := row.Scan(
		&id, &userID, &s.BusinessName, &s.BusinessEmail, &s.BusinessPhone, &s.Country,
		&status, &s.CommissionRate, &s.Description, &s.LogoURL,
		&verifiedAt, &s.CreatedAt, &s.UpdatedAt,
	); err != nil {
		return nil, err
	}
	s.ID = uuid.UUID(id)
	s.UserID = uuid.UUID(userID)
	s.Status = entity.SellerStatus(status)
	s.VerifiedAt = pgTimestamptzPtr(verifiedAt)
	return s, nil
}

func (r *SellerRepository) Create(ctx context.Context, s *entity.Seller) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO sellers (id, user_id, business_name, business_email, business_phone, country,
		                     status, commission_rate, description, logo_url, verified_at, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)`,
		uuidBytes(s.ID), uuidBytes(s.UserID), s.BusinessName, s.BusinessEmail, s.BusinessPhone, s.Country,
		string(s.Status), s.CommissionRate, s.Description, s.LogoURL,
		nullableTime(s.VerifiedAt), s.CreatedAt, s.UpdatedAt,
	)
	return err
}

func (r *SellerRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Seller, error) {
	rows, err := r.pool.Query(ctx, `SELECT `+sellerCols+` FROM sellers WHERE id=$1`, uuidBytes(id))
	if err != nil {
		return nil, err
	}
	return pgx.CollectOneRow(rows, collectSeller)
}

func (r *SellerRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*entity.Seller, error) {
	rows, err := r.pool.Query(ctx, `SELECT `+sellerCols+` FROM sellers WHERE user_id=$1`, uuidBytes(userID))
	if err != nil {
		return nil, err
	}
	return pgx.CollectOneRow(rows, collectSeller)
}

func (r *SellerRepository) Update(ctx context.Context, s *entity.Seller) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE sellers SET business_name=$1, business_email=$2, business_phone=$3, country=$4,
		                   commission_rate=$5, description=$6, logo_url=$7, verified_at=$8, updated_at=NOW()
		WHERE id=$9`,
		s.BusinessName, s.BusinessEmail, s.BusinessPhone, s.Country,
		s.CommissionRate, s.Description, s.LogoURL, nullableTime(s.VerifiedAt), uuidBytes(s.ID),
	)
	return err
}

func (r *SellerRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status entity.SellerStatus) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE sellers SET status=$1, updated_at=NOW() WHERE id=$2`, string(status), uuidBytes(id))
	return err
}

func (r *SellerRepository) List(ctx context.Context, f repository.SellerFilter) ([]*entity.Seller, int64, error) {
	where := []string{"1=1"}
	args := []any{}
	n := 1

	if f.Status != "" {
		where = append(where, fmt.Sprintf("status = $%d", n))
		args = append(args, string(f.Status))
		n++
	}
	if f.Search != "" {
		where = append(where, fmt.Sprintf("(business_name ILIKE $%d OR business_email ILIKE $%d)", n, n))
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

	query := fmt.Sprintf(`
		SELECT `+sellerCols+`, COUNT(*) OVER() AS total
		FROM sellers WHERE %s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d`,
		strings.Join(where, " AND "), n, n+1,
	)
	args = append(args, limit, (page-1)*limit)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var sellers []*entity.Seller
	var total int64
	for rows.Next() {
		s := &entity.Seller{}
		var id, userID [16]byte
		var status string
		var verifiedAt pgtype.Timestamptz
		var rowTotal int64
		if err := rows.Scan(
			&id, &userID, &s.BusinessName, &s.BusinessEmail, &s.BusinessPhone, &s.Country,
			&status, &s.CommissionRate, &s.Description, &s.LogoURL,
			&verifiedAt, &s.CreatedAt, &s.UpdatedAt, &rowTotal,
		); err != nil {
			return nil, 0, err
		}
		s.ID = uuid.UUID(id)
		s.UserID = uuid.UUID(userID)
		s.Status = entity.SellerStatus(status)
		s.VerifiedAt = pgTimestamptzPtr(verifiedAt)
		total = rowTotal
		sellers = append(sellers, s)
	}
	if sellers == nil {
		sellers = []*entity.Seller{}
	}
	return sellers, total, rows.Err()
}
