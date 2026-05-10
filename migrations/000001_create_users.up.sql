CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE user_role   AS ENUM ('buyer', 'seller', 'agent', 'admin');
CREATE TYPE user_status AS ENUM ('active', 'blocked', 'pending');

CREATE TABLE users (
                       id            UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                       email         VARCHAR(255) NOT NULL UNIQUE,
                       phone         VARCHAR(20),
                       password_hash VARCHAR(255),
                       role          user_role   NOT NULL DEFAULT 'buyer',
                       status        user_status NOT NULL DEFAULT 'pending',
                       created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                       updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE user_profiles (
                               id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                               user_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                               first_name  VARCHAR(100),
                               last_name   VARCHAR(100),
                               avatar_url  TEXT,
                               bio         TEXT,
                               verified_at TIMESTAMPTZ,
                               UNIQUE(user_id)
);

CREATE TABLE oauth_accounts (
                                id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                                user_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                                provider    VARCHAR(50)  NOT NULL,
                                provider_id VARCHAR(255) NOT NULL,
                                access_token TEXT,
                                expires_at  TIMESTAMPTZ,
                                UNIQUE(provider, provider_id)
);

CREATE TABLE refresh_tokens (
                                id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                                user_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                                token_hash  VARCHAR(255) NOT NULL UNIQUE,
                                device_info VARCHAR(255),
                                expires_at  TIMESTAMPTZ NOT NULL,
                                created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_users_email        ON users(email);
CREATE INDEX idx_refresh_token_hash ON refresh_tokens(token_hash);
CREATE INDEX idx_refresh_token_user ON refresh_tokens(user_id);