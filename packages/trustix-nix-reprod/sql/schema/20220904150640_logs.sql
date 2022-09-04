-- +goose Up
-- +goose StatementBegin
CREATE TABLE log (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(255) NOT NULL,
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
    UNIQUE (output_id, log_id)
);
CREATE INDEX idx_derivationoutputresult_output_id ON derivationoutputresult (output_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE log;
DROP TABLE derivationoutputresult;
-- +goose StatementEnd
