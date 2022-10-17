-- name: GetDerivationReproducibility :many
SELECT
  drv.drv,
  drvoutput.output,
  drvoutput.store_path,
  json_group_object(
    log.log_id,
    drvoutputresult.output_hash
  ) AS output_results
FROM
  derivationoutput AS drvoutput
  JOIN derivation drv ON drv.id = drvoutput.derivation_id
  LEFT JOIN derivationoutputresult drvoutputresult ON drvoutputresult.store_path = drvoutput.store_path
  LEFT JOIN log log ON log.id = drvoutputresult.log_id
  JOIN derivationrefrecursive refs_recurse ON refs_recurse.drv_id = drvoutput.derivation_id
  JOIN derivation referrer_drv ON referrer_drv.id = refs_recurse.referrer_id
WHERE
  referrer_drv.drv = ?
GROUP BY
  drvoutput.id;

-- name: GetDerivationReproducibilityTimeSeriesByAttr :many
SELECT
  eval.id AS eval_id,
  eval.timestamp AS eval_timestamp,
  drv.drv,
  COUNT(drvoutputresult.id) AS result_count,
  COUNT(drvoutputresult.output_hash) AS output_hash_count
FROM
  derivationoutput AS drvoutput
  JOIN evaluation eval ON eval.id = drveval.eval
  JOIN derivationeval drveval ON drveval.drv = drvoutput.derivation_id
  JOIN derivation drv ON drv.id = drvattr.derivation_id
  LEFT JOIN derivationoutputresult drvoutputresult ON drvoutputresult.store_path = drvoutput.store_path
  JOIN derivationrefrecursive refs_recurse ON refs_recurse.drv_id = drvoutput.derivation_id
  JOIN derivationattr drvattr ON drvattr.derivation_id = refs_recurse.referrer_id
WHERE
  drvattr.attr = ?
  AND eval.timestamp >= ?
  AND eval.timestamp <= ?
  AND eval.channel = ?
GROUP BY
  eval.id,
  drv.id,
  drvoutput.id
ORDER BY
  eval.timestamp
  ;
