package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/muhammadjoni/mfwebapp/internal/domain/entity"
	"github.com/muhammadjoni/mfwebapp/internal/domain/repository"
)

// ── ProductRepository ────────────────────────────────────────────────────────

type ProductRepository struct{ pool *pgxpool.Pool }

func NewProductRepository(pool *pgxpool.Pool) *ProductRepository {
	return &ProductRepository{pool: pool}
}

// salePriceParam returns nil (→ SQL NULL) or the float64 value.
// pgx passes float64 as FLOAT8; PostgreSQL implicitly casts FLOAT8→NUMERIC.
func salePriceParam(v *float64) interface{} {
	if v == nil {
		return nil
	}
	return *v
}

const productSelectCols = `
	p.id, p.seller_id, p.category_id, p.name, p.slug, p.description, p.short_desc,
	p.base_price, p.sale_price, p.currency, p.sku, p.stock, p.status,
	p.images, p.tags, p.specifications, p.weight, p.dimensions,
	p.view_count, p.sold_count, p.rating, p.review_count, p.featured_at,
	p.created_at, p.updated_at`

func collectProduct(row pgx.CollectableRow) (*entity.Product, error) {
	p := &entity.Product{}
	var (
		id, sellerID, catID [16]byte
		salePrice           pgtype.Numeric
		images, tags        pgtype.Array[pgtype.Text]
		specsRaw, dimsRaw   json.RawMessage
		featuredAt          pgtype.Timestamptz
		status              string
	)
	if err := row.Scan(
		&id, &sellerID, &catID, &p.Name, &p.Slug, &p.Description, &p.ShortDesc,
		&p.BasePrice, &salePrice, &p.Currency, &p.SKU, &p.Stock, &status,
		&images, &tags, &specsRaw, &p.Weight, &dimsRaw,
		&p.ViewCount, &p.SoldCount, &p.Rating, &p.ReviewCount, &featuredAt,
		&p.CreatedAt, &p.UpdatedAt,
	); err != nil {
		return nil, err
	}
	p.ID = uuid.UUID(id)
	p.SellerID = uuid.UUID(sellerID)
	p.CategoryID = uuid.UUID(catID)
	p.Status = entity.ProductStatus(status)
	p.SalePrice = numericPtr(salePrice)
	p.Images = fromTextArray(images)
	p.Tags = fromTextArray(tags)
	p.FeaturedAt = pgTimestamptzPtr(featuredAt)
	p.Specifications = map[string]string{}
	if specsRaw != nil {
		_ = json.Unmarshal(specsRaw, &p.Specifications)
	}
	if dimsRaw != nil {
		_ = json.Unmarshal(dimsRaw, &p.Dimensions)
	}
	return p, nil
}

func (r *ProductRepository) Create(ctx context.Context, p *entity.Product) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO products (
			id, seller_id, category_id, name, slug, description, short_desc,
			base_price, sale_price, currency, sku, stock, status,
			images, tags, specifications, weight, dimensions, created_at, updated_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20)`,
		uuidBytes(p.ID), uuidBytes(p.SellerID), uuidBytes(p.CategoryID),
		p.Name, p.Slug, p.Description, p.ShortDesc,
		p.BasePrice, salePriceParam(p.SalePrice),
		p.Currency, p.SKU, p.Stock, string(p.Status),
		toTextArray(p.Images), toTextArray(p.Tags),
		marshalJSON(p.Specifications), p.Weight, marshalJSON(p.Dimensions),
		p.CreatedAt, p.UpdatedAt,
	)
	return err
}

func (r *ProductRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Product, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT `+productSelectCols+` FROM products p WHERE p.id = $1`, uuidBytes(id))
	if err != nil {
		return nil, err
	}
	p, err := pgx.CollectOneRow(rows, collectProduct)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}
	return p, err
}

