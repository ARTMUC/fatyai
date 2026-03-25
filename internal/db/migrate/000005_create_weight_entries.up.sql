CREATE TABLE IF NOT EXISTS weight_entries (
    id          CHAR(36) NOT NULL DEFAULT (UUID()),
    user_id     CHAR(36) NOT NULL,
    weight_kg   DOUBLE   NOT NULL,
    measured_at DATE     NOT NULL,
    created_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at  DATETIME NULL,
    PRIMARY KEY (id),
    UNIQUE KEY uq_weight_user_date (user_id, measured_at),
    KEY idx_weight_user_id_id (user_id, id),
    KEY idx_weight_user_measured_at (user_id, measured_at),
    CONSTRAINT fk_weight_entries_user_id FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
