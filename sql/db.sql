-- Drop tables if exists (untuk development)
DROP TABLE IF EXISTS sale_items CASCADE;
DROP TABLE IF EXISTS sales CASCADE;
DROP TABLE IF EXISTS items CASCADE;
DROP TABLE IF EXISTS racks CASCADE;
DROP TABLE IF EXISTS warehouses CASCADE;
DROP TABLE IF EXISTS categories CASCADE;
DROP TABLE IF EXISTS sessions CASCADE;
DROP TABLE IF EXISTS users CASCADE;

-- =====================================================
-- TABLE: users
-- =====================================================
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(100) NOT NULL,
    role VARCHAR(20) NOT NULL CHECK (role IN ('super_admin', 'admin', 'staff')),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Index untuk performa
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_role ON users(role);

-- =====================================================
-- TABLE: sessions
-- =====================================================
CREATE TABLE sessions (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token UUID UNIQUE NOT NULL DEFAULT gen_random_uuid(),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expired_at TIMESTAMP NOT NULL,
    revoked_at TIMESTAMP NULL,
    last_activity TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    ip_address VARCHAR(50),
    user_agent TEXT
);

-- Index untuk performa
CREATE INDEX idx_sessions_token ON sessions(token);
CREATE INDEX idx_sessions_user_id ON sessions(user_id);
CREATE INDEX idx_sessions_expired_at ON sessions(expired_at);

-- =====================================================
-- TABLE: categories
-- =====================================================
CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    code VARCHAR(20) UNIQUE NOT NULL,
    name VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    is_active BOOLEAN DEFAULT true,
    created_by INTEGER REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_categories_code ON categories(code);
CREATE INDEX idx_categories_is_active ON categories(is_active);

