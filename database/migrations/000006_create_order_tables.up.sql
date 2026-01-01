-- Migration: Create order and payment tables
-- Created: 2026-01-01

-- Orders table
CREATE TABLE IF NOT EXISTS orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_number VARCHAR(50) NOT NULL UNIQUE,
    user_id BIGINT NOT NULL REFERENCES users(id),
    buyer_group_session_id UUID REFERENCES buyer_group_sessions(id) ON DELETE SET NULL,
    seller_id BIGINT NOT NULL REFERENCES sellers(id),
    product_variant_id UUID NOT NULL REFERENCES product_variants(id),
    quantity INTEGER NOT NULL DEFAULT 1,
    price_at_order DECIMAL(15,2) NOT NULL,
    subtotal DECIMAL(15,2) NOT NULL,
    delivery_charge DECIMAL(15,2) DEFAULT 0,
    total_amount DECIMAL(15,2) NOT NULL,
    status TEXT DEFAULT 'pending_payment' CHECK (status IN (
        'pending_payment', 'paid', 'processing', 'shipped', 'delivered', 'cancelled', 'expired'
    )),
    address_id UUID NOT NULL REFERENCES addresses(address_id),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Payments table
CREATE TABLE IF NOT EXISTS payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    amount DECIMAL(15,2) NOT NULL,
    status TEXT DEFAULT 'pending' CHECK (status IN ('pending', 'settlement', 'expire', 'cancel', 'deny')),
    payment_method VARCHAR(50) NOT NULL DEFAULT 'bank_transfer',
    bank_code VARCHAR(20) NOT NULL,
    va_number VARCHAR(50),
    bill_key VARCHAR(50),
    biller_code VARCHAR(20),
    gateway_transaction_id VARCHAR(255),
    expired_at TIMESTAMPTZ NOT NULL,
    paid_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Order Shipping Details (snapshot of address at order time)
CREATE TABLE IF NOT EXISTS order_shipping_details (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID NOT NULL UNIQUE REFERENCES orders(id) ON DELETE CASCADE,
    receiver_name VARCHAR(255) NOT NULL,
    phone VARCHAR(20),
    street_address TEXT NOT NULL,
    rt VARCHAR(5),
    rw VARCHAR(5),
    village VARCHAR(100),
    district VARCHAR(100),
    city VARCHAR(100) NOT NULL,
    province VARCHAR(100) NOT NULL,
    postal_code VARCHAR(10) NOT NULL,
    notes TEXT
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_orders_user_id ON orders(user_id);
CREATE INDEX IF NOT EXISTS idx_orders_seller_id ON orders(seller_id);
CREATE INDEX IF NOT EXISTS idx_orders_status ON orders(status);
CREATE INDEX IF NOT EXISTS idx_orders_buyer_group_session ON orders(buyer_group_session_id);
CREATE INDEX IF NOT EXISTS idx_orders_created_at ON orders(created_at DESC);

CREATE INDEX IF NOT EXISTS idx_payments_order_id ON payments(order_id);
CREATE INDEX IF NOT EXISTS idx_payments_status ON payments(status);
CREATE INDEX IF NOT EXISTS idx_payments_expired_at ON payments(expired_at);
CREATE INDEX IF NOT EXISTS idx_payments_gateway_tx_id ON payments(gateway_transaction_id);
