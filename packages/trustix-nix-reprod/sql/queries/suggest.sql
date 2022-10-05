-- name: SuggestAttribute :many
SELECT
  drvattr.attr
FROM
  derivationattr AS drvattr
WHERE
  drvattr.attr LIKE ?
ORDER BY
  drvattr.attr
LIMIT 100
  ;
