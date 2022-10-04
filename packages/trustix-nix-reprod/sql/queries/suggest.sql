-- name: SuggestAttribute :many
SELECT
  drvattr.attr
FROM
  derivationattr AS drvattr
WHERE
  drvattr.attr LIKE ?
LIMIT 100
  ;
