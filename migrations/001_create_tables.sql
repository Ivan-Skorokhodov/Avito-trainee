CREATE TABLE teams (
    team_id     SERIAL PRIMARY KEY,
    team_name   TEXT NOT NULL UNIQUE
);

CREATE TABLE users (
    user_id   SERIAL PRIMARY KEY,
    system_id TEXT NOT NULL UNIQUE,
    user_name TEXT NOT NULL,
    team_id   NOT NULL INT REFERENCES teams(team_id) ON DELETE CASCADE,
    is_active BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE pull_requests (
    pull_request_id   SERIAL PRIMARY KEY,
    system_id         TEXT NOT NULL UNIQUE,
    pull_request_name TEXT NOT NULL,
    author_id         INT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    status            TEXT NOT NULL REFERENCES statuses(status) ON DELETE RESTRICT,
    created_at        TIMESTAMP NOT NULL DEFAULT NOW(),
    merged_at         TIMESTAMP
);

CREATE TABLE pull_request_reviewers (
    id              SERIAL PRIMARY KEY,
    pull_request_id INT NOT NULL REFERENCES pull_requests(pull_request_id) ON DELETE CASCADE,
    user_id         INT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    UNIQUE(pull_request_id, user_id)
);

CREATE TABLE statuses (
    status TEXT PRIMARY KEY
);


