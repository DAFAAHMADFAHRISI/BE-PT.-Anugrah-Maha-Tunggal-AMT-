-- ====================================================================================
-- RANCANGAN BASIS DATA (DATABASE SCHEMA): SISTEM ERP PT. ANUGRAH MAHA TUNGGAL
-- Sistem Terintegrasi: Persewaan Forklift, Operasional Surat Jalan, dan Billing
-- ====================================================================================
-- 
-- PANDUAN TIPE DATA (UNTUK KLIEN & LAPORAN TUGAS AKHIR):
-- 1. INT           -> Angka Bulat (Nomor Urut / ID Referensi Otomatis)
-- 2. VARCHAR(n)    -> Teks Singkat (Panjang maksimal n karakter huruf/angka)
-- 3. TEXT          -> Teks Panjang (Alamat lengkap, rincian pekerjaan, atau catatan)
-- 4. DECIMAL(p, s) -> Angka Desimal / Nilai Mata Uang Rupiah (Akurat tanpa pembulatan)
-- 5. DATE          -> Format Tanggal (Tahun-Bulan-Hari)
-- 6. TIMESTAMP     -> Waktu Otomatis Sistem (Tanggal, Jam, Menit, Detik)
-- 7. ENUM(...)     -> Pilihan Terbatas (Hanya bisa memilih opsi yang telah ditentukan)
-- 8. BOOLEAN      -> Benar / Salah (Aktif = TRUE, Nonaktif = FALSE)
-- ====================================================================================

CREATE DATABASE IF NOT EXISTS erp_amt_db 
CHARACTER SET utf8mb4 
COLLATE utf8mb4_unicode_ci;

USE erp_amt_db;

-- ------------------------------------------------------------------------------------
-- 1. TABEL PENGGUNA (users)
-- Fungsi: Menyimpan data akun login dan hak akses pengguna sistem
-- ------------------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT PRIMARY KEY 
        COMMENT 'Nomor ID Pengguna (Angka Bulat Otomatis / Kunci Utama)',
    name VARCHAR(100) NOT NULL 
        COMMENT 'Nama Lengkap Pengguna (Teks Singkat, Maks. 100 Huruf)',
    username VARCHAR(50) NOT NULL UNIQUE 
        COMMENT 'Nama Akun untuk Login (Teks Unik, Maks. 50 Huruf)',
    email VARCHAR(100) NOT NULL UNIQUE 
        COMMENT 'Alamat Surat Elektronik / Email (Teks Unik, Maks. 100 Huruf)',
    password_hash VARCHAR(255) NOT NULL 
        COMMENT 'Kata Sandi Terenkripsi Bcrypt (Kode Pengaman Sistem)',
    role ENUM('direktur', 'admin_staff', 'operasional', 'finance') NOT NULL DEFAULT 'admin_staff' 
        COMMENT 'Peran Pengguna: direktur (Pimpinan), admin_staff (Admin & Keuangan), operasional (Koordinator Lapangan), finance (Kasir/Keuangan)',
    phone VARCHAR(20) 
        COMMENT 'Nomor Telepon / WhatsApp Pengguna (Teks Angka, Maks. 20 Digit)',
    is_active BOOLEAN DEFAULT TRUE 
        COMMENT 'Status Akun: TRUE (Aktif Bisa Login), FALSE (Diblokir)',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP 
        COMMENT 'Waktu Akun Didaftarkan ke Sistem',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP 
        COMMENT 'Waktu Terakhir Informasi Akun Diperbarui'
) ENGINE=InnoDB COMMENT='Tabel Master Akun Pengguna dan Hak Akses Sistem';

-- ------------------------------------------------------------------------------------
-- 2. TABEL PELANGGAN (customers)
-- Fungsi: Menyimpan data rekanan/perusahaan/klien yang menyewa alat berat
-- ------------------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS customers (
    id INT AUTO_INCREMENT PRIMARY KEY 
        COMMENT 'Nomor ID Klien (Angka Bulat Otomatis / Kunci Utama)',
    customer_code VARCHAR(20) NOT NULL UNIQUE 
        COMMENT 'Kode Unik Klien (Contoh: CUST-001)',
    name VARCHAR(150) NOT NULL 
        COMMENT 'Nama Lengkap Perusahaan / Klien (Teks, Maks. 150 Huruf)',
    type ENUM('PT', 'CV', 'Perorangan', 'Instansi') NOT NULL DEFAULT 'PT' 
        COMMENT 'Bentuk Badan Usaha: PT, CV, Perorangan, atau Instansi Pemerintah',
    pic_name VARCHAR(100) NOT NULL 
        COMMENT 'Nama Penanggung Jawab / PIC Klien (Teks, Maks. 100 Huruf)',
    phone VARCHAR(25) NOT NULL 
        COMMENT 'Nomor Kontak PIC / Telepon Kantor (Teks, Maks. 25 Digit)',
    email VARCHAR(100) 
        COMMENT 'Email Resmi Perusahaan Klien (Teks, Maks. 100 Huruf)',
    address TEXT NOT NULL 
        COMMENT 'Alamat Lengkap Domisili / Kantor Klien (Teks Panjang)',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP 
        COMMENT 'Waktu Data Klien Ditambahkan',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP 
        COMMENT 'Waktu Terakhir Data Klien Diperbarui'
) ENGINE=InnoDB COMMENT='Tabel Master Rekanan / Pelanggan Penyewa';

