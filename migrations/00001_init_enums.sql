-- +goose Up
-- +goose StatementBegin
CREATE TYPE user_status as ENUM('active','inactive','suspended');
CREATE TYPE product_status as ENUM('draft','active','inactive','archived');
CREATE TYPE cart_status as ENUM('active','abandoned','converted');
CREATE TYPE order_status as ENUM('pending','confirmed','processing','shipped','delivered','cancelled','completed','refunded');
CREATE TYPE payment_status as ENUM('pending','paid','failed','expired','cancelled','refunded','partially_refunded');
CREATE TYPE payment_gateway as ENUM('midtrans','xendit');
CREATE TYPE payment_transaction_status as ENUM('pending','processing','success','failed','expired','cancelled');
CREATE TYPE address_type as ENUM('billing','shipping');
CREATE TYPE discount_type as ENUM('percentage','fixed');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TYPE IF EXISTS discount_type;
DROP TYPE IF EXISTS address_type;
DROP TYPE IF EXISTS payment_transaction_status;
DROP TYPE IF EXISTS payment_gateway;
DROP TYPE IF EXISTS payment_status;
DROP TYPE IF EXISTS order_status;
DROP TYPE IF EXISTS cart_status;
DROP TYPE IF EXISTS product_status;
DROP TYPE IF EXISTS user_status;
-- +goose StatementEnd
