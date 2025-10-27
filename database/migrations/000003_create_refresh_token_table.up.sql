CREATE TABLE IF NOT EXISTS refresh_tokens (
    token_id           UUID        PRIMARY KEY,
    user_id            BIGINT      NOT NULL,
    token_hash         VARCHAR(255) NOT NULL,
    role               user_role    NOT NULL DEFAULT 'user',
    device_info        VARCHAR(255),
    ip_address         VARCHAR(45),
    is_revoked            BOOLEAN     NOT NULL DEFAULT FALSE,
    expires_at         TIMESTAMPTZ NOT NULL,
    created_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_tokens_user
    FOREIGN KEY (user_id) REFERENCES users(id)
    ON UPDATE CASCADE ON DELETE CASCADE,

    CONSTRAINT ux_tokens_hash UNIQUE (token_hash),
    CONSTRAINT ux_tokens_user_device UNIQUE (user_id, device_info)
    );

CREATE INDEX IF NOT EXISTS ix_refresh_tokens_user_id    ON refresh_tokens (user_id);
CREATE INDEX IF NOT EXISTS ix_refresh_tokens_revoked    ON refresh_tokens (is_revoked);
CREATE INDEX IF NOT EXISTS ix_refresh_tokens_expires_at ON refresh_tokens (expires_at);