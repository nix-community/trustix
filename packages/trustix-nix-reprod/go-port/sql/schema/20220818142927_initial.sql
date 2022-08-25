-- +goose Up
-- +goose StatementBegin
CREATE TABLE evaluation (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    commit_sha VARCHAR(40) NOT NULL,
    timestamp TIMESTAMP NOT NULL  DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(commit_sha)
);

CREATE TABLE derivation (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    drv VARCHAR(255) NOT NULL,
    system VARCHAR(255) NOT NULL,
    UNIQUE(drv)
);

CREATE TABLE IF NOT EXISTS derivationrefdirect (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    drv_id INTEGER NOT NULL REFERENCES derivation (id) ON DELETE CASCADE,
    referrer_id INTEGER NOT NULL REFERENCES derivation (id) ON DELETE CASCADE,
    UNIQUE(drv_id, referrer_id)
);
CREATE INDEX idx_derivationrefdirect_drv_id ON derivationrefdirect (drv_id);
CREATE INDEX idx_derivationrefdirect_referrer_id ON derivationrefdirect (referrer_id);

CREATE TABLE IF NOT EXISTS derivationrefrecursive (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    drv_id INTEGER NOT NULL REFERENCES derivation (id) ON DELETE CASCADE,
    referrer_id INTEGER NOT NULL REFERENCES derivation (id) ON DELETE CASCADE,
    UNIQUE(drv_id, referrer_id)
);
CREATE INDEX derivationrefrecursive_idx_drv_id ON derivationrefrecursive (drv_id);
CREATE INDEX derivationrefrecursive_idx_referrer_id ON derivationrefrecursive (referrer_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE evaluation;

DROP TABLE derivation;
-- +goose StatementEnd
