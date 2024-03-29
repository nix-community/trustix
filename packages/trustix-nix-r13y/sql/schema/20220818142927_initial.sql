-- +goose Up
-- +goose StatementBegin
CREATE TABLE evaluation (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    channel VARCHAR(40) NOT NULL,
    timestamp TIMESTAMP NOT NULL  DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_evaluation_timestamp ON evaluation (timestamp);

CREATE TABLE hydraevaluation (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    evaluation INTEGER NOT NULL REFERENCES evaluation (id) ON DELETE CASCADE,
    hydra_eval_id INTEGER NOT NULL,
    revision VARCHAR(40) NOT NULL,
    UNIQUE(evaluation)
);
CREATE INDEX idx_hydraevaluation_hydra_eval_id ON hydraevaluation (hydra_eval_id);
CREATE INDEX idx_hydraevaluation_evaluation ON hydraevaluation (evaluation);

CREATE TABLE derivation (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    drv VARCHAR(255) NOT NULL,
    system VARCHAR(255) NOT NULL,
    UNIQUE(drv)
);

CREATE TABLE derivationrefdirect (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    drv_id INTEGER NOT NULL REFERENCES derivation (id) ON DELETE CASCADE,
    referrer_id INTEGER NOT NULL REFERENCES derivation (id) ON DELETE CASCADE,
    UNIQUE(drv_id, referrer_id)
);
CREATE INDEX idx_derivationrefdirect_drv_id ON derivationrefdirect (drv_id);
CREATE INDEX idx_derivationrefdirect_referrer_id ON derivationrefdirect (referrer_id);

CREATE TABLE derivationrefrecursive (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    drv_id INTEGER NOT NULL REFERENCES derivation (id) ON DELETE CASCADE,
    referrer_id INTEGER NOT NULL REFERENCES derivation (id) ON DELETE CASCADE,
    UNIQUE(drv_id, referrer_id)
);
CREATE INDEX derivationrefrecursive_idx_drv_id ON derivationrefrecursive (drv_id);
CREATE INDEX derivationrefrecursive_idx_referrer_id ON derivationrefrecursive (referrer_id);

CREATE TABLE derivationoutput (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    output VARCHAR(255) NOT NULL,
    store_path VARCHAR(255) NOT NULL,
    derivation_id INTEGER NOT NULL REFERENCES derivation (id) ON DELETE CASCADE,
    UNIQUE (derivation_id, output)
);
CREATE INDEX idx_derivationoutput_output ON derivationoutput (output);
CREATE INDEX idx_derivationoutput_store_path ON derivationoutput (store_path);

CREATE TABLE derivationeval (
    drv INTEGER NOT NULL REFERENCES derivation (id) ON DELETE CASCADE,
    eval INTEGER NOT NULL REFERENCES evaluation (id) ON DELETE CASCADE
);
CREATE INDEX idx_derivationeval_drv ON derivationeval (drv);

CREATE TABLE derivationattr (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    attr VARCHAR(255) NOT NULL,
    derivation_id INTEGER NOT NULL REFERENCES derivation (id) ON DELETE CASCADE,
    UNIQUE (derivation_id, attr)
);
CREATE INDEX idx_derivationattr_attr ON derivationattr (attr);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE evaluation;
DROP TABLE derivation;
DROP TABLE derivationrefdirect;
DROP TABLE derivationrefrecursive;
DROP TABLE derivationoutput;
DROP TABLE derivationattr;
-- +goose StatementEnd
