-- Migration: 00002_create_triggers (down)
-- Description: Remove triggers for history logging

-- Drop triggers
DROP TRIGGER IF EXISTS trigger_items_insert ON items;
DROP TRIGGER IF EXISTS trigger_items_update ON items;
DROP TRIGGER IF EXISTS trigger_items_delete ON items;

-- Drop function
DROP FUNCTION IF EXISTS log_item_change();
