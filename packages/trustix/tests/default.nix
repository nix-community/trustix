{ pkgs, trustix }:

let
  inherit (pkgs) lib;

  mkTest = name: command: pkgs.runCommand "trustix-test-${name}"
    {
      nativeBuildInputs = [ trustix pkgs.systemfd ];
      meta = {
        # Does not work on darwin because of socket paths being incorrect
        platforms = builtins.filter (system: ! lib.elem system lib.platforms.darwin) lib.platforms.unix;
      };
    }
    (lib.concatStringsSep "\n" [
      ''
        export TRUSTIX_TOKEN=${../dev/token-priv}

        export HOME=$(mktemp -d)
        ln -s ${./fixtures} fixtures
        set -x
      ''
      command
      "set +x && touch $out"
    ]);

in
{
  # A simple submit/get test
  submission = mkTest "submit" ''
    key="bc63f28a4e8dda15107f687e6c3a8848492e89e3bc7726a56a0f1ee68dd9350d"
    value="28899cec2bd12feeabb5d82a3b1eafd23221798ac30a20f449144015746e2321"
    expected="5768f7201db3dccf3ec8c5ec2be5108c411396ad8c1351d89294f515456cdc23"
    log_id="5fea3cb44ef951dfb2a2ec37ebfd759174003ea9300756e26128dceb0987a30a"

    export TRUSTIX_RPC=unix://$(pwd)/sock.socket

    systemfd -s $(pwd)/sock.socket -- trustix daemon --config ${./config-simple.toml} &

    trustix --log-id "$log_id" submit --key "$key" --value "$value"
    trustix --log-id "$log_id" flush

    output=$(trustix --log-id "$log_id" query --key "$key" | cut -d" " -f 3)
    test "$output" = "$expected"
  '';

  # Test comparing multiple logs
  comparison = mkTest "compare" ''
    export TRUSTIX_TOKEN=${../dev/token-priv}

    key="bc63f28a4e8dda15107f687e6c3a8848492e89e3bc7726a56a0f1ee68dd9350d"
    output_hash="28899cec2bd12feeabb5d82a3b1eafd23221798ac30a20f449144015746e2321"
    evil_hash="053e399dbbdd74b10ad6d2cfa28ab4aab7e342d613a731c7dc4b66c2283e0757"

    log_id_1="e0e746c2a911eb67d3c68b685cd7904aeb88bf4a505225799f6e1091b76d36fa"
    log_id_2="421c8bb7aeb86eeb426cd8094e9e7cd0ad4171ce5ccd550ae622ceee7631d97c"
    log_id_3="d76ddf6362f03ae5f357fb90c37e402299578912331fa8fcd560e3df686831cc"

    build_dir=$NIX_BUILD_TOP

    # Spin up 3 log instances
    (cd ${compare-fixtures/log1}; systemfd -s $build_dir/1.socket -- trustix daemon --state $TMPDIR/log1-state --config ./config.toml) &

    (cd ${compare-fixtures/log2}; systemfd -s $build_dir/2.socket -- trustix daemon --state $TMPDIR/log2-state --config ./config.toml) &

    (cd ${compare-fixtures/log3}; systemfd -s $build_dir/3.socket -- trustix daemon --state $TMPDIR/log3-state --config ./config.toml) &

    # Submit hashes
    trustix --log-id "$log_id_1" submit --key "$key" --value "$output_hash" --address "unix://$build_dir/1.socket"
    trustix --log-id "$log_id_1" flush --address "unix://$build_dir/1.socket"

    trustix --log-id "$log_id_2" submit --key "$key" --value "$output_hash" --address "unix://$build_dir/2.socket"
    trustix --log-id "$log_id_2" flush --address "unix://$build_dir/2.socket"

    trustix --log-id "$log_id_3" submit --key "$key" --value "$evil_hash" --address "unix://$build_dir/3.socket"
    trustix --log-id "$log_id_3" flush --address "unix://$build_dir/3.socket"

    (cd ${compare-fixtures/log-agg}; systemfd -s $build_dir/agg.socket -- trustix daemon --interval 1 --state $TMPDIR/log-agg-state --config ./config.toml) &

    # Allow the aggregate to sync
    # Ideally waiting for a synchronised state should be exposed somehow but I'm uncertain
    # about what that would look like
    sleep 5

    trustix decide --protocol test --key "$key" --address "unix://$build_dir/agg.socket" > output

    # Assert correct output
    grep "Found mismatched digest '7ab45a4e40d2c0e72291ad824f8a4b208b2921e44c283022a66e87ab7c61ee38' in log '$log_id_3'" output > /dev/null
    grep "Decided on output digest '5768f7201db3dccf3ec8c5ec2be5108c411396ad8c1351d89294f515456cdc23' with confidence 66" output > /dev/null
  '';
}
