-- name: GetLog :one
SELECT
  *
FROM
  log
WHERE
  log_id = ?
LIMIT
  1;

-- name: CreateLog :one
INSERT INTO
  log(log_id, tree_size)
VALUES
  (?, 0) RETURNING *;

-- name: CreateDerivationOutputResult :one
INSERT INTO
  derivationoutputresult (output_hash, store_path, log_id)
VALUES
  (?, ?, ?) RETURNING *;

-- name: SetTreeSize :exec
UPDATE
  log
SET
  tree_size = ?
WHERE
  id = ?;
