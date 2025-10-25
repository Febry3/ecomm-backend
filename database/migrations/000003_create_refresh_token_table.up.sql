CREATE TABLE IF NOT EXISTS refresh_tokens (
  token_id           UUID        PRIMARY KEY,                
  user_id            BIGINT      NOT NULL,
  refresh_token_hash VARCHAR(255) NOT NULL,
  role               user_role    NOT NULL DEFAULT 'user',
  device_info        VARCHAR(255),
  ip_address         VARCHAR(45),
  revoked            BOOLEAN     NOT NULL DEFAULT FALSE,
  expired_at         TIMESTAMPTZ NOT NULL,
  created_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),

  CONSTRAINT fk_refresh_tokens_user
    FOREIGN KEY (user_id) REFERENCES users(id)
    ON UPDATE CASCADE ON DELETE CASCADE,

  CONSTRAINT ux_refresh_tokens_hash UNIQUE (refresh_token_hash)
);

CREATE INDEX IF NOT EXISTS ix_refresh_tokens_user_id    ON refresh_tokens (user_id);
CREATE INDEX IF NOT EXISTS ix_refresh_tokens_revoked    ON refresh_tokens (revoked);
CREATE INDEX IF NOT EXISTS ix_refresh_tokens_expired_at ON refresh_tokens (expired_at);
