CREATE TABLE IF NOT EXISTS activities (
    id              CHAR(36)     NOT NULL DEFAULT (UUID()),
    user_id         CHAR(36)     NOT NULL,
    activity_type   VARCHAR(100) NOT NULL DEFAULT '',
    duration_min    INT          NOT NULL DEFAULT 0,
    intensity       VARCHAR(20)  NOT NULL DEFAULT '',
    calories_burned DOUBLE       NOT NULL DEFAULT 0,
    logged_at       DATETIME     NOT NULL,
    created_at      DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at      DATETIME     NULL,
    PRIMARY KEY (id),
    KEY idx_activities_user_id_id (user_id, id),
    KEY idx_activities_user_logged_at (user_id, logged_at),
    CONSTRAINT fk_activities_user_id FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