-- ------------------------------------------------------------------------------------
-- 3. TABEL ARMADA FORKLIFT (forklift_units)
-- Fungsi: Menyimpan aset armada alat berat forklift milik PT. Anugrah Maha Tunggal
-- ------------------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS forklift_units (
    id INT AUTO_INCREMENT PRIMARY KEY 
        COMMENT 'Nomor ID Unit (Angka Bulat Otomatis / Kunci Utama)',
    unit_code VARCHAR(30) NOT NULL UNIQUE 
        COMMENT 'Nomor Lambung / Kode Unik Unit (Contoh: FL-01, FL-02)',
    brand VARCHAR(50) NOT NULL 
        COMMENT 'Merk Pabrikan Forklift (Contoh: Toyota, Komatsu, Mitsubishi, TCM)',
    model VARCHAR(50) 
        COMMENT 'Tipe / Seri Model Forklift (Contoh: 8FD30, FD50AY-10)',
    capacity_ton DECIMAL(4, 1) NOT NULL 
        COMMENT 'Kapasitas Angkat dalam Satuan Ton (Contoh: 2.5, 3.0, 5.0, 10.0 Ton)',
    fuel_type ENUM('Diesel', 'Electric', 'Gasoline', 'LPG') NOT NULL DEFAULT 'Diesel' 
        COMMENT 'Jenis Bahan Bakar: Solar (Diesel), Baterai Listrik (Electric), Bensin (Gasoline), Gas (LPG)',
    manufacture_year INT 
        COMMENT 'Tahun Pembuatan Unit (Angka Bulat Tahun, Contoh: 2021)',
    hourly_rate DECIMAL(12, 2) DEFAULT 0 
        COMMENT 'Tarif Sewa Per Jam dalam Rupiah (Angka Desimal Nilai Uang)',
    daily_rate DECIMAL(12, 2) DEFAULT 0 
        COMMENT 'Tarif Sewa Per Hari dalam Rupiah (Angka Desimal Nilai Uang)',
    monthly_rate DECIMAL(12, 2) DEFAULT 0 
        COMMENT 'Tarif Sewa Per Bulan dalam Rupiah (Angka Desimal Nilai Uang)',
    status ENUM('AVAILABLE', 'RENTED', 'MAINTENANCE') NOT NULL DEFAULT 'AVAILABLE' 
        COMMENT 'Status Kondisi Unit: AVAILABLE (Siap Pakai di Pool), RENTED (Sedang Bertugas/Disewa), MAINTENANCE (Sedang Rusak/Servis)',
    notes TEXT 
        COMMENT 'Catatan Khusus Kondisi Mesin / Kelengkapan Unit',
    hour_meter DECIMAL(10, 2) NOT NULL DEFAULT 0 
        COMMENT 'C1: Total Jam Kerja Operasional / Hour Meter (Satuan: Jam)',
    machine_age DECIMAL(5, 2) NOT NULL DEFAULT 0 
        COMMENT 'C2: Usia Operasional Mesin (Satuan: Tahun)',
    breakdown_frequency DECIMAL(5, 2) NOT NULL DEFAULT 0 
        COMMENT 'C3: Frekuensi Kerusakan 3 Bulan Terakhir (Satuan: Kali)',
    idle_days DECIMAL(5, 2) NOT NULL DEFAULT 0 
        COMMENT 'C4: Sisa Hari Bebas Tugas / Idle Time di Pool (Satuan: Hari)',
    last_calculated_score DECIMAL(8, 4) NOT NULL DEFAULT 0 
        COMMENT 'Nilai Preferensi Vi SAW Terkini (Rentang 0 s.d 1)',
    last_risk_status VARCHAR(20) NOT NULL DEFAULT 'SAFE' 
        COMMENT 'Status Kategori Risiko Terkini: CRITICAL, WARNING, SAFE',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP 
        COMMENT 'Waktu Pendaftaran Unit ke Sistem',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP 
        COMMENT 'Waktu Terakhir Status Unit Diperbarui'
) ENGINE=InnoDB COMMENT='Tabel Master Aset Alat Berat Forklift';

