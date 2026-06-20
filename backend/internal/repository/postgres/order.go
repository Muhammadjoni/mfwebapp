package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/muhammadjoni/mfwebapp/internal/domain/entity"
	"github.com/muhammadjoni/mfwebapp/internal/domain/repository"
)

// ── OrderRepository ──────────────────────────────────────────────────────────

type OrderRepository struct{ pool *pgxpool.Pool }

func NewOrderRepository(pool *pgxpool.Pool) *OrderRepository {
	return &OrderRepository{pool: pool}
}

func (r *OrderRepository) Create(ctx context.Context, o *entity.Order) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx) //nolint:errcheck

	addrJSON := marshalJSON(o.ShippingAddress)

	_, err = tx.Exec(ctx, `
		INSERT INTO orders (id, user_id, status, subtotal, shipping_cost, tax, total, currency,
		                    shipping_address, notes, tracking_number, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)`,
		uuidBytes(o.ID), uuidBytes(o.UserID), string(o.Status),
		o.SubTotal, o.ShippingCost, o.Tax, o.Total, o.Currency,
		addrJSON, o.Notes, o.TrackingNumber,
		o.CreatedAt, o.UpdatedAt,
	)
	if err != nil {
		return err
	}

	batch := &pgx.Batch{}
	for _, item := range o.Items {
		batch.Queue(`
			INSERT INTO order_items (id, order_id, product_id, variant_id, seller_id,
			                        name, sku, image_url, quantity, unit_price, total_price)
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`,
			uuidBytes(item.ID), uuidBytes(o.ID), uuidBytes(item.ProductID),
			nullableUUIDParam(item.VariantID), uuidBytes(item.SellerID),
			item.Name, item.SKU, item.ImageURL,
			item.Quantity, item.UnitPrice, item.TotalPrice,
		)
	}
	if len(o.Items) > 0 {
		br := tx.SendBatch(ctx, batch)
		for range o.Items {
			if _, err := br.Exec(); err != nil {
				br.Close()
				return err
			}
		}
		br.Close()
	}

	return tx.Commit(ctx)
}

func (r *OrderRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Order, error) {
	o := &entity.Order{}
	var (
		oid, uid  [16]byte
		addrRaw   json.RawMessage
		paymentID pgtype.UUID
		status    string
	)
	err := r.pool.QueryRow(ctx, `
		SELECT id, user_id, status, subtotal, shipping_cost, tax, total, currency,
		       shipping_address, notes, tracking_number, payment_id, created_at, updated_at
		FROM orders WHERE id=$1`, uuidBytes(id)).Scan(
		&oid, &uid, &status, &o.SubTotal, &o.ShippingCost, &o.Tax, &o.Total, &o.Currency,
		&addrRaw, &o.Notes, &o.TrackingNumber, &paymentID, &o.CreatedAt, &o.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	o.ID = uuid.UUID(oid)
	o.UserID = uuid.UUID(uid)
	o.Status = entity.OrderStatus(status)
	o.PaymentID = pgUUIDPtr(paymentID)
	if addrRaw != nil {
		_ = json.Unmarshal(addrRaw, &o.ShippingAddress)
	}

	items, err := r.fetchItems(ctx, o.ID)
	if err != nil {
		return nil, err
	}
	o.Items = items
	return o, nil
}

func (r *OrderRepository) fetchItems(ctx context.Context, orderID uuid.UUID) ([]entity.OrderItem, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, order_id, product_id, variant_id, seller_id,
		       name, sku, image_url, quantity, unit_price, total_price
		FROM order_items WHERE order_id=$1`, uuidBytes(orderID))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []entity.OrderItem
	for rows.Next() {
		var item entity.OrderItem
		var id, oid, pid, sid [16]byte
		var vid pgtype.UUID
		if err := rows.Scan(&id, &oid, &pid, &vid, &sid,
			&item.Name, &item.SKU, &item.ImageURL,
			&item.Quantity, &item.UnitPrice, &item.TotalPrice); err != nil {
			return nil, err
		}
		item.ID = uuid.UUID(id)
		item.OrderID = uuid.UUID(oid)
		item.ProductID = uuid.UUID(pid)
		item.SellerID = uuid.UUID(sid)
		item.VariantID = pgUUIDPtr(vid)
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *OrderRepository) Update(ctx context.Context, o *entity.Order) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE orders SET status=$1, tracking_number=$2, payment_id=$3, updated_at=NOW()
		WHERE id=$4`,
		string(o.Status), o.TrackingNumber, nullableUUIDParam(o.PaymentID), uuidBytes(o.ID),
	)
	return err
}