-- =====================================================
-- TABLE: warehouses
-- =====================================================
CREATE TABLE warehouses (
    id SERIAL PRIMARY KEY,
    code VARCHAR(20) UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL,
    address TEXT,
    city VARCHAR(50),
    province VARCHAR(50),
    postal_code VARCHAR(10),
    phone VARCHAR(20),
    is_active BOOLEAN DEFAULT true,
    created_by INTEGER REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_warehouses_code ON warehouses(code);
CREATE INDEX idx_warehouses_is_active ON warehouses(is_active);

-- =====================================================
-- TABLE: racks
-- =====================================================
CREATE TABLE racks (
    id SERIAL PRIMARY KEY,
    warehouse_id INTEGER NOT NULL REFERENCES warehouses(id) ON DELETE CASCADE,
    code VARCHAR(20) UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL,
    location VARCHAR(100),
    capacity INTEGER DEFAULT 0,
    description TEXT,
    is_active BOOLEAN DEFAULT true,
    created_by INTEGER REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_racks_code ON racks(code);
CREATE INDEX idx_racks_warehouse_id ON racks(warehouse_id);
CREATE INDEX idx_racks_is_active ON racks(is_active);

-- =====================================================
-- TABLE: items
-- =====================================================
CREATE TABLE items (
    id SERIAL PRIMARY KEY,
    category_id INTEGER NOT NULL REFERENCES categories(id) ON DELETE RESTRICT,
    rack_id INTEGER REFERENCES racks(id) ON DELETE SET NULL,
    sku VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(200) NOT NULL,
    description TEXT,
    unit VARCHAR(20) NOT NULL DEFAULT 'pcs',
    price DECIMAL(15,2) NOT NULL CHECK (price >= 0),
    cost DECIMAL(15,2) DEFAULT 0 CHECK (cost >= 0),
    stock INTEGER NOT NULL DEFAULT 0 CHECK (stock >= 0),
    minimum_stock INTEGER NOT NULL DEFAULT 5,
    weight DECIMAL(10,2) DEFAULT 0,
    dimensions VARCHAR(50),
    is_active BOOLEAN DEFAULT true,
    created_by INTEGER REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_items_sku ON items(sku);
CREATE INDEX idx_items_category_id ON items(category_id);
CREATE INDEX idx_items_rack_id ON items(rack_id);
CREATE INDEX idx_items_is_active ON items(is_active);
CREATE INDEX idx_items_stock ON items(stock);
CREATE INDEX idx_items_minimum_stock ON items(stock, minimum_stock);

-- =====================================================
-- TABLE: sales
-- =====================================================
CREATE TABLE sales (
    id SERIAL PRIMARY KEY,
    invoice_number VARCHAR(50) UNIQUE NOT NULL,
    customer_name VARCHAR(100),
    customer_phone VARCHAR(20),
    customer_email VARCHAR(100),
    sale_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    total_amount DECIMAL(15,2) NOT NULL DEFAULT 0 CHECK (total_amount >= 0),
    discount DECIMAL(15,2) DEFAULT 0 CHECK (discount >= 0),
    tax DECIMAL(15,2) DEFAULT 0 CHECK (tax >= 0),
    grand_total DECIMAL(15,2) NOT NULL DEFAULT 0 CHECK (grand_total >= 0),
    payment_method VARCHAR(50),
    payment_status VARCHAR(20) DEFAULT 'pending' CHECK (payment_status IN ('pending', 'paid', 'cancelled')),
    notes TEXT,
    created_by INTEGER NOT NULL REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_sales_invoice_number ON sales(invoice_number);
CREATE INDEX idx_sales_sale_date ON sales(sale_date);
CREATE INDEX idx_sales_created_by ON sales(created_by);
CREATE INDEX idx_sales_payment_status ON sales(payment_status);

-- =====================================================
-- TABLE: sale_items
-- =====================================================
CREATE TABLE sale_items (
    id SERIAL PRIMARY KEY,
    sale_id INTEGER NOT NULL REFERENCES sales(id) ON DELETE CASCADE,
    item_id INTEGER NOT NULL REFERENCES items(id) ON DELETE RESTRICT,
    quantity INTEGER NOT NULL CHECK (quantity > 0),
    unit_price DECIMAL(15,2) NOT NULL CHECK (unit_price >= 0),
    subtotal DECIMAL(15,2) NOT NULL CHECK (subtotal >= 0),
    discount DECIMAL(15,2) DEFAULT 0 CHECK (discount >= 0),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_sale_items_sale_id ON sale_items(sale_id);
CREATE INDEX idx_sale_items_item_id ON sale_items(item_id);

-- =====================================================
-- TRIGGERS
-- =====================================================

-- Trigger untuk update updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_categories_updated_at BEFORE UPDATE ON categories
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_warehouses_updated_at BEFORE UPDATE ON warehouses
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_racks_updated_at BEFORE UPDATE ON racks
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_items_updated_at BEFORE UPDATE ON items
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_sales_updated_at BEFORE UPDATE ON sales
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- =====================================================
-- VIEWS untuk reporting
-- =====================================================

-- View untuk laporan barang dengan stok minimum
CREATE OR REPLACE VIEW v_low_stock_items AS
SELECT 
    i.id,
    i.sku,
    i.name,
    c.name as category_name,
    r.name as rack_name,
    w.name as warehouse_name,
    i.stock,
    i.minimum_stock,
    i.price,
    (i.minimum_stock - i.stock) as stock_shortage
FROM items i
LEFT JOIN categories c ON i.category_id = c.id
LEFT JOIN racks r ON i.rack_id = r.id
LEFT JOIN warehouses w ON r.warehouse_id = w.id
WHERE i.stock < i.minimum_stock AND i.is_active = true
ORDER BY (i.minimum_stock - i.stock) DESC;

-- View untuk laporan penjualan
CREATE OR REPLACE VIEW v_sales_report AS
SELECT 
    s.id,
    s.invoice_number,
    s.sale_date,
    s.customer_name,
    s.grand_total,
    s.payment_status,
    u.full_name as created_by_name,
    COUNT(si.id) as total_items
FROM sales s
LEFT JOIN sale_items si ON s.id = si.sale_id
LEFT JOIN users u ON s.created_by = u.id
GROUP BY s.id, s.invoice_number, s.sale_date, s.customer_name, 
         s.grand_total, s.payment_status, u.full_name;

-- =====================================================
-- COMMENTS
-- =====================================================
COMMENT ON TABLE users IS 'Tabel untuk menyimpan data pengguna sistem';
COMMENT ON TABLE sessions IS 'Tabel untuk menyimpan session token pengguna';
COMMENT ON TABLE categories IS 'Tabel untuk kategori barang';
COMMENT ON TABLE warehouses IS 'Tabel untuk gudang penyimpanan';
COMMENT ON TABLE racks IS 'Tabel untuk rak penyimpanan di gudang';
COMMENT ON TABLE items IS 'Tabel untuk barang/produk';
COMMENT ON TABLE sales IS 'Tabel untuk transaksi penjualan';
COMMENT ON TABLE sale_items IS 'Tabel untuk detail item penjualan';