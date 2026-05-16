CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE SEQUENCE IF NOT EXISTS letter_number_seq
    START WITH 1
    INCREMENT BY 1
    MINVALUE 1
    NO MAXVALUE
    CACHE 1;

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    username VARCHAR(100) NOT NULL UNIQUE,
    full_name VARCHAR(150) NOT NULL,
    password_hash TEXT NOT NULL,

    role VARCHAR(30) NOT NULL CHECK (role IN ('superuser', 'editor', 'readonly')),

    is_active BOOLEAN NOT NULL DEFAULT TRUE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_users_is_active ON users(is_active);

CREATE TABLE letters (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    letter_number BIGINT NOT NULL UNIQUE,

    title VARCHAR(255) NOT NULL,
    letter_date DATE NOT NULL,

    registrar_name VARCHAR(150) NOT NULL,
    destination VARCHAR(255) NOT NULL,

    description TEXT,

    created_by UUID NOT NULL REFERENCES users(id),
    updated_by UUID REFERENCES users(id),
    deleted_by UUID REFERENCES users(id),

    is_deleted BOOLEAN NOT NULL DEFAULT FALSE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_letters_letter_number ON letters(letter_number);
CREATE INDEX idx_letters_title ON letters(title);
CREATE INDEX idx_letters_letter_date ON letters(letter_date);
CREATE INDEX idx_letters_destination ON letters(destination);
CREATE INDEX idx_letters_registrar_name ON letters(registrar_name);
CREATE INDEX idx_letters_is_deleted ON letters(is_deleted);
CREATE INDEX idx_letters_created_by ON letters(created_by);
CREATE INDEX idx_letters_created_at ON letters(created_at);

CREATE TABLE audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    actor_user_id UUID REFERENCES users(id),

    action VARCHAR(100) NOT NULL,
    entity_type VARCHAR(100) NOT NULL,
    entity_id UUID,

    old_value JSONB,
    new_value JSONB,

    ip_address INET,
    user_agent TEXT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_audit_logs_actor_user_id ON audit_logs(actor_user_id);
CREATE INDEX idx_audit_logs_action ON audit_logs(action);
CREATE INDEX idx_audit_logs_entity_type ON audit_logs(entity_type);
CREATE INDEX idx_audit_logs_entity_id ON audit_logs(entity_id);
CREATE INDEX idx_audit_logs_created_at ON audit_logs(created_at);

CREATE TABLE app_settings (
    key VARCHAR(100) PRIMARY KEY,
    value JSONB NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_app_settings_updated_at ON app_settings(updated_at);