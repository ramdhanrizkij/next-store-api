Project ecommerce {
  database_type: "PostgreSQL"
  Note: "Single-store ecommerce with product variants, guest/login checkout, RBAC, and payment gateway integration"
}

/*
|--------------------------------------------------------------------------
| ENUMS
|--------------------------------------------------------------------------
*/

Enum user_status {
  active
  inactive
  suspended
}

Enum product_status {
  draft
  active
  inactive
  archived
}

Enum cart_status {
  active
  abandoned
  converted
}

Enum order_status {
  pending
  confirmed
  processing
  shipped
  delivered
  cancelled
  completed
  refunded
}

Enum payment_status {
  pending
  paid
  failed
  expired
  cancelled
  refunded
  partially_refunded
}

Enum payment_gateway {
  midtrans
  xendit
}

Enum payment_transaction_status {
  pending
  processing
  success
  failed
  expired
  cancelled
}

Enum address_type {
  billing
  shipping
}

Enum discount_type {
  percentage
  fixed
}

/*
|--------------------------------------------------------------------------
| USERS & RBAC
|--------------------------------------------------------------------------
*/

Table users {
  id uuid [pk]

  name varchar(150) [not null]
  email varchar(255) [not null, unique]
  password_hash varchar(255)
  phone varchar(50)

  status user_status [not null, default: 'active']

  email_verified_at timestamptz
  last_login_at timestamptz

  created_at timestamptz [not null, default: `now()`]
  updated_at timestamptz [not null, default: `now()`]
  deleted_at timestamptz

  indexes {
    email
    phone
  }
}

Table roles {
  id uuid [pk]

  name varchar(100) [not null, unique]
  slug varchar(100) [not null, unique]
  description text

  created_at timestamptz [not null, default: `now()`]
  updated_at timestamptz [not null, default: `now()`]
  deleted_at timestamptz
}

Table permissions {
  id uuid [pk]

  name varchar(150) [not null]
  slug varchar(150) [not null, unique]
  description text

  created_at timestamptz [not null, default: `now()`]
  updated_at timestamptz [not null, default: `now()`]
  deleted_at timestamptz
}

Table user_roles {
  id uuid [pk]

  user_id uuid [not null]
  role_id uuid [not null]

  created_at timestamptz [not null, default: `now()`]
  updated_at timestamptz [not null, default: `now()`]
  deleted_at timestamptz

  indexes {
    (user_id, role_id) [unique]
  }
}

Table role_permissions {
  id uuid [pk]

  role_id uuid [not null]
  permission_id uuid [not null]

  created_at timestamptz [not null, default: `now()`]
  updated_at timestamptz [not null, default: `now()`]
  deleted_at timestamptz

  indexes {
    (role_id, permission_id) [unique]
  }
}

/*
|--------------------------------------------------------------------------
| CATALOG
|--------------------------------------------------------------------------
*/

Table brands {
  id uuid [pk]

  name varchar(150) [not null]
  slug varchar(180) [not null, unique]
  description text
  logo_url text

  is_active boolean [not null, default: true]

  created_at timestamptz [not null, default: `now()`]
  updated_at timestamptz [not null, default: `now()`]
  deleted_at timestamptz

  indexes {
    slug
  }
}

Table categories {
  id uuid [pk]

  parent_id uuid

  name varchar(150) [not null]
  slug varchar(180) [not null, unique]
  description text
  image_url text

  is_active boolean [not null, default: true]

  created_at timestamptz [not null, default: `now()`]
  updated_at timestamptz [not null, default: `now()`]
  deleted_at timestamptz

  indexes {
    parent_id
    slug
  }
}

Table products {
  id uuid [pk]

  brand_id uuid
  category_id uuid

  name varchar(255) [not null]
  slug varchar(280) [not null, unique]
  description text

  status product_status [not null, default: 'draft']

  base_price numeric(19,2) [not null, default: 0]
  base_weight integer [not null, default: 0]

  sku varchar(100) [unique]

  created_at timestamptz [not null, default: `now()`]
  updated_at timestamptz [not null, default: `now()`]
  deleted_at timestamptz

  indexes {
    brand_id
    category_id
    slug
    status
  }
}

Table product_images {
  id uuid [pk]

  product_id uuid [not null]

  image_url text [not null]
  alt_text varchar(255)

  sort_order integer [not null, default: 0]
  is_primary boolean [not null, default: false]

  created_at timestamptz [not null, default: `now()`]
  updated_at timestamptz [not null, default: `now()`]
  deleted_at timestamptz

  indexes {
    product_id
  }
}

/*
|--------------------------------------------------------------------------
| PRODUCT VARIANTS
|--------------------------------------------------------------------------
*/

