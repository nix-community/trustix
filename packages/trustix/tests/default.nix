{ pkgs ? import ../../../nix }:
let
  inherit (pkgs) trustix lib;

  mkTest = name: command: pkgs.runCommand "trustix-test-${name}"
    {
      nativeBuildInputs = [ trustix pkgs.systemfd ];
    }
    (lib.concatStringsSep "\n" [
      ''
        export HOME=$(mktemp -d)
        ln -s ${./fixtures} fixtures
        set -x
      ''
      command
      "set +x && touch $out"
    ]);

in
{

  inherit trustix;

  # A simple submit/get test
  submission = mkTest "submit" ''
    key="bc63f28a4e8dda15107f687e6c3a8848492e89e3bc7726a56a0f1ee68dd9350d"
    value="28899cec2bd12feeabb5d82a3b1eafd23221798ac30a20f449144015746e2321"
    expected="5768f7201db3dccf3ec8c5ec2be5108c411396ad8c1351d89294f515456cdc23"
    log_id="e675d5665a41a51841ccc51e21d8b1ad54cd48a36623d1287ed29bb0a043d36b"

    export TRUSTIX_SOCK=./sock

    systemfd -s ./sock -- trustix --config ${./config-simple.toml} &

    trustix --log-id "$log_id" submit --input-hash "$key" --output-hash "$value"
    trustix --log-id "$log_id" flush

    echo "Checking input equality"
    test $(trustix --log-id "$log_id" query --input-hash "$key" | cut -d' ' -f 3) = "$expected"
  '';

  # Test comparing multiple logs
  comparison = mkTest "compare" ''
    input_hash="bc63f28a4e8dda15107f687e6c3a8848492e89e3bc7726a56a0f1ee68dd9350d"
    output_hash="28899cec2bd12feeabb5d82a3b1eafd23221798ac30a20f449144015746e2321"
    evil_hash="053e399dbbdd74b10ad6d2cfa28ab4aab7e342d613a731c7dc4b66c2283e0757"

    log_id_1="8d240243184ee561e6becd60c9ee7d51a1367f7ba0d11129bb2aa080184ecc8c"
    log_id_2="c4d3f0c4c0648d4d7c80a720931378a0268c24712964049c99b504a9403b8045"
    log_id_3="b3ee46fde4e9ead1fe7bc9a3af26b9208d5c4f7732b2d1484742562deab7a21d"

    build_dir=$(pwd)

    # Spin up 3 log instances
    (cd ${compare-fixtures/log1}; systemfd -s $build_dir/1.sock -- trustix --state $TMPDIR/log1-state --config ./config.toml) &

    (cd ${compare-fixtures/log2}; systemfd -s $build_dir/2.sock -- trustix --state $TMPDIR/log2-state --config ./config.toml) &

    (cd ${compare-fixtures/log3}; systemfd -s $build_dir/3.sock -- trustix --state $TMPDIR/log3-state --config ./config.toml) &

    # Submit hashes
    trustix --log-id "$log_id_1" submit --input-hash "$input_hash" --output-hash "$output_hash" --address "unix://./1.sock"
    trustix --log-id "$log_id_1" flush --address "unix://./1.sock"

    trustix --log-id "$log_id_2" submit --input-hash "$input_hash" --output-hash "$output_hash" --address "unix://./2.sock"
    trustix --log-id "$log_id_2" flush --address "unix://./2.sock"

    trustix --log-id "$log_id_3" submit --input-hash "$input_hash" --output-hash "$evil_hash" --address "unix://./3.sock"
    trustix --log-id "$log_id_3" flush --address "unix://./3.sock"

    (cd ${compare-fixtures/log-agg}; systemfd -s $build_dir/agg.sock -- trustix --state $TMPDIR/log-agg-state --config ./config.toml) &

    # Allow the aggregate to sync
    # Ideally waiting for a synchronised state should be exposed somehow but I'm uncertain
    # about what that would look like
    sleep 5

    trustix decide --input-hash "$input_hash" --address "unix://./agg.sock" > output

    echo "Decision output:"
    cat output
    echo "---"

    # Assert correct output
    grep "Found mismatched digest '7ab45a4e40d2c0e72291ad824f8a4b208b2921e44c283022a66e87ab7c61ee38' in log '$log_id_3'" output > /dev/null
    grep "Decided on output digest '5768f7201db3dccf3ec8c5ec2be5108c411396ad8c1351d89294f515456cdc23' with confidence 66" output > /dev/null
  '';
}
