-- name: GetEval :one
SELECT * FROM evaluation
WHERE commit_sha = ? LIMIT 1;

-- name: GetDerivation :one
SELECT * FROM derivation
WHERE drv = ? LIMIT 1;
