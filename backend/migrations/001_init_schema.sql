-- +goose Up
-- +goose StatementBegin

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pg_trgm";

-- ─── ENUMS ───────────────────────────────────────────────────────────────────

CREATE TYPE user_role     AS ENUM ('admin', 'user', 'seller');
CREATE TYPE user_status   AS ENUM ('active', 'inactive', 'banned', 'pending');
CREATE TYPE lang          AS ENUM ('ru', 'tj');
CREATE TYPE product_status AS ENUM ('draft', 'pending', 'active', 'inactive', 'rejected');
CREATE TYPE order_status  AS ENUM ('created', 'paid', 'processing', 'shipped', 'delivered', 'closed', 'cancelled', 'refunded');
CREATE TYPE payment_provider AS ENUM ('stripe', 'visa', 'alif');
CREATE TYPE payment_status   AS ENUM ('pending', 'succeeded', 'failed', 'refunded', 'cancelled');
CREATE TYPE seller_status    AS ENUM ('pending', 'approved', 'rejected', 'suspended');

-- ─── USERS ───────────────────────────────────────────────────────────────────

CREATE TABLE users (
    id            UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email         VARCHAR(255) NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    first_name    VARCHAR(100) NOT NULL,
    last_name     VARCHAR(100) NOT NULL DEFAULT '',
    phone         VARCHAR(30)  NOT NULL DEFAULT '',
    role          user_role    NOT NULL DEFAULT 'user',
    status        user_status  NOT NULL DEFAULT 'active',
    language      lang         NOT NULL DEFAULT 'ru',
    avatar_url    TEXT         NOT NULL DEFAULT '',
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_users_email  ON users(email);
CREATE INDEX idx_users_role   ON users(role);
CREATE INDEX idx_users_status ON users(status);

-- ─── REFRESH TOKENS ──────────────────────────────────────────────────────────

CREATE TABLE refresh_tokens (
    id         UUID        PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id    UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token      TEXT        NOT NULL UNIQUE,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens(user_id);
CREATE INDEX idx_refresh_tokens_token   ON refresh_tokens(token);

-- ─── SELLERS ─────────────────────────────────────────────────────────────────

CREATE TABLE sellers (
    id              UUID         PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id         UUID         NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    business_name   VARCHAR(255) NOT NULL,
    business_email  VARCHAR(255) NOT NULL,
    business_phone  VARCHAR(30)  NOT NULL DEFAULT '',
    country         VARCHAR(100) NOT NULL DEFAULT '',
    status          seller_status NOT NULL DEFAULT 'pending',
    commission_rate NUMERIC(5,2) NOT NULL DEFAULT 10.00,
    description     TEXT         NOT NULL DEFAULT '',
    logo_url        TEXT         NOT NULL DEFAULT '',
    verified_at     TIMESTAMPTZ,
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_sellers_user_id ON sellers(user_id);
CREATE INDEX idx_sellers_status  ON sellers(status);

-- ─── CATEGORIES ──────────────────────────────────────────────────────────────

CREATE TABLE categories (
    id          UUID        PRIMARY KEY DEFAULT uuid_generate_v4(),
    parent_id   UUID        REFERENCES categories(id) ON DELETE SET NULL,
    name        VARCHAR(255) NOT NULL,
    slug        VARCHAR(255) NOT NULL UNIQUE,
    description TEXT         NOT NULL DEFAULT '',
    image_url   TEXT         NOT NULL DEFAULT '',
    sort_order  INT          NOT NULL DEFAULT 0,
    is_active   BOOLEAN      NOT NULL DEFAULT TRUE,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_categories_parent_id ON categories(parent_id);
CREATE INDEX idx_categories_slug      ON categories(slug);

-- ─── PRODUCTS ────────────────────────────────────────────────────────────────

CREATE TABLE products (
    id              UUID           PRIMARY KEY DEFAULT uuid_generate_v4(),
    seller_id       UUID           NOT NULL REFERENCES sellers(id) ON DELETE CASCADE,
    category_id     UUID           NOT NULL REFERENCES categories(id) ON DELETE RESTRICT,
    name            VARCHAR(500)   NOT NULL,
    slug            VARCHAR(600)   NOT NULL UNIQUE,
    description     TEXT           NOT NULL DEFAULT '',
    short_desc      VARCHAR(1000)  NOT NULL DEFAULT '',
    base_price      NUMERIC(12,2)  NOT NULL,
    sale_price      NUMERIC(12,2),
    currency        VARCHAR(3)     NOT NULL DEFAULT 'USD',
    sku             VARCHAR(100)   NOT NULL DEFAULT '',
    stock           INT            NOT NULL DEFAULT 0,
    status          product_status NOT NULL DEFAULT 'pending',
    images          TEXT[]         NOT NULL DEFAULT '{}',
    tags            TEXT[]         NOT NULL DEFAULT '{}',
    specifications  JSONB          NOT NULL DEFAULT '{}',
    weight          NUMERIC(8,3)   NOT NULL DEFAULT 0,
    dimensions      JSONB          NOT NULL DEFAULT '{}',
    view_count      INT            NOT NULL DEFAULT 0,
    sold_count      INT            NOT NULL DEFAULT 0,
    rating          NUMERIC(3,2)   NOT NULL DEFAULT 0,
    review_count    INT            NOT NULL DEFAULT 0,
    featured_at     TIMESTAMPTZ,
    created_at      TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ    NOT NULL DEFAULT NOW(),

    CONSTRAINT chk_base_price_positive CHECK (base_price >= 0),
    CONSTRAINT chk_stock_non_negative  CHECK (stock >= 0)
);

CREATE INDEX idx_products_seller_id   ON products(seller_id);
CREATE INDEX idx_products_category_id ON products(category_id);
CREATE INDEX idx_products_status      ON products(status);
CREATE INDEX idx_products_base_price  ON products(base_price);
CREATE INDEX idx_products_name_trgm   ON products USING GIN (name gin_trgm_ops);
CREATE INDEX idx_products_tags        ON products USING GIN (tags);

-- ─── PRODUCT VARIANTS ────────────────────────────────────────────────────────

CREATE TABLE product_variants (
    id         UUID          PRIMARY KEY DEFAULT uuid_generate_v4(),
    product_id UUID          NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    name       VARCHAR(255)  NOT NULL,
    sku        VARCHAR(100)  NOT NULL,
    price      NUMERIC(12,2) NOT NULL,
    stock      INT           NOT NULL DEFAULT 0,
    attributes JSONB         NOT NULL DEFAULT '{}',
    image_url  TEXT          NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ   NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_product_variants_product_id ON product_variants(product_id);

-- ─── CARTS ───────────────────────────────────────────────────────────────────

CREATE TABLE carts (
    id         UUID        PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id    UUID        NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE cart_items (
    id         UUID        PRIMARY KEY DEFAULT uuid_generate_v4(),
    cart_id    UUID        NOT NULL REFERENCES carts(id) ON DELETE CASCADE,
    product_id UUID        NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    variant_id UUID        REFERENCES product_variants(id) ON DELETE SET NULL,
    quantity   INT         NOT NULL DEFAULT 1 CHECK (quantity > 0),
    added_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(cart_id, product_id, variant_id)
);

CREATE INDEX idx_cart_items_cart_id ON cart_items(cart_id);

-- ─── ORDERS ──────────────────────────────────────────────────────────────────

CREATE TABLE orders (
    id               UUID          PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id          UUID          NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    status           order_status  NOT NULL DEFAULT 'created',
    subtotal         NUMERIC(12,2) NOT NULL DEFAULT 0,
    shipping_cost    NUMERIC(12,2) NOT NULL DEFAULT 0,
    tax              NUMERIC(12,2) NOT NULL DEFAULT 0,
    total            NUMERIC(12,2) NOT NULL DEFAULT 0,
    currency         VARCHAR(3)    NOT NULL DEFAULT 'USD',
    shipping_address JSONB         NOT NULL DEFAULT '{}',
    notes            TEXT          NOT NULL DEFAULT '',
    tracking_number  VARCHAR(100)  NOT NULL DEFAULT '',
    payment_id       UUID,
    created_at       TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ   NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_orders_user_id    ON orders(user_id);
CREATE INDEX idx_orders_status     ON orders(status);
CREATE INDEX idx_orders_created_at ON orders(created_at DESC);

CREATE TABLE order_items (
    id          UUID          PRIMARY KEY DEFAULT uuid_generate_v4(),
    order_id    UUID          NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    product_id  UUID          NOT NULL REFERENCES products(id) ON DELETE RESTRICT,
    variant_id  UUID          REFERENCES product_variants(id) ON DELETE SET NULL,
    seller_id   UUID          NOT NULL REFERENCES sellers(id) ON DELETE RESTRICT,
    name        VARCHAR(500)  NOT NULL,
    sku         VARCHAR(100)  NOT NULL DEFAULT '',
    image_url   TEXT          NOT NULL DEFAULT '',
    quantity    INT           NOT NULL CHECK (quantity > 0),
    unit_price  NUMERIC(12,2) NOT NULL,
    total_price NUMERIC(12,2) NOT NULL
);

CREATE INDEX idx_order_items_order_id  ON order_items(order_id);
CREATE INDEX idx_order_items_seller_id ON order_items(seller_id);

CREATE TABLE order_status_history (
    id         UUID         PRIMARY KEY DEFAULT uuid_generate_v4(),
    order_id   UUID         NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    status     order_status NOT NULL,
    changed_by UUID         REFERENCES users(id) ON DELETE SET NULL,
    note       TEXT         NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_order_status_history_order_id ON order_status_history(order_id);

-- ─── PAYMENTS ────────────────────────────────────────────────────────────────

CREATE TABLE payments (
    id                UUID             PRIMARY KEY DEFAULT uuid_generate_v4(),
    order_id          UUID             NOT NULL REFERENCES orders(id) ON DELETE RESTRICT,
    user_id           UUID             NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    provider          payment_provider NOT NULL,
    status            payment_status   NOT NULL DEFAULT 'pending',
    amount            NUMERIC(12,2)    NOT NULL,
    currency          VARCHAR(3)       NOT NULL DEFAULT 'USD',
    external_id       TEXT             NOT NULL DEFAULT '',
    idempotency_key   TEXT             NOT NULL UNIQUE,
    provider_metadata JSONB            NOT NULL DEFAULT '{}',
    failure_reason    TEXT             NOT NULL DEFAULT '',
    refunded_amount   NUMERIC(12,2)    NOT NULL DEFAULT 0,
    created_at        TIMESTAMPTZ      NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ      NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_payments_order_id   ON payments(order_id);
CREATE INDEX idx_payments_user_id    ON payments(user_id);
CREATE INDEX idx_payments_external   ON payments(provider, external_id);
CREATE INDEX idx_payments_status     ON payments(status);

-- Deferred FK from orders → payments (avoids circular creation issue)
ALTER TABLE orders
    ADD CONSTRAINT fk_orders_payment_id
    FOREIGN KEY (payment_id) REFERENCES payments(id) ON DELETE SET NULL;

-- ─── AUTO updated_at TRIGGER ─────────────────────────────────────────────────

CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER LANGUAGE plpgsql AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$;

CREATE TRIGGER trg_users_updated_at            BEFORE UPDATE ON users            FOR EACH ROW EXECUTE FUNCTION set_updated_at();
CREATE TRIGGER trg_sellers_updated_at          BEFORE UPDATE ON sellers          FOR EACH ROW EXECUTE FUNCTION set_updated_at();
CREATE TRIGGER trg_categories_updated_at       BEFORE UPDATE ON categories       FOR EACH ROW EXECUTE FUNCTION set_updated_at();
CREATE TRIGGER trg_products_updated_at         BEFORE UPDATE ON products         FOR EACH ROW EXECUTE FUNCTION set_updated_at();
CREATE TRIGGER trg_product_variants_updated_at BEFORE UPDATE ON product_variants FOR EACH ROW EXECUTE FUNCTION set_updated_at();
CREATE TRIGGER trg_carts_updated_at            BEFORE UPDATE ON carts            FOR EACH ROW EXECUTE FUNCTION set_updated_at();
CREATE TRIGGER trg_orders_updated_at           BEFORE UPDATE ON orders           FOR EACH ROW EXECUTE FUNCTION set_updated_at();
CREATE TRIGGER trg_payments_updated_at         BEFORE UPDATE ON payments         FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS order_status_history;
DROP TABLE IF EXISTS order_items;
ALTER TABLE orders DROP CONSTRAINT IF EXISTS fk_orders_payment_id;
DROP TABLE IF EXISTS payments;
DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS cart_items;
DROP TABLE IF EXISTS carts;
DROP TABLE IF EXISTS product_variants;
DROP TABLE IF EXISTS products;
DROP TABLE IF EXISTS categories;
DROP TABLE IF EXISTS sellers;
DROP TABLE IF EXISTS refresh_tokens;
DROP TABLE IF EXISTS users;
DROP FUNCTION IF EXISTS set_updated_at;
DROP TYPE IF EXISTS seller_status;
DROP TYPE IF EXISTS payment_status;
DROP TYPE IF EXISTS payment_provider;
DROP TYPE IF EXISTS order_status;
DROP TYPE IF EXISTS product_status;
DROP TYPE IF EXISTS lang;
DROP TYPE IF EXISTS user_status;
DROP TYPE IF EXISTS user_role;

-- +goose StatementEnd
