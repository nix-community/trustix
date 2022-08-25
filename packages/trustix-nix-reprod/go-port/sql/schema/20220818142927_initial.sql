-- +goose Up
-- +goose StatementBegin
CREATE TABLE evaluation (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    commit_sha VARCHAR(40) NOT NULL,
    timestamp TIMESTAMP NOT NULL  DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(commit_sha)
);

CREATE TABLE IF NOT EXISTS derivation (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    drv VARCHAR(255) NOT NULL,
    system VARCHAR(255) NOT NULL,
    UNIQUE(drv)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE evaluation;

DROP TABLE derivation;
-- +goose StatementEnd
