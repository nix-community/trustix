-- name: GetDerivationOutputsRecursive :many
SELECT derivationoutput.*
    FROM derivationoutput
    JOIN derivationrefrecursive refs_recurse ON refs_recurse.drv_id = derivationoutput.derivation_id
    JOIN derivation referrer_drv ON referrer_drv.id = refs_recurse.referrer_id
    WHERE referrer_drv.drv = ?
    ;

-- name: GetDerivationOutputResultsRecursive :many
SELECT derivationoutputresult.*
    FROM derivationoutputresult
    JOIN derivationoutput drvoutput ON drvoutput.store_path = derivationoutputresult.store_path
    JOIN derivationrefrecursive refs_recurse ON refs_recurse.drv_id = drvoutput.derivation_id
    JOIN derivation referrer_drv ON referrer_drv.id = refs_recurse.referrer_id
    WHERE referrer_drv.drv = ?
    ;
