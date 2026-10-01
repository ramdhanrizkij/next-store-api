-- +goose Up
CREATE TABLE IF NOT EXISTS product_images(
    id UUID PRIMARY KEY default gen_random_uuid(),
    product_id UUID NOT NULL REFERENCES products(id),
    image_url text NOT NULL,
    alt_text varchar(255),
    sort_order integer not null default 0,
    is_primary boolean not null default false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_product_image_id ON product_images(product_id);

-- +goose Down
DROP INDEX IF EXISTS idx_product_image_id;
DROP TABLE IF EXISTS product_images;
