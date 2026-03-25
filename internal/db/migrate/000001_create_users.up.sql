CREATE TABLE IF NOT EXISTS users (
    id               CHAR(36)     NOT NULL DEFAULT (UUID()),
    name             VARCHAR(255) NOT NULL DEFAULT '',
    email            VARCHAR(255) NOT NULL DEFAULT '',
    password_hash    VARCHAR(255) NOT NULL DEFAULT '',
    active           TINYINT(1)   NOT NULL DEFAULT 0,
    verification_token VARCHAR(64) NOT NULL DEFAULT '',
    created_at       DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at       DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at       DATETIME     NULL,
    PRIMARY KEY (id),
    UNIQUE KEY uq_users_email (email),
    KEY idx_users_verification_token (verification_token)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
