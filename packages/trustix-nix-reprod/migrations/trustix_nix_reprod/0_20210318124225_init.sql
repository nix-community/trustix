CREATE TABLE IF NOT EXISTS "derivation" (
    "drv" VARCHAR(255) NOT NULL  PRIMARY KEY,
    "system" VARCHAR(255) NOT NULL
);
CREATE INDEX IF NOT EXISTS "idx_derivation_system_5c2dd2" ON "derivation" ("system");

CREATE TABLE IF NOT EXISTS "derivationattr" (
    "id" SERIAL NOT NULL PRIMARY KEY,
    "attr" VARCHAR(255) NOT NULL,
    "derivation_id" VARCHAR(255) NOT NULL REFERENCES "derivation" ("drv") ON DELETE CASCADE,
    CONSTRAINT "uid_derivationa_derivat_b9328d" UNIQUE ("derivation_id", "attr")
);
CREATE INDEX IF NOT EXISTS "idx_derivationa_attr_f37b80" ON "derivationattr" ("attr");

CREATE TABLE IF NOT EXISTS "derivationoutput" (
    "input_hash" VARCHAR(25) NOT NULL  PRIMARY KEY,
    "output" VARCHAR(255) NOT NULL,
    "store_path" VARCHAR(255) NOT NULL,
    "derivation_id" VARCHAR(255) NOT NULL REFERENCES "derivation" ("drv") ON DELETE CASCADE,
    CONSTRAINT "uid_derivationo_derivat_808b7b" UNIQUE ("derivation_id", "output")
);
CREATE INDEX IF NOT EXISTS "idx_derivationo_output_8c7711" ON "derivationoutput" ("output");
CREATE INDEX IF NOT EXISTS "idx_derivationo_store_p_84d64d" ON "derivationoutput" ("store_path");

CREATE TABLE IF NOT EXISTS "derivationrefdirect" (
    "id" SERIAL NOT NULL PRIMARY KEY,
    "drv_id" VARCHAR(255) NOT NULL REFERENCES "derivation" ("drv") ON DELETE CASCADE,
    "referrer_id" VARCHAR(255) NOT NULL REFERENCES "derivation" ("drv") ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS "idx_derivationr_drv_id_cb1ce7" ON "derivationrefdirect" ("drv_id");
CREATE INDEX IF NOT EXISTS "idx_derivationr_referre_03f0ee" ON "derivationrefdirect" ("referrer_id");

CREATE TABLE IF NOT EXISTS "derivationrefrecursive" (
    "id" SERIAL NOT NULL PRIMARY KEY,
    "drv_id" VARCHAR(255) NOT NULL REFERENCES "derivation" ("drv") ON DELETE CASCADE,
    "referrer_id" VARCHAR(255) NOT NULL REFERENCES "derivation" ("drv") ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS "idx_derivationr_drv_id_ece8f8" ON "derivationrefrecursive" ("drv_id");
CREATE INDEX IF NOT EXISTS "idx_derivationr_referre_5ccd9c" ON "derivationrefrecursive" ("referrer_id");

CREATE TABLE IF NOT EXISTS "evaluation" (
    "commit" VARCHAR(40) NOT NULL  PRIMARY KEY,
    "timestamp" TIMESTAMPTZ NOT NULL  DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS "derivationeval" (
    "id" SERIAL NOT NULL PRIMARY KEY,
    "drv_id" VARCHAR(255) NOT NULL REFERENCES "derivation" ("drv") ON DELETE CASCADE,
    "eval_id" VARCHAR(40) NOT NULL REFERENCES "evaluation" ("commit") ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS "idx_derivatione_drv_id_304e08" ON "derivationeval" ("drv_id");
CREATE INDEX IF NOT EXISTS "idx_derivatione_eval_id_d1000a" ON "derivationeval" ("eval_id");

CREATE TABLE IF NOT EXISTS "log" (
    "id" SERIAL NOT NULL PRIMARY KEY,
    "name" VARCHAR(255) NOT NULL,
    "tree_size" INT NOT NULL
);
CREATE INDEX IF NOT EXISTS "idx_log_name_1bf001" ON "log" ("name");

CREATE TABLE IF NOT EXISTS "derivationoutputresult" (
    "id" SERIAL NOT NULL PRIMARY KEY,
    "output_hash" VARCHAR(40) NOT NULL,
    "output_id" VARCHAR(25) NOT NULL,
    "log_id" INT NOT NULL REFERENCES "log" ("id") ON DELETE CASCADE,
    CONSTRAINT "uid_derivationo_output__ffeb49" UNIQUE ("output_id", "log_id")
);
CREATE INDEX IF NOT EXISTS "idx_derivationo_output__5fdba2" ON "derivationoutputresult" ("output_id");

CREATE TABLE IF NOT EXISTS "aerich" (
    "id" SERIAL NOT NULL PRIMARY KEY,
    "version" VARCHAR(255) NOT NULL,
    "app" VARCHAR(20) NOT NULL,
    "content" JSONB NOT NULL
);
