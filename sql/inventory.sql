CREATE TABLE inventories (
    code VARCHAR(50) PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    stock INT NOT NULL DEFAULT 0,
    description TEXT,
    status VARCHAR(50) NOT NULL DEFAULT ''
)

INSERT INTO inventories (code, name, stock, description, status) VALUES
('INV001', 'Laptop', 25, 'Dell Latitude 5420', 'active'),
('INV002', 'Mouse', 100, 'Logitech wireless mouse', 'active'),
('INV003', 'Keyboard', 75, 'Mechanical keyboard with RGB lights', 'active'),
('INV004', 'Monitor', 30, '27-inch 4K UHD monitor', 'active'),
('INV005', 'Printer', 10, 'HP LaserJet Pro multifunction printer', 'active'),
('INV006', 'Desk Chair', 40, 'Ergonomic office chair', 'active'),
('INV007', 'Webcam', 60, 'HD webcam with built-in mic', 'active'),
('INV008', 'Router', 20, 'Wi-Fi 6 Dual-Band Router', 'active'),
('INV009', 'USB Hub', 85, '7-port powered USB hub', 'active'),
('INV010', 'External HDD', 15, '2TB Seagate USB 3.0 external hard drive', 'active'),
('INV011', 'Projector', 5, 'Full HD business projector', 'inactive'),
('INV012', 'Scanner', 8, 'Flatbed document scanner', 'inactive'),
('INV013', 'Desk Lamp', 50, 'LED lamp with brightness control', 'active'),
('INV014', 'Headphones', 35, 'Noise-cancelling over-ear headphones', 'active'),
('INV015', 'Laptop Stand', 45, 'Adjustable aluminum laptop stand', 'active');