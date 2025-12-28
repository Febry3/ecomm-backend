-- Migration: Add missing ERD entities and fields
-- Created: 2025-12-28

-- 1. Create stock_reservations table
CREATE TABLE IF NOT EXISTS stock_reservations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_variant_id UUID NOT NULL REFERENCES product_variants(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    order_id UUID REFERENCES orders(id) ON DELETE SET NULL,
    quantity INTEGER NOT NULL,
    status TEXT DEFAULT 'pending' CHECK (status IN ('pending', 'completed', 'expired')),
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_stock_reservations_status ON stock_reservations(status);
CREATE INDEX idx_stock_reservations_expires_at ON stock_reservations(expires_at);
CREATE INDEX idx_stock_reservations_product_variant ON stock_reservations(product_variant_id);

-- 2. Add version column to product_variant_stocks for optimistic locking
ALTER TABLE product_variant_stocks 
ADD COLUMN IF NOT EXISTS version INTEGER DEFAULT 1;

-- 3. Add current_participants and final_tier_id columns to group_buy_sessions
ALTER TABLE group_buy_sessions 
ADD COLUMN IF NOT EXISTS current_participants INTEGER DEFAULT 0;

ALTER TABLE group_buy_sessions 
ADD COLUMN IF NOT EXISTS final_tier_id UUID REFERENCES group_buy_tiers(id) ON DELETE SET NULL;

-- 4. Add reference_type and reference_id columns to inventory_ledgers
ALTER TABLE inventory_ledgers 
ADD COLUMN IF NOT EXISTS reference_type TEXT;

ALTER TABLE inventory_ledgers 
ADD COLUMN IF NOT EXISTS reference_id UUID;

CREATE INDEX idx_inventory_ledgers_reference ON inventory_ledgers(reference_type, reference_id);