func (r *OrderRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status entity.OrderStatus, changedBy uuid.UUID, note string) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx) //nolint:errcheck

	if _, err = tx.Exec(ctx,
		`UPDATE orders SET status=$1, updated_at=NOW() WHERE id=$2`, string(status), uuidBytes(id)); err != nil {
		return err
	}
	if _, err = tx.Exec(ctx, `
		INSERT INTO order_status_history (id, order_id, status, changed_by, note, created_at)
		VALUES ($1,$2,$3,$4,$5,$6)`,
		uuidBytes(uuid.New()), uuidBytes(id), string(status), uuidBytes(changedBy), note, time.Now(),
	); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (r *OrderRepository) List(ctx context.Context, f repository.OrderFilter) ([]*entity.Order, int64, error) {
	where := []string{"1=1"}
	args := []any{}
	n := 1

	if f.UserID != nil {
		where = append(where, fmt.Sprintf("user_id = $%d", n))
		args = append(args, uuidBytes(*f.UserID))
		n++
	}
	if f.Status != "" {
		where = append(where, fmt.Sprintf("status = $%d", n))
		args = append(args, string(f.Status))
		n++
	}

	limit := 20
	page := 1
	if f.Limit > 0 {
		limit = f.Limit
	}
	if f.Page > 0 {
		page = f.Page
	}

	query := fmt.Sprintf(`
		SELECT id, user_id, status, subtotal, shipping_cost, tax, total, currency,
		       shipping_address, notes, tracking_number, payment_id, created_at, updated_at,
		       COUNT(*) OVER() AS total
		FROM orders WHERE %s
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

	var orders []*entity.Order
	var total int64
	for rows.Next() {
		o := &entity.Order{}
		var (
			oid, uid  [16]byte
			addrRaw   json.RawMessage
			paymentID pgtype.UUID
			status    string
			rowTotal  int64
		)
		if err := rows.Scan(
			&oid, &uid, &status, &o.SubTotal, &o.ShippingCost, &o.Tax, &o.Total, &o.Currency,
			&addrRaw, &o.Notes, &o.TrackingNumber, &paymentID, &o.CreatedAt, &o.UpdatedAt, &rowTotal,
		); err != nil {
			return nil, 0, err
		}
		o.ID = uuid.UUID(oid)
		o.UserID = uuid.UUID(uid)
		o.Status = entity.OrderStatus(status)
		o.PaymentID = pgUUIDPtr(paymentID)
		if addrRaw != nil {
			_ = json.Unmarshal(addrRaw, &o.ShippingAddress)
		}
		total = rowTotal
		orders = append(orders, o)
	}
	return orders, total, rows.Err()
}

func (r *OrderRepository) GetByUserID(ctx context.Context, userID uuid.UUID, page, limit int) ([]*entity.Order, int64, error) {
	return r.List(ctx, repository.OrderFilter{UserID: &userID, Page: page, Limit: limit})
}

// ── PaymentRepository ────────────────────────────────────────────────────────

type PaymentRepository struct{ pool *pgxpool.Pool }

func NewPaymentRepository(pool *pgxpool.Pool) *PaymentRepository {
	return &PaymentRepository{pool: pool}
}

const paymentCols = `
	id, order_id, user_id, provider, status, amount, currency,
	external_id, idempotency_key, provider_metadata, failure_reason, refunded_amount,
	created_at, updated_at`

func scanPayment(row pgx.CollectableRow) (*entity.Payment, error) {
	p := &entity.Payment{}
	var (
		id, orderID, userID [16]byte
		provider, status    string
		metaRaw             json.RawMessage
	)
	if err := row.Scan(
		&id, &orderID, &userID, &provider, &status, &p.Amount, &p.Currency,
		&p.ExternalID, &p.IdempotencyKey, &metaRaw, &p.FailureReason, &p.RefundedAmount,
		&p.CreatedAt, &p.UpdatedAt,
	); err != nil {
		return nil, err
	}
	p.ID = uuid.UUID(id)
	p.OrderID = uuid.UUID(orderID)
	p.UserID = uuid.UUID(userID)
	p.Provider = entity.PaymentProvider(provider)
	p.Status = entity.PaymentStatus(status)
	p.ProviderMetadata = map[string]string{}
	if metaRaw != nil {
		_ = json.Unmarshal(metaRaw, &p.ProviderMetadata)
	}
	return p, nil
}

func (r *PaymentRepository) Create(ctx context.Context, p *entity.Payment) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO payments (id, order_id, user_id, provider, status, amount, currency,
		                      external_id, idempotency_key, provider_metadata, failure_reason,
		                      refunded_amount, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14)`,
		uuidBytes(p.ID), uuidBytes(p.OrderID), uuidBytes(p.UserID),
		string(p.Provider), string(p.Status), p.Amount, p.Currency,
		p.ExternalID, p.IdempotencyKey, marshalJSON(p.ProviderMetadata),
		p.FailureReason, p.RefundedAmount, p.CreatedAt, p.UpdatedAt,
	)
	return err
}

