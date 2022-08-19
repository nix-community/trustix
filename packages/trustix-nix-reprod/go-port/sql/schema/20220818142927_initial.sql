-- +goose Up
-- +goose StatementBegin
CREATE TABLE evaluation (
    commit_sha VARCHAR(40) NOT NULL  PRIMARY KEY,
    timestamp TIMESTAMP NOT NULL  DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS derivation (
    drv VARCHAR(255) NOT NULL  PRIMARY KEY,
    system VARCHAR(255) NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE evaluation;

DROP TABLE derivation;
-- +goose StatementEnd
