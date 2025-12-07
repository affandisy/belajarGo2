CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE bg_users (
    id VARCHAR(40) PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(100) NOT NULL,
    fullname VARCHAR(100) NOT NULL,
    role VARCHAR(20) NOT NULL CHECK (role IN ('superadmin', 'admin', 'user')),
    is_email_verified BOOLEAN DEFAULT FALSE
)