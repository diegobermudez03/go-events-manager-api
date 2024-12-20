CREATE TABLE IF NOT EXISTS users(
    id          UUID REFERENCES users_auth(id),
    full_name   VARCHAR(100) NOT NULL,
    birth_date  TIMESTAMP   NOT NULL,
    gender      VARCHAR(10) NOT NULL CHECK(gender IN('MALE', 'FEMALE')),
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY(id)
)
