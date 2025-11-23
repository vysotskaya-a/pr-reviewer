CREATE EXTENSION IF NOT EXISTS "pgcrypto";


CREATE TABLE IF NOT EXISTS teams (
    team_name TEXT PRIMARY KEY,
    description TEXT,
    created_at TIMESTAMPTZ DEFAULT now()
);


CREATE TABLE IF NOT EXISTS users (
    user_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username TEXT NOT NULL,
    display_name TEXT,
    is_active BOOLEAN NOT NULL DEFAULT true,
    team_name TEXT REFERENCES teams(team_name) ON DELETE RESTRICT,
    created_at TIMESTAMPTZ DEFAULT now()
);


CREATE TYPE IF NOT EXISTS pr_status AS ENUM ('OPEN','MERGED');


CREATE TABLE IF NOT EXISTS prs (
    pull_request_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    pull_request_name TEXT NOT NULL,
    author_id UUID NOT NULL REFERENCES users(user_id),
    team_name TEXT NOT NULL REFERENCES teams(team_name),
    status pr_status NOT NULL DEFAULT 'OPEN',
    createdAt TIMESTAMPTZ DEFAULT now(),
    mergedAt TIMESTAMPTZ NULL
);


CREATE TABLE IF NOT EXISTS pr_reviewers (
    pull_request_id UUID NOT NULL REFERENCES prs(pull_request_id) ON DELETE CASCADE,
    reviewer_id UUID NOT NULL REFERENCES users(user_id),
    assigned_at TIMESTAMPTZ DEFAULT now(),
    PRIMARY KEY (pull_request_id, reviewer_id)
);