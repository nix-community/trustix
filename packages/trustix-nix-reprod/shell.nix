{ pkgs ? import ../../pkgs.nix { } }:

let
  goEnv = pkgs.mkGoEnv {
    pwd = ./.;
  };

in
pkgs.mkShell {
  buildInputs = [
    pkgs.nix-eval-jobs
    goEnv
    pkgs.sqlite
  ];

  shellHook = ''
    export TRUSTIX_NIX_REPROD_STATE_DIR="$STATE_DIR/nix-reprod"
  '';
}
