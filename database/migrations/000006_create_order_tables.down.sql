-- Rollback: Drop order and payment tables

DROP INDEX IF EXISTS idx_payments_gateway_tx_id;
DROP INDEX IF EXISTS idx_payments_expired_at;
DROP INDEX IF EXISTS idx_payments_status;
DROP INDEX IF EXISTS idx_payments_order_id;

DROP INDEX IF EXISTS idx_orders_created_at;
DROP INDEX IF EXISTS idx_orders_buyer_group_session;
DROP INDEX IF EXISTS idx_orders_status;
DROP INDEX IF EXISTS idx_orders_seller_id;
DROP INDEX IF EXISTS idx_orders_user_id;

DROP TABLE IF EXISTS order_shipping_details;
DROP TABLE IF EXISTS payments;
DROP TABLE IF EXISTS orders;
