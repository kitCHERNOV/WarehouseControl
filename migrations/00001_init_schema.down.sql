-- Migration: 00001_init_schema (down)
-- Description: Rollback database schema by dropping tables and indexes

-- Drop indexes (if they exist)
DROP INDEX IF EXISTS idx_items_sku;
DROP INDEX IF EXISTS idx_item_history_changed_at;
DROP INDEX IF EXISTS idx_item_history_user_id;
DROP INDEX IF EXISTS idx_item_history_item_id;

-- Drop tables (order matters due to foreign key constraints)
DROP TABLE IF EXISTS item_history;
DROP TABLE IF EXISTS items;
DROP TABLE IF EXISTS users;