func (r *PaymentRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Payment, error) {
	rows, err := r.pool.Query(ctx, `SELECT `+paymentCols+` FROM payments WHERE id=$1`, uuidBytes(id))
	if err != nil {
		return nil, err
	}
	return pgx.CollectOneRow(rows, scanPayment)
}

func (r *PaymentRepository) GetByOrderID(ctx context.Context, orderID uuid.UUID) (*entity.Payment, error) {
	rows, err := r.pool.Query(ctx, `SELECT `+paymentCols+` FROM payments WHERE order_id=$1`, uuidBytes(orderID))
	if err != nil {
		return nil, err
	}
	return pgx.CollectOneRow(rows, scanPayment)
}

func (r *PaymentRepository) GetByExternalID(ctx context.Context, provider entity.PaymentProvider, externalID string) (*entity.Payment, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT `+paymentCols+` FROM payments WHERE provider=$1 AND external_id=$2`,
		string(provider), externalID)
	if err != nil {
		return nil, err
	}
	return pgx.CollectOneRow(rows, scanPayment)
}

func (r *PaymentRepository) Update(ctx context.Context, p *entity.Payment) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE payments SET
			status=$1, external_id=$2, provider_metadata=$3, failure_reason=$4,
			refunded_amount=$5, updated_at=NOW()
		WHERE id=$6`,
		string(p.Status), p.ExternalID, marshalJSON(p.ProviderMetadata),
		p.FailureReason, p.RefundedAmount, uuidBytes(p.ID),
	)
	return err
}

