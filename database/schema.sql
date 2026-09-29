-- ====================================================================
-- DATABASE SCHEMA: ERP PT. ANUGRAH MAHA TUNGGAL
-- Sistem Manajemen Sewa Forklift, Surat Jalan, dan Operasional
-- ====================================================================

CREATE DATABASE IF NOT EXISTS erp_amt_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE erp_amt_db;

-- 1. TABEL PENGGUNA & HAK AKSES (RBAC)
CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(100) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    role ENUM('direktur', 'admin_staff', 'operasional', 'finance') NOT NULL DEFAULT 'admin_staff',
    phone VARCHAR(20),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB;

-- 2. TABEL PELANGGAN (CUSTOMERS)
CREATE TABLE IF NOT EXISTS customers (
    id INT AUTO_INCREMENT PRIMARY KEY,
    customer_code VARCHAR(20) NOT NULL UNIQUE,
    name VARCHAR(150) NOT NULL,
    type ENUM('PT', 'CV', 'Perorangan', 'Instansi') NOT NULL DEFAULT 'PT',
    pic_name VARCHAR(100) NOT NULL,
    phone VARCHAR(25) NOT NULL,
    email VARCHAR(100),
    address TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB;

-- 3. TABEL MASTER UNIT FORKLIFT
CREATE TABLE IF NOT EXISTS forklift_units (
    id INT AUTO_INCREMENT PRIMARY KEY,
    unit_code VARCHAR(30) NOT NULL UNIQUE, -- Contoh: FL-01, FL-02
    brand VARCHAR(50) NOT NULL,           -- Contoh: Toyota, Komatsu, Mitsubishi, TCM
    model VARCHAR(50),
    capacity_ton DECIMAL(4, 1) NOT NULL,  -- Contoh: 2.5, 3.0, 5.0, 7.0, 10.0 Ton
    fuel_type ENUM('Diesel', 'Electric', 'Gasoline', 'LPG') NOT NULL DEFAULT 'Diesel',
    manufacture_year INT,
    hourly_rate DECIMAL(12, 2) DEFAULT 0,
    daily_rate DECIMAL(12, 2) DEFAULT 0,
    monthly_rate DECIMAL(12, 2) DEFAULT 0,
    status ENUM('AVAILABLE', 'RENTED', 'MAINTENANCE') NOT NULL DEFAULT 'AVAILABLE',
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB;

-- 4. TABEL MASTER OPERATOR (DRIVER ALAT BERAT)
CREATE TABLE IF NOT EXISTS operators (
    id INT AUTO_INCREMENT PRIMARY KEY,
    nip VARCHAR(30) UNIQUE,
    name VARCHAR(100) NOT NULL,
    phone VARCHAR(25) NOT NULL,
    sio_number VARCHAR(50),               -- Nomor Lisensi SIO (Surat Izin Operator)
    status ENUM('STANDBY', 'ON_DUTY', 'OFF') NOT NULL DEFAULT 'STANDBY',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB;

-- 5. TABEL PESANAN SEWA (RENTAL ORDERS)
CREATE TABLE IF NOT EXISTS rental_orders (
    id INT AUTO_INCREMENT PRIMARY KEY,
    order_number VARCHAR(40) NOT NULL UNIQUE, -- Contoh: ORD-202609-0001
    customer_id INT NOT NULL,
    order_date DATE NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE,
    rental_duration_type ENUM('JAM', 'HARI', 'BULAN') NOT NULL DEFAULT 'HARI',
    duration_value INT NOT NULL DEFAULT 1,
    project_location TEXT NOT NULL,          -- Alamat tujuan kerja forklift
    location_pic_name VARCHAR(100),          -- Kontak penerima di lapangan
    location_pic_phone VARCHAR(25),
    required_capacity_ton DECIMAL(4, 1) NOT NULL,
    unit_id INT NULL,                         -- Unit forklift yang dialokasikan
    agreed_price DECIMAL(12, 2) NOT NULL DEFAULT 0,
    status ENUM('PENDING', 'APPROVED', 'ACTIVE', 'COMPLETED', 'CANCELLED') NOT NULL DEFAULT 'PENDING',
    po_reference VARCHAR(50),                 -- Nomor PO dari Customer jika ada
    notes TEXT,
    created_by INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE RESTRICT,
    FOREIGN KEY (unit_id) REFERENCES forklift_units(id) ON DELETE SET NULL,
    FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL
) ENGINE=InnoDB;

-- 6. TABEL SURAT JALAN OPERASIONAL (DELIVERY LETTERS)
CREATE TABLE IF NOT EXISTS delivery_letters (
    id INT AUTO_INCREMENT PRIMARY KEY,
    letter_number VARCHAR(40) NOT NULL UNIQUE, -- Contoh: SJ-202609-0001
    rental_order_id INT NOT NULL,
    unit_id INT NOT NULL,
    operator_id INT NOT NULL,
    issue_date DATE NOT NULL,
    departure_time TIME,
    job_description TEXT,
    recipient_name VARCHAR(100),
    operational_status ENUM('ASSIGNED', 'ON_THE_WAY', 'WORKING', 'FINISHED') NOT NULL DEFAULT 'ASSIGNED',
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (rental_order_id) REFERENCES rental_orders(id) ON DELETE RESTRICT,
    FOREIGN KEY (unit_id) REFERENCES forklift_units(id) ON DELETE RESTRICT,
    FOREIGN KEY (operator_id) REFERENCES operators(id) ON DELETE RESTRICT
) ENGINE=InnoDB;

-- 7. TABEL LOG PERGERAKAN & STATUS OPERASIONAL (AUDIT TRAIL)
CREATE TABLE IF NOT EXISTS unit_movement_logs (
    id INT AUTO_INCREMENT PRIMARY KEY,
    unit_id INT NOT NULL,
    delivery_letter_id INT,
    status_from VARCHAR(50),
    status_to VARCHAR(50) NOT NULL,
    notes TEXT,
    logged_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (unit_id) REFERENCES forklift_units(id) ON DELETE CASCADE,
    FOREIGN KEY (delivery_letter_id) REFERENCES delivery_letters(id) ON DELETE SET NULL
) ENGINE=InnoDB;

-- 8. TABEL TAGIHAN / INVOICE
CREATE TABLE IF NOT EXISTS invoices (
    id INT AUTO_INCREMENT PRIMARY KEY,
    invoice_number VARCHAR(40) NOT NULL UNIQUE,
    rental_order_id INT NOT NULL,
    total_amount DECIMAL(12, 2) NOT NULL,
    paid_amount DECIMAL(12, 2) NOT NULL DEFAULT 0,
    payment_status ENUM('UNPAID', 'PARTIAL', 'PAID') NOT NULL DEFAULT 'UNPAID',
    due_date DATE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (rental_order_id) REFERENCES rental_orders(id) ON DELETE RESTRICT
) ENGINE=InnoDB;

-- ====================================================================
-- SEEDER DATA AWAL (UNTUK TESTING & DEMO PROYEK AKHIR)
-- ====================================================================

-- 1. User default untuk Setiap Role (Password default semua akun: 'pw123')
-- Bcrypt Hash: $2a$10$K2cWUtAqm19gYRsHC4s9dO4m4Tq4xhvRTUYg4yap.hpOmUUV6HBX6
-- 
-- REFERENSI AKUN TESTING:
-- ---------------------------------------------------------------------------------------------------------
-- Role           | Username       | Password | Email                | Akses Utama
-- ---------------------------------------------------------------------------------------------------------
-- direktur       | direktur       | pw123    | Ghofur@gmail.com     | Dashboard Eksekutif, Laporan, Approval Pesanan
-- admin_staff    | admin          | pw123    | dafa@gmail.com       | Manajemen Order, Customer, Master Unit, Surat Jalan
-- operasional    | operasional    | pw123    | luluk@gmail.com      | Jadwal Surat Jalan, Monitoring Unit & Operator Lapangan
-- finance        | finance        | pw123    | risal@gmail.com      | Tagihan, Invoice Sewa, Status Pembayaran & Kwitansi
-- ---------------------------------------------------------------------------------------------------------

INSERT INTO users (name, username, email, password_hash, role, phone, is_active) VALUES
-- Role: Direktur
('Ghofur (Direktur)', 'direktur', 'Ghofur@gmail.com', '$2a$10$K2cWUtAqm19gYRsHC4s9dO4m4Tq4xhvRTUYg4yap.hpOmUUV6HBX6', 'direktur', '081122334455', TRUE),

-- Role: Admin Staff
('Dafa (Admin Staff)', 'admin', 'dafa@gmail.com', '$2a$10$K2cWUtAqm19gYRsHC4s9dO4m4Tq4xhvRTUYg4yap.hpOmUUV6HBX6', 'admin_staff', '081234567890', TRUE),

-- Role: Operasional
('Luluk (Operasional)', 'operasional', 'luluk@gmail.com', '$2a$10$K2cWUtAqm19gYRsHC4s9dO4m4Tq4xhvRTUYg4yap.hpOmUUV6HBX6', 'operasional', '085611223344', TRUE),

-- Role: Finance
('Risal (Finance)', 'finance', 'risal@gmail.com', '$2a$10$K2cWUtAqm19gYRsHC4s9dO4m4Tq4xhvRTUYg4yap.hpOmUUV6HBX6', 'finance', '081377889900', TRUE)
ON DUPLICATE KEY UPDATE 
    name = VALUES(name),
    email = VALUES(email),
    role = VALUES(role),
    phone = VALUES(phone),
    is_active = VALUES(is_active);

-- 2. Master Customer
INSERT INTO customers (customer_code, name, type, pic_name, phone, email, address) VALUES
('CUST-001', 'PT. Pelabuhan Samudera Raya', 'PT', 'Pak Rahmat', '081234567890', 'purchasing@samuderaraya.co.id', 'Kawasan Industri Cilegon Kav 12, Banten'),
('CUST-002', 'PT. Mega Baja Mandiri', 'PT', 'Ibu Maya', '081398765432', 'logistik@megabaja.com', 'Jl. Raya Narogong Km 14, Bekasi'),
('CUST-003', 'CV. Makmur Jaya Abadi', 'CV', 'Pak Joko', '085711223344', 'makmurjaya@gmail.com', 'Jl. Industri Pergudangan No. 45, Tangerang')
ON DUPLICATE KEY UPDATE name=VALUES(name);

-- 3. Master Forklift Units
INSERT INTO forklift_units (unit_code, brand, model, capacity_ton, fuel_type, manufacture_year, daily_rate, status, notes) VALUES
('FL-01', 'Toyota', '8FD30', 3.0, 'Diesel', 2021, 1200000, 'AVAILABLE', 'Kondisi prima, servis rutin'),
('FL-02', 'Komatsu', 'FD50AY-10', 5.0, 'Diesel', 2020, 1800000, 'RENTED', 'Sedang bertugas di proyek PT. Pelabuhan Samudera Raya'),
('FL-03', 'Mitsubishi', 'FD70N', 7.0, 'Diesel', 2019, 2500000, 'AVAILABLE', 'Siap pakai'),
('FL-04', 'Toyota', '8FB25', 2.5, 'Electric', 2022, 1100000, 'AVAILABLE', 'Unit baterai, cocok untuk indoor gudang'),
('FL-05', 'TCM', 'FD100-2', 10.0, 'Diesel', 2018, 3800000, 'MAINTENANCE', 'Jadwal ganti oli hidrolik')
ON DUPLICATE KEY UPDATE brand=VALUES(brand);

-- 4. Master Operator
INSERT INTO operators (nip, name, phone, sio_number, status) VALUES
('OPR-01', 'Agus Prasetyo', '081299887766', 'SIO-K3-2022-0091', 'ON_DUTY'),
('OPR-02', 'Doni Kurniawan', '081344556677', 'SIO-K3-2023-0145', 'STANDBY'),
('OPR-03', 'Rian Hidayat', '085677889900', 'SIO-K3-2021-0032', 'STANDBY')
ON DUPLICATE KEY UPDATE name=VALUES(name);

-- 5. Contoh Transaksi Pesanan & Surat Jalan Aktif
INSERT INTO rental_orders (id, order_number, customer_id, order_date, start_date, end_date, rental_duration_type, duration_value, project_location, location_pic_name, location_pic_phone, required_capacity_ton, unit_id, agreed_price, status, po_reference) VALUES
(1, 'ORD-202609-0001', 1, '2026-09-24', '2026-09-24', '2026-09-26', 'HARI', 3, 'Kawasan Industri Cilegon Kav 12, Area Dermaga 3', 'Pak Herman', '081288990011', 5.0, 2, 5400000, 'ACTIVE', 'PO-PSR/IX/2026/88')
ON DUPLICATE KEY UPDATE order_number=VALUES(order_number);

INSERT INTO delivery_letters (id, letter_number, rental_order_id, unit_id, operator_id, issue_date, departure_time, job_description, recipient_name, operational_status) VALUES
(1, 'SJ-202609-0001', 1, 2, 1, '2026-09-24', '08:30:00', 'Bongkar muat kontainer baja plat dan pemindahan material dermaga', 'Pak Herman', 'WORKING')
ON DUPLICATE KEY UPDATE letter_number=VALUES(letter_number);