func (r *ProductRepository) GetBySlug(ctx context.Context, slug string) (*entity.Product, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT `+productSelectCols+` FROM products p WHERE p.slug = $1`, slug)
	if err != nil {
		return nil, err
	}
	return pgx.CollectOneRow(rows, collectProduct)
}

func (r *ProductRepository) Update(ctx context.Context, p *entity.Product) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE products SET
			name=$1, description=$2, short_desc=$3, base_price=$4, sale_price=$5,
			currency=$6, sku=$7, stock=$8, images=$9, tags=$10,
			specifications=$11, weight=$12, dimensions=$13, updated_at=NOW()
		WHERE id=$14`,
		p.Name, p.Description, p.ShortDesc, p.BasePrice, salePriceParam(p.SalePrice),
		p.Currency, p.SKU, p.Stock,
		toTextArray(p.Images), toTextArray(p.Tags),
		marshalJSON(p.Specifications), p.Weight, marshalJSON(p.Dimensions),
		uuidBytes(p.ID),
	)
	return err
}

func (r *ProductRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status entity.ProductStatus) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE products SET status=$1, updated_at=NOW() WHERE id=$2`, string(status), uuidBytes(id))
	return err
}

func (r *ProductRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM products WHERE id=$1`, uuidBytes(id))
	return err
}

func (r *ProductRepository) List(ctx context.Context, f repository.ProductFilter) ([]*entity.Product, int64, error) {
	where := []string{"1=1"}
	args := []any{}
	n := 1

	if f.Status != "" {
		where = append(where, fmt.Sprintf("p.status = $%d", n))
		args = append(args, string(f.Status))
		n++
	}
	if f.CategoryID != nil {
		where = append(where, fmt.Sprintf("p.category_id = $%d", n))
		args = append(args, uuidBytes(*f.CategoryID))
		n++
	}
	if f.SellerID != nil {
		where = append(where, fmt.Sprintf("p.seller_id = $%d", n))
		args = append(args, uuidBytes(*f.SellerID))
		n++
	}
	if f.Search != "" {
		where = append(where, fmt.Sprintf("p.name ILIKE $%d", n))
		args = append(args, "%"+f.Search+"%")
		n++
	}
	if f.MinPrice > 0 {
		where = append(where, fmt.Sprintf("p.base_price >= $%d", n))
		args = append(args, f.MinPrice)
		n++
	}
	if f.MaxPrice > 0 {
		where = append(where, fmt.Sprintf("p.base_price <= $%d", n))
		args = append(args, f.MaxPrice)
		n++
	}

	allowedSort := map[string]string{
		"created_at": "p.created_at", "base_price": "p.base_price",
		"rating": "p.rating", "view_count": "p.view_count",
		"sold_count": "p.sold_count", "name": "p.name",
	}
	sortCol := "p.created_at"
	if col, ok := allowedSort[f.SortBy]; ok {
		sortCol = col
	}
	sortDir := "DESC"
	if strings.ToUpper(f.SortDir) == "ASC" {
		sortDir = "ASC"
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
		SELECT `+productSelectCols+`, COUNT(*) OVER() AS total
		FROM products p WHERE %s
		ORDER BY %s %s
		LIMIT $%d OFFSET $%d`,
		strings.Join(where, " AND "), sortCol, sortDir, n, n+1,
	)
	args = append(args, limit, (page-1)*limit)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var products []*entity.Product
	var total int64
	for rows.Next() {
		p := &entity.Product{}
		var (
			id, sellerID, catID [16]byte
			salePrice           pgtype.Numeric
			images, tags        pgtype.Array[pgtype.Text]
			specsRaw, dimsRaw   json.RawMessage
			featuredAt          pgtype.Timestamptz
			status              string
			rowTotal            int64
		)
		if err := rows.Scan(
			&id, &sellerID, &catID, &p.Name, &p.Slug, &p.Description, &p.ShortDesc,
			&p.BasePrice, &salePrice, &p.Currency, &p.SKU, &p.Stock, &status,
			&images, &tags, &specsRaw, &p.Weight, &dimsRaw,
			&p.ViewCount, &p.SoldCount, &p.Rating, &p.ReviewCount, &featuredAt,
			&p.CreatedAt, &p.UpdatedAt, &rowTotal,
		); err != nil {
			return nil, 0, err
		}
		p.ID = uuid.UUID(id)
		p.SellerID = uuid.UUID(sellerID)
		p.CategoryID = uuid.UUID(catID)
		p.Status = entity.ProductStatus(status)
		p.SalePrice = numericPtr(salePrice)
		p.Images = fromTextArray(images)
		p.Tags = fromTextArray(tags)
		p.FeaturedAt = pgTimestamptzPtr(featuredAt)
		p.Specifications = map[string]string{}
		if specsRaw != nil {
			_ = json.Unmarshal(specsRaw, &p.Specifications)
		}
		if dimsRaw != nil {
			_ = json.Unmarshal(dimsRaw, &p.Dimensions)
		}
		total = rowTotal
		products = append(products, p)
	}
	if products == nil {
		products = []*entity.Product{}
	}
	return products, total, rows.Err()
}

func (r *ProductRepository) IncrementViewCount(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `UPDATE products SET view_count=view_count+1 WHERE id=$1`, uuidBytes(id))
	return err
}

func (r *ProductRepository) DecrementStock(ctx context.Context, id uuid.UUID, qty int) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE products SET stock=stock-$1 WHERE id=$2 AND stock >= $1`, qty, uuidBytes(id))
	return err
}

