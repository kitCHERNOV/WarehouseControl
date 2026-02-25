-- Migration: 00003_seed_data (down)
-- Description: Remove seed data

-- Delete seed items
DELETE FROM items WHERE sku IN (
    'SKU-LAPTOP-001',
    'SKU-MOUSE-001',
    'SKU-KEYBOARD-001',
    'SKU-MONITOR-001',
    'SKU-HUB-001'
);

-- Delete seed users
DELETE FROM users WHERE username IN ('admin', 'manager', 'viewer');
