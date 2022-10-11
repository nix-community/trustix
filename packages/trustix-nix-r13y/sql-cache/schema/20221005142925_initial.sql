-- +goose Up
-- +goose StatementBegin
CREATE TABLE diffoscope (
    key VARCHAR(255) PRIMARY KEY,
    html BLOB NOT NULL,
    UNIQUE(key)
);
CREATE INDEX idx_diffoscope_key ON diffoscope (key);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE diffoscope;
-- +goose StatementEnd
