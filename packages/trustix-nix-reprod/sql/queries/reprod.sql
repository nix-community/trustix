-- name: GetDerivationOutputs :many
SELECT derivationoutput.*
    FROM derivationoutput
    JOIN derivation referrer_drv ON referrer_drv.id = derivationoutput.derivation_id
    WHERE referrer_drv.drv = "/nix/store/vl9sbkzc0926kdzsw3vsc8acg5gxdc0h-jq-1.6.drv"
    ;

-- name: GetDerivationOutputResultsRecursive :many
SELECT derivationoutputresult.*
    FROM derivationoutputresult
    JOIN derivationoutput drvoutput ON drvoutput.store_path = derivationoutputresult.store_path
    JOIN derivation drv ON drv.id = drvoutput.derivation_id
    JOIN derivationrefrecursive refs_recurse ON refs_recurse.drv_id = drv.id
    JOIN derivation referrer_drv ON referrer_drv.id = refs_recurse.referrer_id
    WHERE referrer_drv.drv = ?
    ;
