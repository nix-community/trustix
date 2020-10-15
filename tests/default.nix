{ pkgs ? import <nixpkgs> {} }:

let
  inherit (pkgs) lib;

  trustix = import ../default.nix { };

  mkTest = name: command: pkgs.runCommand "trustix-test-${name}" {
    nativeBuildInputs = [ trustix ];
  } (lib.concatStringsSep "\n" [
    ''
      export HOME=$(mktemp -d)
      ln -s ${./fixtures} fixtures
    ''
    command
    "touch $out"
  ]);

in {

  # A simple submit/get test
  submission = mkTest "submit" ''
    input_hash="bc63f28a4e8dda15107f687e6c3a8848492e89e3bc7726a56a0f1ee68dd9350d"
    output_hash="28899cec2bd12feeabb5d82a3b1eafd23221798ac30a20f449144015746e2321"

    trustix --config ${./config-simple.toml} &

    trustix submit --input-hash "$input_hash" --output-hash "$output_hash"

    echo "Checking input equality"
    test $(trustix query --input-hash "$input_hash" | cut -d' ' -f 3) = "$output_hash"
  '';

  # Test comparing multiple logs
  comparison = mkTest "compare" ''
    input_hash="bc63f28a4e8dda15107f687e6c3a8848492e89e3bc7726a56a0f1ee68dd9350d"
    output_hash="28899cec2bd12feeabb5d82a3b1eafd23221798ac30a20f449144015746e2321"
    evil_hash="053e399dbbdd74b10ad6d2cfa28ab4aab7e342d613a731c7dc4b66c2283e0757"

    # Spin up 3 log instances
    (cd ${compare-fixtures/log1}; trustix --state $TMPDIR/log1-state --config ./config.toml --address ":8081") &
    (cd ${compare-fixtures/log2}; trustix --state $TMPDIR/log1-state --config ./config.toml --address ":8082") &
    (cd ${compare-fixtures/log3}; trustix --state $TMPDIR/log1-state --config ./config.toml --address ":8083") &

    sleep 2

    # Submit hashes
    trustix submit --input-hash "$input_hash" --output-hash "$output_hash" --address ":8081"
    trustix submit --input-hash "$input_hash" --output-hash "$output_hash" --address ":8082"
    trustix submit --input-hash "$input_hash" --output-hash "$evil_hash" --address ":8083"

    (cd ${compare-fixtures/log-agg}; trustix --state $TMPDIR/log-agg-state --config ./config.toml --address ":8080") &

    trustix decide --input-hash "$input_hash" --address ":8080"
  '';

}
