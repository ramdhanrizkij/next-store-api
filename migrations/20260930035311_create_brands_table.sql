-- +goose Up
CREATE TABLE IF NOT EXISTS brands(
    id UUID NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    name varchar(150) NOT NULL,
    slug varchar(180) NOT NULL UNIQUE,
    description text,
    logo_url text,
    is_active boolean NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_brands_slug on brands(slug);

-- +goose Down
DROP INDEX IF EXISTS idx_brands_slug;
DROP TABLE IF EXISTS brands;
