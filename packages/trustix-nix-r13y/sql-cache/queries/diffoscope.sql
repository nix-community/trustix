-- name: GetDiffoscopeHTML :one
SELECT
  html
FROM
  diffoscope
WHERE
  key = ?
LIMIT
  1;

-- name: CreateDiffoscope :one
INSERT OR REPLACE INTO
  diffoscope(key, html)
VALUES
  (?, ?) RETURNING *;
