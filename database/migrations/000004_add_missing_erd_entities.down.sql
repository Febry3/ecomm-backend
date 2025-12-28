-- Rollback: Remove missing ERD entities and fields

-- 4. Remove reference columns from inventory_ledgers
DROP INDEX IF EXISTS idx_inventory_ledgers_reference;
ALTER TABLE inventory_ledgers DROP COLUMN IF EXISTS reference_id;
ALTER TABLE inventory_ledgers DROP COLUMN IF EXISTS reference_type;

-- 3. Remove columns from group_buy_sessions
ALTER TABLE group_buy_sessions DROP COLUMN IF EXISTS final_tier_id;
ALTER TABLE group_buy_sessions DROP COLUMN IF EXISTS current_participants;

-- 2. Remove version column from product_variant_stocks
ALTER TABLE product_variant_stocks DROP COLUMN IF EXISTS version;

-- 1. Drop stock_reservations table
DROP INDEX IF EXISTS idx_stock_reservations_product_variant;
DROP INDEX IF EXISTS idx_stock_reservations_expires_at;
DROP INDEX IF EXISTS idx_stock_reservations_status;
DROP TABLE IF EXISTS stock_reservations;
