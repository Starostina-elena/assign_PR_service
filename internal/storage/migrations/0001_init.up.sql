CREATE TABLE IF NOT EXISTS teams (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(256) NOT NULL,
    team_id BIGINT REFERENCES teams(id) ON DELETE SET NULL,
    is_active BOOLEAN NOT NULL DEFAULT true
);

CREATE TABLE IF NOT EXISTS pull_requests (
    id SERIAL PRIMARY KEY,
    title VARCHAR(256) NOT NULL,
    author_id BIGINT REFERENCES users(id),
    is_opened BOOLEAN NOT NULL DEFAULT true,
    reviewer1_id BIGINT REFERENCES users(id),
    reviewer2_id BIGINT REFERENCES users(id)
);