-- ------------------------------------------------------------------------------------
-- 4. TABEL OPERATOR ALAT BERAT (operators)
-- Fungsi: Menyimpan data pengemudi / teknisi operator forklift berlisensi K3
-- ------------------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS operators (
    id INT AUTO_INCREMENT PRIMARY KEY 
        COMMENT 'Nomor ID Operator (Angka Bulat Otomatis / Kunci Utama)',
    nip VARCHAR(30) UNIQUE 
        COMMENT 'Nomor Induk Pegawai / Karyawan (Teks Unik, Maks. 30 Karakter)',
    name VARCHAR(100) NOT NULL 
        COMMENT 'Nama Lengkap Operator (Teks Singkat, Maks. 100 Huruf)',
    phone VARCHAR(25) NOT NULL 
        COMMENT 'Nomor Telepon / HP Operator (Teks Angka, Maks. 25 Digit)',
    sio_number VARCHAR(50) 
        COMMENT 'Nomor Lisensi SIO (Surat Izin Operator K3 Kemnaker RI)',
    status ENUM('STANDBY', 'ON_DUTY', 'OFF') NOT NULL DEFAULT 'STANDBY' 
        COMMENT 'Status Tugas Operator: STANDBY (Siap Berangkat), ON_DUTY (Sedang Bertugas di Proyek), OFF (Libur/Cuti)',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP 
        COMMENT 'Waktu Pendaftaran Operator',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP 
        COMMENT 'Waktu Terakhir Data Operator Diperbarui'
) ENGINE=InnoDB COMMENT='Tabel Master Operator Forklift Berlisensi SIO';

-- ------------------------------------------------------------------------------------
-- 5. TABEL PESANAN SEWA (rental_orders)
-- Fungsi: Mencatat order permintaan sewa dari pelanggan (PO Sewa)
-- ------------------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS rental_orders (
    id INT AUTO_INCREMENT PRIMARY KEY 
        COMMENT 'Nomor ID Pesanan (Angka Bulat Otomatis / Kunci Utama)',
    order_number VARCHAR(40) NOT NULL UNIQUE 
        COMMENT 'Nomor Unik Kontrak Pesanan (Contoh: ORD-202609-0001)',
    customer_id INT NOT NULL 
        COMMENT 'Nomor ID Klien Penyewa (Terhubung ke Tabel Pelanggan)',
    order_date DATE NOT NULL 
        COMMENT 'Tanggal Pesanan Diterima (Format: Tahun-Bulan-Hari)',
    start_date DATE NOT NULL 
        COMMENT 'Tanggal Mulai Pemakaian Sewa (Format: Tahun-Bulan-Hari)',
    end_date DATE 
        COMMENT 'Tanggal Selesai Pemakaian Sewa (Format: Tahun-Bulan-Hari)',
    rental_duration_type ENUM('JAM', 'HARI', 'BULAN') NOT NULL DEFAULT 'HARI' 
        COMMENT 'Skema Satuan Sewa: JAM (Shift Singkat), HARI (Harian), BULAN (Bulanan Kontrak)',
    duration_value INT NOT NULL DEFAULT 1 
        COMMENT 'Jumlah Satuan Durasi Sewa (Angka Bulat, Contoh: 3 Hari)',
    project_location TEXT NOT NULL 
        COMMENT 'Alamat Lengkap Lokasi Proyek Tempat Forklift Bekerja',
    location_pic_name VARCHAR(100) 
        COMMENT 'Nama Penanggung Jawab di Lapangan Proyek',
    location_pic_phone VARCHAR(25) 
        COMMENT 'Nomor Telepon PIC Lapangan Proyek',
    required_capacity_ton DECIMAL(4, 1) NOT NULL 
        COMMENT 'Kapasitas Tonase Forklift yang Diminta Klien (Contoh: 5.0 Ton)',
    unit_id INT NULL 
        COMMENT 'ID Unit Forklift yang Dialokasikan (Terhubung ke Tabel Unit)',
    agreed_price DECIMAL(12, 2) NOT NULL DEFAULT 0 
        COMMENT 'Kesepakatan Total Nilai Harga Sewa dalam Rupiah (Mata Uang)',
    status ENUM('PENDING', 'APPROVED', 'ACTIVE', 'COMPLETED', 'CANCELLED') NOT NULL DEFAULT 'PENDING' 
        COMMENT 'Status Alur Order: PENDING (Menunggu Persetujuan), APPROVED (Disetujui), ACTIVE (Sedang Berjalan), COMPLETED (Selesai), CANCELLED (Dibatalkan)',
    po_reference VARCHAR(50) 
        COMMENT 'Nomor Dokumen Purchase Order (PO) Fisik dari Klien',
    notes TEXT 
        COMMENT 'Instruksi Tambahan / Catatan Khusus Pesanan',
    created_by INT 
        COMMENT 'ID Pengguna / Staff yang Memasukkan Pesanan',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP 
        COMMENT 'Waktu Pencatatan Pesanan Masuk',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP 
        COMMENT 'Waktu Terakhir Perubahan Data Pesanan',
    FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE RESTRICT,
    FOREIGN KEY (unit_id) REFERENCES forklift_units(id) ON DELETE SET NULL,
    FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL
) ENGINE=InnoDB COMMENT='Tabel Transaksi Kontrak Pesanan Sewa (Purchase Order)';

