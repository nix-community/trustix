-- name: GetDerivationReproducibilityRecursive :many
SELECT
    drv.drv
    , drvoutput.output
    , drvoutput.store_path
    , json_group_object(drvoutputresult.output_hash, drvoutputresult.log_id) AS output_results
    FROM derivationoutput AS drvoutput
    JOIN derivation drv ON drv.id = drvoutput.derivation_id
    LEFT JOIN derivationoutputresult drvoutputresult ON drvoutputresult.store_path = drvoutput.store_path
    JOIN derivationrefrecursive refs_recurse ON refs_recurse.drv_id = drvoutput.derivation_id
    JOIN derivation referrer_drv ON referrer_drv.id = refs_recurse.referrer_id
    WHERE referrer_drv.drv = ?
    GROUP BY drvoutput.id
    ;

-- name: GetDerivationReproducibilityByAttr :many
SELECT
    drv.drv
    , drvoutput.output
    , drvoutput.store_path
    , json_group_object(drvoutputresult.output_hash, drvoutputresult.log_id) AS output_results
    FROM derivationoutput AS drvoutput
    JOIN derivation drv ON drv.id = drvoutput.derivation_id
    LEFT JOIN derivationoutputresult drvoutputresult ON drvoutputresult.store_path = drvoutput.store_path
    JOIN derivationrefrecursive refs_recurse ON refs_recurse.drv_id = drvoutput.derivation_id
    JOIN derivation referrer_drv ON referrer_drv.id = refs_recurse.referrer_id
    WHERE referrer_drv.drv = ?
    GROUP BY drvoutput.id
    ;

-- name: GetDerivationReproducibilityTimeSeriesByAttr :many
SELECT
    eval.id
    , eval.timestamp
    , drv.drv
    , drvoutput.output
    , drvoutput.store_path
    , json_group_object(drvoutputresult.output_hash, drvoutputresult.log_id) AS output_results
    FROM derivationoutput AS drvoutput
    JOIN evaluation eval ON eval.id = drveval.eval
    JOIN derivationeval drveval ON drveval.drv = drvoutput.derivation_id
    JOIN derivation drv ON drv.id = drvoutput.derivation_id
    LEFT JOIN derivationoutputresult drvoutputresult ON drvoutputresult.store_path = drvoutput.store_path
    JOIN derivationrefrecursive refs_recurse ON refs_recurse.drv_id = drvoutput.derivation_id
    JOIN derivationattr drvattr ON drvattr.derivation_id = refs_recurse.referrer_id
    WHERE drvattr.attr = ? AND eval.timestamp >= ? AND eval.timestamp <= ?
    GROUP BY drvoutput.id
    ;
