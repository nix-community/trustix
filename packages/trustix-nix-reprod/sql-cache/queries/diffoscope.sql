-- name: GetDiffoscopeHTML :one
SELECT
  html
FROM
  diffoscope
WHERE
  key = ?
ORDER BY
  timestamp
LIMIT
  1;

-- name: CreateDiffoscope :one
INSERT INTO
  diffoscope(key, html)
VALUES
  (?, ?) RETURNING *;
