CREATE INDEX IF NOT EXISTS idx_pull_requests_is_opened ON pull_requests (is_opened);
CREATE INDEX IF NOT EXISTS idx_pr_reviewer1_open ON pull_requests (reviewer1_id) WHERE is_opened = TRUE;
CREATE INDEX IF NOT EXISTS idx_pr_reviewer2_open ON pull_requests (reviewer2_id) WHERE is_opened = TRUE;

CREATE INDEX IF NOT EXISTS idx_pr_author_id ON pull_requests (author_id);

CREATE INDEX IF NOT EXISTS idx_users_is_active ON users (is_active);
