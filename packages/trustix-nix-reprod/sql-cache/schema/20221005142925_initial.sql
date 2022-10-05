-- +goose Up
-- +goose StatementBegin
CREATE TABLE diffoscope (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    key VARCHAR(255) NOT NULL,
    timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    html BLOB NOT NULL
);
CREATE INDEX idx_diffoscope_key ON diffoscope (key);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE diffoscope;
-- +goose StatementEnd
