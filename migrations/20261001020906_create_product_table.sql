-- +goose Up
CREATE TABLE IF NOT EXISTS products(
    id UUID PRIMARY KEY default gen_random_uuid(),
    brand_id UUID REFERENCES brands(id),
    category_id UUID REFERENCES categories(id),
    name varchar(255) not null,
    slug varchar(280) not null unique,
    description text,
    status product_status not null default 'draft',
    base_price decimal(19,2) not null default 0,
    base_weight integer not null default(0),
    sku varchar(100) unique,
     created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);
CREATE INDEX idx_product_brand_category ON products(brand_id, category_id, slug, status);

-- +goose Down
DROP INDEX IF EXISTS idx_product_brand_category;
DROP TABLE IF EXISTS products;
