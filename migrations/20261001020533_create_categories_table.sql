-- +goose Up
CREATE TABLE IF NOT EXISTS categories(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    parent_id UUID REFERENCES categories(id),
    name varchar(150) not null,
    slug varchar(180) not null unique,
    description text,
    image_url text,
    is_active boolean not null default true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_categories_id ON categories(parent_id,slug);

-- +goose Down
DROP INDEX IF EXISTS idx_categories_id;
DROP TABLE IF EXISTS categories;