-- ------------------------------------------------------------------------------------
-- 6. TABEL SURAT JALAN OPERASIONAL (delivery_letters)
-- Fungsi: Mencatat penerbitan surat jalan resmi untuk unit dan operator berangkat
-- ------------------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS delivery_letters (
    id INT AUTO_INCREMENT PRIMARY KEY 
        COMMENT 'Nomor ID Surat Jalan (Angka Bulat Otomatis / Kunci Utama)',
    letter_number VARCHAR(40) NOT NULL UNIQUE 
        COMMENT 'Nomor Registrasi Surat Jalan Resmi (Contoh: SJ-202609-0001)',
    rental_order_id INT NOT NULL 
        COMMENT 'Nomor ID Pesanan Sewa (Terhubung ke Tabel Pesanan Sewa)',
    unit_id INT NOT NULL 
        COMMENT 'Nomor ID Forklift yang Ditugaskan (Terhubung ke Tabel Unit)',
    operator_id INT NOT NULL 
        COMMENT 'Nomor ID Operator yang Bertugas (Terhubung ke Tabel Operator)',
    issue_date DATE NOT NULL 
        COMMENT 'Tanggal Dokumen Surat Jalan Diterbitkan',
    departure_time VARCHAR(50) 
        COMMENT 'Jam Berangkat Unit dari Pool (Contoh: 08:30 WIB)',
    job_description TEXT 
        COMMENT 'Uraian Instruksi Kerja (Bongkar Muat, Pindah Barang, dll)',
    recipient_name VARCHAR(100) 
        COMMENT 'Nama Pihak Penerima di Lokasi Proyek',
    operational_status ENUM('ASSIGNED', 'ON_THE_WAY', 'WORKING', 'FINISHED') NOT NULL DEFAULT 'ASSIGNED' 
        COMMENT 'Status Fisik Lapangan: ASSIGNED (Unit Ditugaskan di Pool), ON_THE_WAY (Sedang Berjalan ke Lokasi), WORKING (Sedang Bekerja di Lokasi Proyek), FINISHED (Pekerjaan Selesai Kembali ke Pool)',
    notes TEXT 
        COMMENT 'Catatan Riil Lapangan / Kondisi Selama Tugas',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP 
        COMMENT 'Waktu Penerbitan Surat Jalan oleh Sistem',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP 
        COMMENT 'Waktu Terakhir Status Lapangan Diperbarui',
    FOREIGN KEY (rental_order_id) REFERENCES rental_orders(id) ON DELETE RESTRICT,
    FOREIGN KEY (unit_id) REFERENCES forklift_units(id) ON DELETE RESTRICT,
    FOREIGN KEY (operator_id) REFERENCES operators(id) ON DELETE RESTRICT
) ENGINE=InnoDB COMMENT='Tabel Operasional Surat Jalan dan Penugasan Lapangan';

-- ------------------------------------------------------------------------------------
-- 7. TABEL LOG PERGERAKAN UNIT (unit_movement_logs)
-- Fungsi: Rekam jejak audit perubahan status unit forklift secara kronologis
-- ------------------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS unit_movement_logs (
    id INT AUTO_INCREMENT PRIMARY KEY 
        COMMENT 'Nomor ID Log Audit (Angka Bulat Otomatis / Kunci Utama)',
    unit_id INT NOT NULL 
        COMMENT 'Nomor ID Unit Forklift yang Mengalami Perubahan Status',
    delivery_letter_id INT 
        COMMENT 'Nomor ID Surat Jalan Terkait (Bila Sedang Bertugas)',
    status_from VARCHAR(50) 
        COMMENT 'Status Awal Unit Sebelum Berubah (Contoh: AVAILABLE, ASSIGNED)',
    status_to VARCHAR(50) NOT NULL 
        COMMENT 'Status Baru Unit Setelah Perubahan (Contoh: WORKING, FINISHED)',
    notes TEXT 
        COMMENT 'Keterangan Tambahan / Alasan Perubahan Status',
    logged_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP 
        COMMENT 'Waktu Tepat Tercatatnya Perubahan Status dalam Sistem',
    FOREIGN KEY (unit_id) REFERENCES forklift_units(id) ON DELETE CASCADE,
    FOREIGN KEY (delivery_letter_id) REFERENCES delivery_letters(id) ON DELETE SET NULL
) ENGINE=InnoDB COMMENT='Tabel Riwayat Audit Trail Status Lapangan';

