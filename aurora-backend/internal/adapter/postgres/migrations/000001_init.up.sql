-- AURORA VPN Panel: Initial Schema

-- ─── ADMIN / AUTH ───
CREATE TABLE admins (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username      VARCHAR(64) NOT NULL UNIQUE,
    password_hash VARCHAR(256) NOT NULL,
    role          VARCHAR(16) NOT NULL DEFAULT 'admin',
    failed_attempts INT NOT NULL DEFAULT 0,
    locked_until  TIMESTAMPTZ,
    last_login_at TIMESTAMPTZ,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- ─── NODES ───
CREATE TABLE nodes (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name          VARCHAR(128) NOT NULL,
    host          VARCHAR(256) NOT NULL,
    port          INT NOT NULL DEFAULT 443,
    api_port      INT NOT NULL DEFAULT 10085,
    api_key       TEXT NOT NULL,
    status        VARCHAR(16) NOT NULL DEFAULT 'offline',
    version       VARCHAR(32),
    location      VARCHAR(128),
    weight        INT NOT NULL DEFAULT 1,
    enabled       BOOLEAN NOT NULL DEFAULT true,
    last_ping_at  TIMESTAMPTZ,
    metrics       JSONB DEFAULT '{}',
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- ─── INBOUNDS ───
CREATE TABLE inbounds (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    node_id          UUID NOT NULL REFERENCES nodes(id) ON DELETE CASCADE,
    tag              VARCHAR(128) NOT NULL,
    protocol         VARCHAR(32) NOT NULL,
    port             INT NOT NULL,
    listen           VARCHAR(64) NOT NULL DEFAULT '0.0.0.0',
    transport        VARCHAR(16) NOT NULL DEFAULT 'tcp',
    security         VARCHAR(16) NOT NULL DEFAULT 'none',
    enable           BOOLEAN NOT NULL DEFAULT true,
    settings         JSONB NOT NULL DEFAULT '{}',
    stream_settings  JSONB NOT NULL DEFAULT '{}',
    sniffing         JSONB NOT NULL DEFAULT '{}',
    user_count       INT NOT NULL DEFAULT 0,
    upload_bytes     BIGINT NOT NULL DEFAULT 0,
    download_bytes   BIGINT NOT NULL DEFAULT 0,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(node_id, tag)
);

-- ─── USERS ───
CREATE TABLE users (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username            VARCHAR(128) NOT NULL,
    email               VARCHAR(256) NOT NULL UNIQUE,
    status              VARCHAR(16) NOT NULL DEFAULT 'active',
    traffic_limit_bytes BIGINT NOT NULL DEFAULT 0,
    traffic_used_bytes  BIGINT NOT NULL DEFAULT 0,
    expire_at           TIMESTAMPTZ,
    max_ips             INT NOT NULL DEFAULT 1,
    concurrent_ips      INT NOT NULL DEFAULT 0,
    subscription_token  VARCHAR(64) UNIQUE,
    notes               TEXT,
    last_seen_at        TIMESTAMPTZ,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- ─── USER ↔ INBOUND (M:N) ───
CREATE TABLE user_inbounds (
    user_id    UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    inbound_id UUID NOT NULL REFERENCES inbounds(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, inbound_id)
);

-- ─── USER PROTOCOLS ───
CREATE TABLE user_protocols (
    user_id   UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    protocol  VARCHAR(32) NOT NULL,
    PRIMARY KEY (user_id, protocol)
);

-- ─── CONNECTION LOGS (partitioned) ───
CREATE TABLE connection_logs (
    id             BIGSERIAL,
    user_id        UUID,
    email          VARCHAR(256) NOT NULL,
    inbound_id     UUID,
    node_id        UUID,
    ip_address     INET NOT NULL,
    user_agent     TEXT,
    connected_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    disconnected_at TIMESTAMPTZ,
    upload_bytes   BIGINT NOT NULL DEFAULT 0,
    download_bytes BIGINT NOT NULL DEFAULT 0,
    PRIMARY KEY (id, connected_at)
) PARTITION BY RANGE (connected_at);

CREATE TABLE connection_logs_2026_06 PARTITION OF connection_logs
    FOR VALUES FROM ('2026-06-01') TO ('2026-07-01');
CREATE TABLE connection_logs_2026_07 PARTITION OF connection_logs
    FOR VALUES FROM ('2026-07-01') TO ('2026-08-01');
CREATE TABLE connection_logs_2026_08 PARTITION OF connection_logs
    FOR VALUES FROM ('2026-08-01') TO ('2026-09-01');

-- ─── INBOUND TRAFFIC (partitioned) ───
CREATE TABLE inbound_traffic (
    id            BIGSERIAL,
    inbound_id    UUID NOT NULL,
    user_email    VARCHAR(256) NOT NULL,
    upload        BIGINT NOT NULL DEFAULT 0,
    download      BIGINT NOT NULL DEFAULT 0,
    collected_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (id, collected_at)
) PARTITION BY RANGE (collected_at);

CREATE TABLE inbound_traffic_2026_06 PARTITION OF inbound_traffic
    FOR VALUES FROM ('2026-06-01') TO ('2026-07-01');
CREATE TABLE inbound_traffic_2026_07 PARTITION OF inbound_traffic
    FOR VALUES FROM ('2026-07-01') TO ('2026-08-01');
CREATE TABLE inbound_traffic_2026_08 PARTITION OF inbound_traffic
    FOR VALUES FROM ('2026-08-01') TO ('2026-09-01');

-- ─── SUBSCRIPTIONS ───
CREATE TABLE subscriptions (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token           VARCHAR(64) NOT NULL UNIQUE,
    url             VARCHAR(512) NOT NULL,
    format          VARCHAR(16) NOT NULL DEFAULT 'base64',
    enabled         BOOLEAN NOT NULL DEFAULT true,
    last_request_at TIMESTAMPTZ,
    last_user_agent TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- ─── SETTINGS ───
CREATE TABLE settings (
    key         VARCHAR(128) PRIMARY KEY,
    value       TEXT NOT NULL,
    description TEXT,
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- ─── FAIL2BAN / BANNED IPs ───
CREATE TABLE banned_ips (
    ip_address  INET PRIMARY KEY,
    attempts    INT NOT NULL DEFAULT 0,
    banned_until TIMESTAMPTZ NOT NULL,
    reason      VARCHAR(64) NOT NULL DEFAULT 'auth_failure',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- ─── INDEXES ───
CREATE INDEX idx_users_status ON users(status);
CREATE INDEX idx_users_expire ON users(expire_at) WHERE expire_at IS NOT NULL;
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_subscription_token ON users(subscription_token);
CREATE INDEX idx_inbounds_node ON inbounds(node_id);
CREATE INDEX idx_inbounds_protocol ON inbounds(protocol);
CREATE INDEX idx_nodes_status ON nodes(status);
CREATE INDEX idx_banned_banned_until ON banned_ips(banned_until);
CREATE INDEX idx_conn_logs_user ON connection_logs(user_id, connected_at DESC);
CREATE INDEX idx_conn_logs_email ON connection_logs(email, connected_at DESC);
CREATE INDEX idx_inbound_traffic_user ON inbound_traffic(user_email, collected_at DESC);
CREATE INDEX idx_subscriptions_token ON subscriptions(token);
CREATE INDEX idx_subscriptions_user ON subscriptions(user_id);

-- ─── DEFAULT DATA ───
INSERT INTO settings (key, value, description) VALUES
    ('panel_name', '"AURORA VPN Panel"', 'Panel display name'),
    ('public_domain', '"aurora.example.com"', 'Public domain for subscriptions'),
    ('log_retention_days', '30', 'Days to keep connection logs'),
    ('fail2ban_enabled', 'true', 'Enable Fail2Ban protection'),
    ('fail2ban_max_attempts', '5', 'Max failed login attempts before ban'),
    ('fail2ban_ban_minutes', '15', 'Ban duration in minutes');