func (r *PaymentRepository) List(ctx context.Context, f repository.PaymentFilter) ([]*entity.Payment, int64, error) {
	where := []string{"1=1"}
	args := []any{}
	n := 1

	if f.Status != "" {
		where = append(where, fmt.Sprintf("status = $%d", n))
		args = append(args, string(f.Status))
		n++
	}
	if f.Provider != "" {
		where = append(where, fmt.Sprintf("provider = $%d", n))
		args = append(args, string(f.Provider))
		n++
	}

	limit := 20
	page := 1
	if f.Limit > 0 {
		limit = f.Limit
	}
	if f.Page > 0 {
		page = f.Page
	}

	query := fmt.Sprintf(`
		SELECT `+paymentCols+`, COUNT(*) OVER() AS total
		FROM payments WHERE %s
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

	var payments []*entity.Payment
	var total int64
	for rows.Next() {
		p := &entity.Payment{}
		var (
			id, orderID, userID [16]byte
			provider, status    string
			metaRaw             json.RawMessage
			rowTotal            int64
		)
		if err := rows.Scan(
			&id, &orderID, &userID, &provider, &status, &p.Amount, &p.Currency,
			&p.ExternalID, &p.IdempotencyKey, &metaRaw, &p.FailureReason, &p.RefundedAmount,
			&p.CreatedAt, &p.UpdatedAt, &rowTotal,
		); err != nil {
			return nil, 0, err
		}
		p.ID = uuid.UUID(id)
		p.OrderID = uuid.UUID(orderID)
		p.UserID = uuid.UUID(userID)
		p.Provider = entity.PaymentProvider(provider)
		p.Status = entity.PaymentStatus(status)
		p.ProviderMetadata = map[string]string{}
		if metaRaw != nil {
			_ = json.Unmarshal(metaRaw, &p.ProviderMetadata)
		}
		total = rowTotal
		payments = append(payments, p)
	}
	return payments, total, rows.Err()
}

// ── CartRepository ───────────────────────────────────────────────────────────

type CartRepository struct{ pool *pgxpool.Pool }

func NewCartRepository(pool *pgxpool.Pool) *CartRepository {
	return &CartRepository{pool: pool}
}

func (r *CartRepository) GetOrCreate(ctx context.Context, userID uuid.UUID) (*entity.Cart, error) {
	cart := &entity.Cart{}
	var id, uid [16]byte

	err := r.pool.QueryRow(ctx, `
		INSERT INTO carts (id, user_id) VALUES ($1,$2)
		ON CONFLICT (user_id) DO UPDATE SET updated_at=NOW()
		RETURNING id, user_id, created_at, updated_at`,
		uuidBytes(uuid.New()), uuidBytes(userID),
	).Scan(&id, &uid, &cart.CreatedAt, &cart.UpdatedAt)
	if err != nil {
		return nil, err
	}
	cart.ID = uuid.UUID(id)
	cart.UserID = uuid.UUID(uid)
	cart.Items, err = r.fetchCartItems(ctx, cart.ID)
	return cart, err
}

func (r *CartRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*entity.Cart, error) {
	cart := &entity.Cart{}
	var id, uid [16]byte
	err := r.pool.QueryRow(ctx,
		`SELECT id, user_id, created_at, updated_at FROM carts WHERE user_id=$1`, uuidBytes(userID)).
		Scan(&id, &uid, &cart.CreatedAt, &cart.UpdatedAt)
	if err != nil {
		return nil, err
	}
	cart.ID = uuid.UUID(id)
	cart.UserID = uuid.UUID(uid)
	cart.Items, err = r.fetchCartItems(ctx, cart.ID)
	return cart, err
}

func (r *CartRepository) fetchCartItems(ctx context.Context, cartID uuid.UUID) ([]entity.CartItem, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, cart_id, product_id, variant_id, quantity, added_at FROM cart_items WHERE cart_id=$1`,
		uuidBytes(cartID))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []entity.CartItem
	for rows.Next() {
		var item entity.CartItem
		var id, cid, pid [16]byte
		var vid pgtype.UUID
		if err := rows.Scan(&id, &cid, &pid, &vid, &item.Quantity, &item.AddedAt); err != nil {
			return nil, err
		}
		item.ID = uuid.UUID(id)
		item.CartID = uuid.UUID(cid)
		item.ProductID = uuid.UUID(pid)
		item.VariantID = pgUUIDPtr(vid)
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *CartRepository) UpsertItem(ctx context.Context, cartID, productID uuid.UUID, variantID *uuid.UUID, qty int) error {
	if variantID != nil {
		_, err := r.pool.Exec(ctx, `
			INSERT INTO cart_items (id, cart_id, product_id, variant_id, quantity, added_at)
			VALUES ($1,$2,$3,$4,$5,NOW())
			ON CONFLICT (cart_id, product_id, variant_id) DO UPDATE SET quantity = cart_items.quantity + EXCLUDED.quantity`,
			uuidBytes(uuid.New()), uuidBytes(cartID), uuidBytes(productID), uuidBytes(*variantID), qty,
		)
		return err
	}
	// NULL variant: try update first, insert if no row exists.
	tag, err := r.pool.Exec(ctx,
		`UPDATE cart_items SET quantity=quantity+$1 WHERE cart_id=$2 AND product_id=$3 AND variant_id IS NULL`,
		qty, uuidBytes(cartID), uuidBytes(productID))
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		_, err = r.pool.Exec(ctx,
			`INSERT INTO cart_items (id, cart_id, product_id, variant_id, quantity, added_at) VALUES ($1,$2,$3,NULL,$4,NOW())`,
			uuidBytes(uuid.New()), uuidBytes(cartID), uuidBytes(productID), qty)
	}
	return err
}

func (r *CartRepository) RemoveItem(ctx context.Context, cartID, productID uuid.UUID, variantID *uuid.UUID) error {
	if variantID != nil {
		_, err := r.pool.Exec(ctx,
			`DELETE FROM cart_items WHERE cart_id=$1 AND product_id=$2 AND variant_id=$3`,
			uuidBytes(cartID), uuidBytes(productID), uuidBytes(*variantID))
		return err
	}
	_, err := r.pool.Exec(ctx,
		`DELETE FROM cart_items WHERE cart_id=$1 AND product_id=$2 AND variant_id IS NULL`,
		uuidBytes(cartID), uuidBytes(productID))
	return err
}

func (r *CartRepository) Clear(ctx context.Context, cartID uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM cart_items WHERE cart_id=$1`, uuidBytes(cartID))
	return err
}