-- ------------------------------------------------------------------------------------
-- 8. TABEL TAGIHAN / INVOICE (invoices)
-- Fungsi: Mencatat penerbitan tagihan biaya sewa kepada pelanggan dan pembayarannya
-- ------------------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS invoices (
    id INT AUTO_INCREMENT PRIMARY KEY 
        COMMENT 'Nomor ID Tagihan (Angka Bulat Otomatis / Kunci Utama)',
    invoice_number VARCHAR(40) NOT NULL UNIQUE 
        COMMENT 'Nomor Faktur / Invoice Resmi (Contoh: INV-202609-0001)',
    rental_order_id INT NOT NULL 
        COMMENT 'Nomor ID Pesanan Sewa yang Ditagihkan',
    total_amount DECIMAL(12, 2) NOT NULL 
        COMMENT 'Total Nominal Tagihan Biaya Sewa dalam Rupiah',
    paid_amount DECIMAL(12, 2) NOT NULL DEFAULT 0 
        COMMENT 'Nominal Pembayaran yang Telah Diterima dalam Rupiah',
    payment_status ENUM('UNPAID', 'PARTIAL', 'PAID') NOT NULL DEFAULT 'UNPAID' 
        COMMENT 'Status Pembayaran: UNPAID (Belum Bayar), PARTIAL (Bayar Sebagian/Uang Muka), PAID (Sudah Lunas)',
    due_date DATE 
        COMMENT 'Batas Waktu Tanggal Jatuh Tempo Pembayaran',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP 
        COMMENT 'Waktu Pembuatan Faktur Tagihan',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP 
        COMMENT 'Waktu Terakhir Pencatatan Pembayaran',
    FOREIGN KEY (rental_order_id) REFERENCES rental_orders(id) ON DELETE RESTRICT
) ENGINE=InnoDB COMMENT='Tabel Keuangan Penagihan Invoice dan Pembayaran';

