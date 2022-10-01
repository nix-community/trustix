-- name: GetAttrDrvs :many
EXPLAIN QUERY PLAN SELECT drv FROM derivation
  JOIN derivationattr attr ON derivation.id = attr.derivation_id
  JOIN derivationattrd attr ON derivation.id = attr.derivation_id
  WHERE attr.attr = "hello"
  ;
