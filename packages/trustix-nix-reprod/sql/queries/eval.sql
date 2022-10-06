-- name: CreateEval :one
INSERT INTO
  evaluation (channel, timestamp)
VALUES
  (?, ?) RETURNING *;

-- name: CreateHydraEval :one
INSERT INTO
  hydraevaluation (evaluation, hydra_eval_id, revision)
VALUES
  (?, ?, ?) RETURNING *;

-- name: GetLatesteval :one
SELECT
  *
FROM
  hydraevaluation AS hydraeval
  JOIN evaluation eval ON eval.id = drveval.eval
WHERE
  eval.channel = ?
  AND hydraeval.hydra_eval_id = ?
ORDER BY
  timestamp
LIMIT
  1;

-- name: GetEvalByHydraID :one
SELECT
  eval.*
FROM
  hydraevaluation AS hydraeval
  JOIN evaluation eval ON eval.id = hydraeval.evaluation
WHERE
  eval.channel = ?
  AND hydraeval.hydra_eval_id = ?
LIMIT
  1;
