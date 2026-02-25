-- Migration: 00003_seed_data
-- Description: Insert seed data for testing (users and initial items)

-- Insert test users
-- Note: Passwords are NOT hashed - this is for development only!
-- In production, use bcrypt hashed passwords
INSERT INTO users (username, password, role) VALUES
('admin', 'admin123', 'admin'),
('manager', 'manager123', 'manager'),
('viewer', 'viewer123', 'viewer')
ON CONFLICT (username) DO NOTHING;

-- Insert initial items for testing
INSERT INTO items (name, quantity, price, description, sku) VALUES
('Laptop Dell XPS 15', 10, 1299.99, 'High-performance laptop with 16GB RAM', 'SKU-LAPTOP-001'),
('Mouse Logitech MX Master', 50, 99.99, 'Wireless ergonomic mouse', 'SKU-MOUSE-001'),
('Keyboard Keychron K2', 30, 89.99, 'Mechanical keyboard with RGB backlight', 'SKU-KEYBOARD-001'),
('Monitor LG 27UK850', 15, 349.99, '27" 4K UHD IPS monitor', 'SKU-MONITOR-001'),
('USB-C Hub Adapter', 100, 29.99, 'Multi-port USB-C hub with HDMI and USB 3.0', 'SKU-HUB-001')
ON CONFLICT (sku) DO NOTHING;
