{ pkgs ? import <nixpkgs> {} }:

let
  inherit (pkgs) lib;

  trustix = import ../default.nix { };

  mkTest = name: command: pkgs.runCommand "trustix-test-${name}" {
    nativeBuildInputs = [ trustix pkgs.systemfd ];
  } (lib.concatStringsSep "\n" [
    ''
      export HOME=$(mktemp -d)
      ln -s ${./fixtures} fixtures
    ''
    command
    "touch $out"
  ]);

in {

  inherit trustix;

  # A simple submit/get test
  submission = mkTest "submit" ''
    input_hash="bc63f28a4e8dda15107f687e6c3a8848492e89e3bc7726a56a0f1ee68dd9350d"
    output_hash="28899cec2bd12feeabb5d82a3b1eafd23221798ac30a20f449144015746e2321"

    export TRUSTIX_SOCK=./sock

    systemfd -s ./sock -- trustix --config ${./config-simple.toml} &

    trustix submit --input-hash "$input_hash" --output-hash "$output_hash"

    echo "Checking input equality"
    test $(trustix query --input-hash "$input_hash" | cut -d' ' -f 3) = "$output_hash"
  '';

  # Test comparing multiple logs
  comparison = mkTest "compare" ''
    input_hash="bc63f28a4e8dda15107f687e6c3a8848492e89e3bc7726a56a0f1ee68dd9350d"
    output_hash="28899cec2bd12feeabb5d82a3b1eafd23221798ac30a20f449144015746e2321"
    evil_hash="053e399dbbdd74b10ad6d2cfa28ab4aab7e342d613a731c7dc4b66c2283e0757"

    build_dir=$(pwd)

    # Spin up 3 log instances
    (cd ${compare-fixtures/log1}; systemfd -s $build_dir/1.sock -- trustix --state $TMPDIR/log1-state --config ./config.toml --listen "tcp://:8081") &
    (cd ${compare-fixtures/log2}; systemfd -s $build_dir/2.sock -- trustix --state $TMPDIR/log2-state --config ./config.toml --listen "tcp://:8082") &
    (cd ${compare-fixtures/log3}; systemfd -s $build_dir/3.sock -- trustix --state $TMPDIR/log3-state --config ./config.toml --listen "tcp://:8083") &

    # Submit hashes
    trustix submit --input-hash "$input_hash" --output-hash "$output_hash" --address "unix://./1.sock"
    trustix submit --input-hash "$input_hash" --output-hash "$output_hash" --address "unix://./2.sock"
    trustix submit --input-hash "$input_hash" --output-hash "$evil_hash" --address "unix://./3.sock"

    (cd ${compare-fixtures/log-agg}; systemfd -s $build_dir/agg.sock -- trustix --state $TMPDIR/log-agg-state --config ./config.toml) &

    trustix decide --input-hash "$input_hash" --address "unix://./agg.sock" > output

    # Assert correct output
    grep "Found mismatched hash '053e399dbbdd74b10ad6d2cfa28ab4aab7e342d613a731c7dc4b66c2283e0757' in log 'trustix-test-follower3'" output > /dev/null
    grep "Decided on output hash '28899cec2bd12feeabb5d82a3b1eafd23221798ac30a20f449144015746e2321' with confidence 66" output > /dev/null
  '';
}
# Storage (one test per storage type)
//
(
  let
    configTemplate = builtins.fromTOML (builtins.readFile ./config-simple.toml);
    old = lib.elemAt configTemplate.log 0;

  in lib.listToAttrs (builtins.map (storageType: let
    name = "storage-${storageType}";

    config = let
      drv = pkgs.writeText "conf-${storageType}" (builtins.toJSON {
        log = [ (old // { storage = { type = "native"; }; }) ];
      });
    in drv.overrideAttrs(old: {
      buildCommand = old.buildCommand + ''
        ${pkgs.remarshal}/bin/remarshal -i $out -o config.toml -if json -of toml
        mv config.toml $out
      '';
    });

  in {
    inherit name;
    value = mkTest name ''
      input_hash="bc63f28a4e8dda15107f687e6c3a8848492e89e3bc7726a56a0f1ee68dd9350d"
      output_hash="28899cec2bd12feeabb5d82a3b1eafd23221798ac30a20f449144015746e2321"

      systemfd -s ./log.sock -- trustix --config ${config} &

      export TRUSTIX_SOCK=./log.sock

      trustix submit --input-hash "$input_hash" --output-hash "$output_hash"

      echo "Checking input equality"
      test $(trustix query --input-hash "$input_hash" | cut -d' ' -f 3) = "$output_hash"
    '';
  }) [
    "native"
    "memory"
    "git"
  ])
)