Table attributes {
  id uuid [pk]

  name varchar(100) [not null]
  slug varchar(120) [not null, unique]

  created_at timestamptz [not null, default: `now()`]
  updated_at timestamptz [not null, default: `now()`]
  deleted_at timestamptz
}

Table attribute_values {
  id uuid [pk]

  attribute_id uuid [not null]

  value varchar(150) [not null]
  slug varchar(170) [not null]

  created_at timestamptz [not null, default: `now()`]
  updated_at timestamptz [not null, default: `now()`]
  deleted_at timestamptz

  indexes {
    (attribute_id, slug) [unique]
  }
}

Table product_variants {
  id uuid [pk]

  product_id uuid [not null]

  sku varchar(100) [not null, unique]

  name varchar(255)

  price numeric(19,2) [not null, default: 0]
  compare_at_price numeric(19,2)

  stock_quantity integer [not null, default: 0]
  weight integer [not null, default: 0]

  status product_status [not null, default: 'active']

  created_at timestamptz [not null, default: `now()`]
  updated_at timestamptz [not null, default: `now()`]
  deleted_at timestamptz

  indexes {
    product_id
    sku
    status
  }
}

Table variant_attribute_values {
  id uuid [pk]

  variant_id uuid [not null]
  attribute_value_id uuid [not null]

  created_at timestamptz [not null, default: `now()`]
  updated_at timestamptz [not null, default: `now()`]
  deleted_at timestamptz

  indexes {
    (variant_id, attribute_value_id) [unique]
  }
}

/*
|--------------------------------------------------------------------------
| CART
|--------------------------------------------------------------------------
*/

Table carts {
  id uuid [pk]

  user_id uuid

  guest_token varchar(255) [unique]

  status cart_status [not null, default: 'active']

  expires_at timestamptz

  created_at timestamptz [not null, default: `now()`]
  updated_at timestamptz [not null, default: `now()`]
  deleted_at timestamptz

  indexes {
    user_id
    guest_token
    status
  }
}

Table cart_items {
  id uuid [pk]

  cart_id uuid [not null]
  product_id uuid [not null]
  variant_id uuid

  quantity integer [not null, default: 1]

  unit_price numeric(19,2) [not null, default: 0]

  created_at timestamptz [not null, default: `now()`]
  updated_at timestamptz [not null, default: `now()`]
  deleted_at timestamptz

  indexes {
    cart_id
    product_id
    variant_id
  }
}

/*
|--------------------------------------------------------------------------
| CUSTOMER ADDRESSES
|--------------------------------------------------------------------------
*/

Table addresses {
  id uuid [pk]

  user_id uuid

  recipient_name varchar(150) [not null]
  phone varchar(50) [not null]

  address_line text [not null]
  province varchar(150)
  city varchar(150)
  district varchar(150)
  postal_code varchar(20)
  country varchar(100) [not null, default: 'Indonesia']

  notes text

  is_default boolean [not null, default: false]

  created_at timestamptz [not null, default: `now()`]
  updated_at timestamptz [not null, default: `now()`]
  deleted_at timestamptz

  indexes {
    user_id
  }
}

/*
|--------------------------------------------------------------------------
| ORDERS
|--------------------------------------------------------------------------
*/

Table orders {
  id uuid [pk]

  user_id uuid

  order_number varchar(50) [not null, unique]

  guest_email varchar(255)
  guest_phone varchar(50)

  status order_status [not null, default: 'pending']
  payment_status payment_status [not null, default: 'pending']

  subtotal numeric(19,2) [not null, default: 0]
  discount_amount numeric(19,2) [not null, default: 0]
  shipping_amount numeric(19,2) [not null, default: 0]
  tax_amount numeric(19,2) [not null, default: 0]

  total_amount numeric(19,2) [not null, default: 0]

  currency varchar(10) [not null, default: 'IDR']

  notes text

  created_at timestamptz [not null, default: `now()`]
  updated_at timestamptz [not null, default: `now()`]
  deleted_at timestamptz

  indexes {
    user_id
    order_number
    status
    payment_status
    created_at
  }
}

Table order_items {
  id uuid [pk]

  order_id uuid [not null]

  product_id uuid
  variant_id uuid

  product_name varchar(255) [not null]
  variant_name varchar(255)

  sku varchar(100)

  quantity integer [not null]

  unit_price numeric(19,2) [not null]
  subtotal numeric(19,2) [not null]

  created_at timestamptz [not null, default: `now()`]
  updated_at timestamptz [not null, default: `now()`]
  deleted_at timestamptz

  indexes {
    order_id
    product_id
    variant_id
  }
}

