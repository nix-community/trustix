{ pkgs ? import ../../pkgs.nix { } }:

pkgs.mkShell {
  buildInputs = [
    pkgs.nix-eval-jobs
    pkgs.sqlite
    pkgs.diffoscope
    pkgs.sqlc
    pkgs.goose
    pkgs.protoc-gen-go
    pkgs.protoc-gen-connect-go
  ];

  shellHook = ''
    export TRUSTIX_NIX_REPROD_STATE_DIR="$STATE_DIR/nix-reprod"
  '';
}
