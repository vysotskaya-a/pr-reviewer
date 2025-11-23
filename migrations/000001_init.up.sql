-- 000001_init.up.sql
CREATE TABLE teams (
                       team_name TEXT PRIMARY KEY,
                       description TEXT,
                       created_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE users (
                       user_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                       username TEXT NOT NULL UNIQUE,
                       display_name TEXT,
                       is_active BOOLEAN NOT NULL DEFAULT true,
                       team_name TEXT REFERENCES teams(team_name) ON DELETE RESTRICT,
                       created_at TIMESTAMPTZ DEFAULT now()
);

CREATE TYPE pr_status AS ENUM ('OPEN', 'MERGED');

CREATE TABLE prs (
                     pull_request_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                     pull_request_name TEXT NOT NULL,
                     author_id UUID NOT NULL REFERENCES users(user_id),
                     team_name TEXT NOT NULL REFERENCES teams(team_name),
                     status pr_status NOT NULL DEFAULT 'OPEN',
                     "createdAt" TIMESTAMPTZ DEFAULT now(),
                     "mergedAt" TIMESTAMPTZ
);

CREATE TABLE pr_reviewers (
                              pull_request_id UUID NOT NULL REFERENCES prs(pull_request_id) ON DELETE CASCADE,
                              reviewer_id UUID NOT NULL REFERENCES users(user_id),
                              assigned_at TIMESTAMPTZ DEFAULT now(),
                              PRIMARY KEY (pull_request_id, reviewer_id)
);