-- +goose Up
-- +goose StatementBegin
CREATE TABLE evaluation (
    commit_sha VARCHAR(40) NOT NULL  PRIMARY KEY,
    timestamp TIMESTAMP NOT NULL  DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE evaluation;
-- +goose StatementEnd