func (r *ProductRepository) IncrementSoldCount(ctx context.Context, id uuid.UUID, qty int) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE products SET sold_count=sold_count+$1 WHERE id=$2`, qty, uuidBytes(id))
	return err
}

// ── CategoryRepository ───────────────────────────────────────────────────────

type CategoryRepository struct{ pool *pgxpool.Pool }

func NewCategoryRepository(pool *pgxpool.Pool) *CategoryRepository {
	return &CategoryRepository{pool: pool}
}

const catCols = `id, parent_id, name, slug, description, image_url, sort_order, is_active, created_at, updated_at`

func collectCategory(row pgx.CollectableRow) (*entity.Category, error) {
	c := &entity.Category{}
	var id [16]byte
	var parentID pgtype.UUID
	if err := row.Scan(&id, &parentID, &c.Name, &c.Slug, &c.Description, &c.ImageURL,
		&c.SortOrder, &c.IsActive, &c.CreatedAt, &c.UpdatedAt); err != nil {
		return nil, err
	}
	c.ID = uuid.UUID(id)
	c.ParentID = pgUUIDPtr(parentID)
	return c, nil
}

func (r *CategoryRepository) Create(ctx context.Context, c *entity.Category) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO categories (id, parent_id, name, slug, description, image_url, sort_order, is_active, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`,
		uuidBytes(c.ID), nullableUUIDParam(c.ParentID), c.Name, c.Slug,
		c.Description, c.ImageURL, c.SortOrder, c.IsActive, c.CreatedAt, c.UpdatedAt,
	)
	return err
}

func (r *CategoryRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Category, error) {
	rows, err := r.pool.Query(ctx, `SELECT `+catCols+` FROM categories WHERE id=$1`, uuidBytes(id))
	if err != nil {
		return nil, err
	}
	return pgx.CollectOneRow(rows, collectCategory)
}

func (r *CategoryRepository) GetBySlug(ctx context.Context, slug string) (*entity.Category, error) {
	rows, err := r.pool.Query(ctx, `SELECT `+catCols+` FROM categories WHERE slug=$1`, slug)
	if err != nil {
		return nil, err
	}
	return pgx.CollectOneRow(rows, collectCategory)
}

func (r *CategoryRepository) Update(ctx context.Context, c *entity.Category) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE categories SET name=$1, slug=$2, description=$3, image_url=$4, sort_order=$5, is_active=$6, updated_at=NOW()
		WHERE id=$7`,
		c.Name, c.Slug, c.Description, c.ImageURL, c.SortOrder, c.IsActive, uuidBytes(c.ID),
	)
	return err
}

func (r *CategoryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM categories WHERE id=$1`, uuidBytes(id))
	return err
}

func (r *CategoryRepository) List(ctx context.Context) ([]*entity.Category, error) {
	rows, err := r.pool.Query(ctx, `SELECT `+catCols+` FROM categories ORDER BY sort_order, name`)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, collectCategory)
}
