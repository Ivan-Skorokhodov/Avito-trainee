CREATE TABLE teams (
    team_id     SERIAL PRIMARY KEY,
    team_name   TEXT NOT NULL UNIQUE
);

CREATE TABLE users (
    user_id   SERIAL PRIMARY KEY,
    system_id TEXT NOT NULL,
    user_name TEXT NOT NULL,
    team_id   INT REFERENCES teams(team_id) ON DELETE SET NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE TABLE pull_requests (
    pull_request_id   TEXT PRIMARY KEY,
    pull_request_name TEXT NOT NULL,
    author_id         INT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    status            TEXT NOT NULL,
    created_at        TIMESTAMP NOT NULL DEFAULT NOW(),
    merged_at         TIMESTAMP
);

CREATE TABLE pull_request_reviewers (
    pull_request_id TEXT NOT NULL REFERENCES pull_requests(pull_request_id) ON DELETE CASCADE,
    user_id         INT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    PRIMARY KEY (pull_request_id, user_id)
);

