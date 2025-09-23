CREATE TABLE IF NOT EXISTS regions
(
    region_id
    TEXT
    PRIMARY
    KEY,
    region_type
    TEXT
    NOT
    NULL,
    region_name
    TEXT
    NOT
    NULL,
    region_eng_name
    TEXT
    NOT
    NULL,
    last_update
    TIMESTAMP
    NOT
    NULL
    DEFAULT
    CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS users
(
    id
    BIGSERIAL
    PRIMARY
    KEY,
    username
    TEXT
    NOT
    NULL
    UNIQUE,
    password_hash
    TEXT
    NOT
    NULL,
    email
    TEXT
    NOT
    NULL,
    created_at
    TIMESTAMP
    NOT
    NULL
    DEFAULT
    CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS region_active_alarms
(
    id
    BIGSERIAL
    PRIMARY
    KEY,
    region_id
    TEXT
    NOT
    NULL
    REFERENCES
    regions
(
    region_id
) ON DELETE CASCADE,
    region_type TEXT NOT NULL,
    type TEXT NOT NULL,
    last_update TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE
(
    region_id,
    type
)
    );

CREATE INDEX IF NOT EXISTS idx_region_active_alarms_region_id
    ON region_active_alarms(region_id);
