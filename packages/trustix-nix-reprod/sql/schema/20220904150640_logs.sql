-- +goose Up
-- +goose StatementBegin
CREATE TABLE log (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    log_id VARCHAR(255) NOT NULL,
    tree_size INT NOT NULL
);
CREATE INDEX idx_log_name ON log (name);


CREATE TABLE derivationoutputresult (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    output_hash VARCHAR(40) NOT NULL,

    -- Dont directly reference the derivationoutput_id field as a log may have produced an
    -- output which we might not have indexed yet.
    --
    -- This case needs loose coupling.
    store_path VARCHAR(255) NOT NULL REFERENCES derivationoutput (store_path),

    log_id INT NOT NULL REFERENCES log (id) ON DELETE CASCADE,
    UNIQUE (store_path, log_id)
);
CREATE INDEX idx_derivationoutputresult_store_path ON derivationoutputresult (store_path);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE log;
DROP TABLE derivationoutputresult;
-- +goose StatementEnd