/*
|--------------------------------------------------------------------------
| ORDER ADDRESS SNAPSHOT
|--------------------------------------------------------------------------
*/

Table order_addresses {
  id uuid [pk]

  order_id uuid [not null]
  type address_type [not null]

  recipient_name varchar(150) [not null]
  phone varchar(50) [not null]

  address_line text [not null]
  province varchar(150)
  city varchar(150)
  district varchar(150)
  postal_code varchar(20)
  country varchar(100) [not null, default: 'Indonesia']

  notes text

  created_at timestamptz [not null, default: `now()`]
  updated_at timestamptz [not null, default: `now()`]
  deleted_at timestamptz

  indexes {
    (order_id, type) [unique]
  }
}

/*
|--------------------------------------------------------------------------
| PAYMENT
|--------------------------------------------------------------------------
*/

Table payments {
  id uuid [pk]

  order_id uuid [not null]

  gateway payment_gateway [not null]

  status payment_status [not null, default: 'pending']

  amount numeric(19,2) [not null]
  currency varchar(10) [not null, default: 'IDR']

  external_payment_id varchar(255)
  external_order_id varchar(255)

  payment_method varchar(100)

  expired_at timestamptz
  paid_at timestamptz

  raw_response jsonb

  created_at timestamptz [not null, default: `now()`]
  updated_at timestamptz [not null, default: `now()`]
  deleted_at timestamptz

  indexes {
    order_id
    gateway
    status
    external_payment_id
    external_order_id
  }
}

Table payment_transactions {
  id uuid [pk]

  payment_id uuid [not null]

  transaction_id varchar(255)
  reference_id varchar(255)

  status payment_transaction_status [not null, default: 'pending']

  amount numeric(19,2) [not null]

  payment_method varchar(100)

  response_code varchar(100)
  response_message text

  raw_request jsonb
  raw_response jsonb

  processed_at timestamptz

  created_at timestamptz [not null, default: `now()`]
  updated_at timestamptz [not null, default: `now()`]
  deleted_at timestamptz

  indexes {
    payment_id
    transaction_id
    reference_id
    status
  }
}

/*
|--------------------------------------------------------------------------
| DISCOUNT / COUPON
|--------------------------------------------------------------------------
*/

Table coupons {
  id uuid [pk]

  code varchar(100) [not null, unique]
  name varchar(150) [not null]
  description text

  discount_type discount_type [not null]
  discount_value numeric(19,2) [not null]

  minimum_order_amount numeric(19,2) [default: 0]
  maximum_discount_amount numeric(19,2)

  usage_limit integer
  usage_count integer [not null, default: 0]

  starts_at timestamptz
  expires_at timestamptz

  is_active boolean [not null, default: true]

  created_at timestamptz [not null, default: `now()`]
  updated_at timestamptz [not null, default: `now()`]
  deleted_at timestamptz

  indexes {
    code
    is_active
  }
}

Table order_coupons {
  id uuid [pk]

  order_id uuid [not null]
  coupon_id uuid [not null]

  discount_amount numeric(19,2) [not null]

  created_at timestamptz [not null, default: `now()`]
  updated_at timestamptz [not null, default: `now()`]
  deleted_at timestamptz

  indexes {
    (order_id, coupon_id) [unique]
  }
}

/*
|--------------------------------------------------------------------------
| RELATIONSHIPS
|--------------------------------------------------------------------------
*/

/* RBAC */

Ref: user_roles.user_id > users.id
Ref: user_roles.role_id > roles.id

Ref: role_permissions.role_id > roles.id
Ref: role_permissions.permission_id > permissions.id

/* Catalog */

Ref: categories.parent_id > categories.id

Ref: products.brand_id > brands.id
Ref: products.category_id > categories.id

Ref: product_images.product_id > products.id

Ref: attribute_values.attribute_id > attributes.id

Ref: product_variants.product_id > products.id

Ref: variant_attribute_values.variant_id > product_variants.id
Ref: variant_attribute_values.attribute_value_id > attribute_values.id

/* Cart */

Ref: carts.user_id > users.id

Ref: cart_items.cart_id > carts.id
Ref: cart_items.product_id > products.id
Ref: cart_items.variant_id > product_variants.id

/* Customer */

Ref: addresses.user_id > users.id

/* Orders */

Ref: orders.user_id > users.id

Ref: order_items.order_id > orders.id
Ref: order_items.product_id > products.id
Ref: order_items.variant_id > product_variants.id

Ref: order_addresses.order_id > orders.id

/* Payment */

Ref: payments.order_id > orders.id

Ref: payment_transactions.payment_id > payments.id

/* Coupon */

Ref: order_coupons.order_id > orders.id
Ref: order_coupons.coupon_id > coupons.id