-- ------------------------------------------------------------------------------------
-- 9. TABEL MASTER KRITERIA SPK (kriteria)
-- Fungsi: Menyimpan 4 parameter kriteria penentu keputusan pemeliharaan (C1 s.d C4)
-- ------------------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS kriteria (
    id INT AUTO_INCREMENT PRIMARY KEY 
        COMMENT 'Nomor ID Kriteria (Angka Bulat Otomatis / Kunci Utama)',
    code VARCHAR(10) NOT NULL UNIQUE 
        COMMENT 'Kode Kriteria SPK (Contoh: C1, C2, C3, C4)',
    name VARCHAR(100) NOT NULL 
        COMMENT 'Nama Lengkap Kriteria (Jam Kerja, Usia Mesin, Kerusakan, Hari Bebas)',
    attribute_type ENUM('benefit', 'cost') NOT NULL DEFAULT 'benefit' 
        COMMENT 'Sifat Atribut: benefit (makin besar makin prioritas), cost (makin kecil makin prioritas)',
    unit VARCHAR(20) 
        COMMENT 'Satuan Nilai Kriteria (Jam, Tahun, Kali, Hari)',
    description TEXT 
        COMMENT 'Penjelasan dan Pengaruh Kriteria Terhadap Pemeliharaan',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB COMMENT='Tabel Master 4 Kriteria Pemrosesan SPK SAW';

-- ------------------------------------------------------------------------------------
-- 10. TABEL BOBOT PREFERENSI KRITERIA (bobot_kriteria)
-- Fungsi: Menyimpan nilai bobot preferensi (W1 s.d W4) yang ditentukan manajemen
-- ------------------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS bobot_kriteria (
    id INT AUTO_INCREMENT PRIMARY KEY 
        COMMENT 'Nomor ID Bobot (Angka Bulat Otomatis / Kunci Utama)',
    kriteria_id INT NOT NULL 
        COMMENT 'ID Referensi Kriteria (Terhubung ke Tabel kriteria)',
    weight DECIMAL(6, 4) NOT NULL 
        COMMENT 'Nilai Desimal Bobot Preferensi (Wj, Contoh: 0.3000, Total Harus 1.0)',
    percentage DECIMAL(5, 2) NOT NULL DEFAULT 0 
        COMMENT 'Nilai Persentase Bobot (Contoh: 30.00%)',
    set_by INT 
        COMMENT 'ID User/Direktur yang Menetapkan Bobot',
    notes TEXT 
        COMMENT 'Catatan Pertimbangan Manajemen atas Pembobotan',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (kriteria_id) REFERENCES kriteria(id) ON DELETE CASCADE,
    FOREIGN KEY (set_by) REFERENCES users(id) ON DELETE SET NULL
) ENGINE=InnoDB COMMENT='Tabel Nilai Bobot Preferensi Manajemen (W1..W4)';

-- ------------------------------------------------------------------------------------
-- 11. TABEL LOG OPERASIONAL ARMADA (log_operasional)
-- Fungsi: Merekam variabel mentah harian forklift (Hour Meter, Usia, Kerusakan, Idle)
-- ------------------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS log_operasional (
    id INT AUTO_INCREMENT PRIMARY KEY 
        COMMENT 'Nomor ID Log Operasional (Angka Bulat Otomatis)',
    unit_id INT NOT NULL 
        COMMENT 'ID Forklift yang Dicatat (Terhubung ke Tabel forklift_units)',
    log_date DATE NOT NULL 
        COMMENT 'Tanggal Pencatatan Variabel Operasional Unit',
    hour_meter DECIMAL(10, 2) NOT NULL DEFAULT 0 
        COMMENT 'C1: Total Jam Kerja Operasional / Hour Meter (Satuan: Jam)',
    machine_age DECIMAL(5, 2) NOT NULL DEFAULT 0 
        COMMENT 'C2: Usia Operasional Mesin (Satuan: Tahun)',
    breakdown_freq DECIMAL(5, 2) NOT NULL DEFAULT 0 
        COMMENT 'C3: Frekuensi Kerusakan Historis 3 Bulan Terakhir (Satuan: Kali)',
    idle_days DECIMAL(5, 2) NOT NULL DEFAULT 0 
        COMMENT 'C4: Sisa Hari Bebas Tugas / Idle Time di Pool (Satuan: Hari)',
    recorded_by INT 
        COMMENT 'ID Staf Operasional / Admin yang Menginput',
    notes TEXT 
        COMMENT 'Catatan Fisik Mesin, Suara Aneh, atau Kebocoran Hidrolik',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (unit_id) REFERENCES forklift_units(id) ON DELETE CASCADE,
    FOREIGN KEY (recorded_by) REFERENCES users(id) ON DELETE SET NULL
) ENGINE=InnoDB COMMENT='Tabel Input Variabel Mentah Harian Armada Forklift';

-- ------------------------------------------------------------------------------------
-- 12. TABEL SNAPSHOT HASIL RANKING SAW (hasil_ranking)
-- Fungsi: Menyimpan riwayat kalkulasi ranking global per sesi hitung (Audit Trail)
-- ------------------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS hasil_ranking (
    id INT AUTO_INCREMENT PRIMARY KEY 
        COMMENT 'Nomor ID Sesi Perhitungan (Angka Bulat Otomatis)',
    calculation_code VARCHAR(50) NOT NULL UNIQUE 
        COMMENT 'Kode Unik Sesi Kalkulasi (Contoh: SAW-20261111-001)',
    calculated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP 
        COMMENT 'Waktu Komputasi SAW Dieksekusi',
    calculated_by INT 
        COMMENT 'ID Pengguna yang Memicu Kalkulasi',
    total_units INT NOT NULL DEFAULT 0 
        COMMENT 'Jumlah Armada Forklift yang Dihitung',
    critical_count INT NOT NULL DEFAULT 0 
        COMMENT 'Jumlah Unit Kategori Merah / Kritis (Servis Segera)',
    warning_count INT NOT NULL DEFAULT 0 
        COMMENT 'Jumlah Unit Kategori Kuning / Waspada',
    safe_count INT NOT NULL DEFAULT 0 
        COMMENT 'Jumlah Unit Kategori Hijau / Aman',
    notes TEXT 
        COMMENT 'Ringkasan Rekomendasi Sistem untuk Direksi',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (calculated_by) REFERENCES users(id) ON DELETE SET NULL
) ENGINE=InnoDB COMMENT='Tabel Riwayat Audit Trail Komputasi Algoritma SAW';

-- ------------------------------------------------------------------------------------
-- 13. TABEL DETAIL PERANKINGAN FORKLIFT (detail_ranking)
-- Fungsi: Menyimpan nilai normalisasi Rij, skor akhir Vi, dan ranking unit
-- ------------------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS detail_ranking (
    id INT AUTO_INCREMENT PRIMARY KEY 
        COMMENT 'Nomor ID Detail Ranking (Angka Bulat Otomatis)',
    hasil_ranking_id INT NOT NULL 
        COMMENT 'ID Sesi Perhitungan (Terhubung ke hasil_ranking)',
    unit_id INT NOT NULL 
        COMMENT 'ID Unit Forklift yang Dinilai',
    rank_order INT NOT NULL 
        COMMENT 'Nomor Urut Peringkat Prioritas Pemeliharaan (1, 2, 3...)',
    c1_raw DECIMAL(10, 2) NOT NULL DEFAULT 0 
        COMMENT 'Nilai Mentah Jam Kerja (Hour Meter)',
    c2_raw DECIMAL(5, 2) NOT NULL DEFAULT 0 
        COMMENT 'Nilai Mentah Usia Mesin (Tahun)',
    c3_raw DECIMAL(5, 2) NOT NULL DEFAULT 0 
        COMMENT 'Nilai Mentah Frekuensi Kerusakan (Kali)',
    c4_raw DECIMAL(5, 2) NOT NULL DEFAULT 0 
        COMMENT 'Nilai Mentah Hari Bebas Tugas (Hari)',
    c1_normalized DECIMAL(8, 4) NOT NULL DEFAULT 0 
        COMMENT 'Nilai Hasil Normalisasi R1 (Benefit/Cost)',
    c2_normalized DECIMAL(8, 4) NOT NULL DEFAULT 0 
        COMMENT 'Nilai Hasil Normalisasi R2 (Benefit/Cost)',
    c3_normalized DECIMAL(8, 4) NOT NULL DEFAULT 0 
        COMMENT 'Nilai Hasil Normalisasi R3 (Benefit/Cost)',
    c4_normalized DECIMAL(8, 4) NOT NULL DEFAULT 0 
        COMMENT 'Nilai Hasil Normalisasi R4 (Benefit/Cost)',
    score_v DECIMAL(8, 4) NOT NULL DEFAULT 0 
        COMMENT 'Nilai Preferensi Akhir Vi (Rentang 0.0 s.d 1.0)',
    risk_status VARCHAR(20) NOT NULL DEFAULT 'SAFE' 
        COMMENT 'Status Risiko: CRITICAL, WARNING, SAFE',
    risk_color VARCHAR(20) NOT NULL DEFAULT 'hijau' 
        COMMENT 'Indikator Warna Badge: merah, kuning, hijau',
    recommendation TEXT 
        COMMENT 'Keterangan Rekomendasi Tindakan Bengkel',
    spk_status VARCHAR(20) NOT NULL DEFAULT 'NONE' 
        COMMENT 'Status Tindak Lanjut: NONE (Belum Ada SPK), ISSUED (Sudah Diterbitkan SPK)',
    FOREIGN KEY (hasil_ranking_id) REFERENCES hasil_ranking(id) ON DELETE CASCADE,
    FOREIGN KEY (unit_id) REFERENCES forklift_units(id) ON DELETE CASCADE
) ENGINE=InnoDB COMMENT='Tabel Rincian Skor dan Normalisasi SAW Per Alternatif';

-- ------------------------------------------------------------------------------------
-- 14. TABEL SURAT PERINTAH KERJA (SPK) PEMELIHARAAN (spk_maintenance)
-- Fungsi: Menghubungkan rekomendasi ranking SAW menuju eksekusi modul bengkel ERP
-- ------------------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS spk_maintenance (
    id INT AUTO_INCREMENT PRIMARY KEY 
        COMMENT 'Nomor ID SPK Bengkel (Angka Bulat Otomatis)',
    spk_number VARCHAR(50) NOT NULL UNIQUE 
        COMMENT 'Nomor Registrasi SPK Bengkel Resmi (Contoh: SPK-MT-202611-0001)',
    unit_id INT NOT NULL 
        COMMENT 'ID Unit Forklift yang Diservis (Otomatis Jadi MAINTENANCE)',
    risk_score DECIMAL(8, 4) NOT NULL DEFAULT 0 
        COMMENT 'Skor Risiko Vi SAW Saat SPK Diterbitkan',
    risk_status VARCHAR(20) NOT NULL DEFAULT 'CRITICAL' 
        COMMENT 'Kategori Risiko Saat Diterbitkan (merah/kuning)',
    mechanic_name VARCHAR(100) NOT NULL 
        COMMENT 'Nama Lengkap Mekanik Penanggung Jawab',
    complaint_notes TEXT NOT NULL 
        COMMENT 'Catatan Keluhan / Target Komponen Perbaikan',
    scheduled_date DATE NOT NULL 
        COMMENT 'Rencana Tanggal Pelaksanaan Servis di Bengkel',
    status ENUM('DRAFT', 'IN_PROGRESS', 'COMPLETED', 'CANCELLED') NOT NULL DEFAULT 'IN_PROGRESS' 
        COMMENT 'Status Pekerjaan Servis di Bengkel',
    issued_by INT 
        COMMENT 'ID Direktur / Kepala Mekanik yang Menerbitkan SPK',
    issued_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP 
        COMMENT 'Waktu Penerbitan Dokumen SPK Fisik',
    notes TEXT 
        COMMENT 'Catatan Hasil Uji Coba Pasca-Servis',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (unit_id) REFERENCES forklift_units(id) ON DELETE RESTRICT,
    FOREIGN KEY (issued_by) REFERENCES users(id) ON DELETE SET NULL
) ENGINE=InnoDB COMMENT='Tabel Eksekusi Dokumen SPK Pemeliharaan Bengkel';

-- ====================================================================================
-- DATA AWAL / SEEDER CONTOH (UNTUK DEMO & PENGUJIAN SISTEM)
-- ====================================================================================

-- 1. Data Akun Pengguna Default (Kata sandi bawaan semua akun: 'pw123')
INSERT INTO users (name, username, email, password_hash, role, phone, is_active) VALUES
('Ghofur (Direktur)', 'direktur', 'Ghofur@gmail.com', '$2a$10$K2cWUtAqm19gYRsHC4s9dO4m4Tq4xhvRTUYg4yap.hpOmUUV6HBX6', 'direktur', '081122334455', TRUE),
('Dafa (Staff Admin & Keuangan)', 'admin', 'dafa@gmail.com', '$2a$10$K2cWUtAqm19gYRsHC4s9dO4m4Tq4xhvRTUYg4yap.hpOmUUV6HBX6', 'admin_staff', '081234567890', TRUE),
('Luluk (Koordinator Lapangan)', 'operasional', 'luluk@gmail.com', '$2a$10$K2cWUtAqm19gYRsHC4s9dO4m4Tq4xhvRTUYg4yap.hpOmUUV6HBX6', 'operasional', '085611223344', TRUE),
('Risal (Staff Keuangan)', 'finance', 'risal@gmail.com', '$2a$10$K2cWUtAqm19gYRsHC4s9dO4m4Tq4xhvRTUYg4yap.hpOmUUV6HBX6', 'finance', '081377889900', TRUE)
ON DUPLICATE KEY UPDATE 
    name = VALUES(name),
    email = VALUES(email),
    role = VALUES(role),
    phone = VALUES(phone),
    is_active = VALUES(is_active);

-- 2. Data Pelanggan / Rekanan
INSERT INTO customers (customer_code, name, type, pic_name, phone, email, address) VALUES
('CUST-001', 'PT. Pelabuhan Samudera Raya', 'PT', 'Pak Rahmat', '081234567890', 'purchasing@samuderaraya.co.id', 'Kawasan Industri Cilegon Kav 12, Banten'),
('CUST-002', 'PT. Mega Baja Mandiri', 'PT', 'Ibu Maya', '081398765432', 'logistik@megabaja.com', 'Jl. Raya Narogong Km 14, Bekasi'),
('CUST-003', 'CV. Makmur Jaya Abadi', 'CV', 'Pak Joko', '085711223344', 'makmurjaya@gmail.com', 'Jl. Industri Pergudangan No. 45, Tangerang')
ON DUPLICATE KEY UPDATE name=VALUES(name);

-- 3. Data Armada Forklift
INSERT INTO forklift_units (unit_code, brand, model, capacity_ton, fuel_type, manufacture_year, daily_rate, status, notes) VALUES
('FL-01', 'Toyota', '8FD30', 3.0, 'Diesel', 2021, 1200000, 'AVAILABLE', 'Kondisi prima, servis rutin berkala'),
('FL-02', 'Komatsu', 'FD50AY-10', 5.0, 'Diesel', 2020, 1800000, 'RENTED', 'Sedang bertugas di proyek dermaga PT. Pelabuhan Samudera Raya'),
('FL-03', 'Mitsubishi', 'FD70N', 7.0, 'Diesel', 2019, 2500000, 'AVAILABLE', 'Unit siap kerja untuk beban berat'),
('FL-04', 'Toyota', '8FB25', 2.5, 'Electric', 2022, 1100000, 'AVAILABLE', 'Unit baterai ramah lingkungan untuk pergudangan tertutup'),
('FL-05', 'TCM', 'FD100-2', 10.0, 'Diesel', 2018, 3800000, 'MAINTENANCE', 'Sedang dalam jadwal servis berkala')
ON DUPLICATE KEY UPDATE brand=VALUES(brand);

-- 4. Data Operator Berlisensi SIO
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
