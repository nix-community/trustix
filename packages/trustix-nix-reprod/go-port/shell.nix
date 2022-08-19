{ pkgs ? import ../../../pkgs.nix { } }:

let
  rootShell = import ../../../shell.nix;

  goEnv = pkgs.mkGoEnv {
    pwd = ./.;
  };

in
pkgs.mkShell {

  inherit (rootShell) TRUSTIX_RPC TRUSTIX_ROOT;

  NIX_REPROD_STATE_DIR = "${rootShell.STATE_DIR}/nix-reprod";

  CGO_ENABLED = false;

  buildInputs = [
    pkgs.nixpkgs-fmt
    pkgs.nix-eval-jobs

    goEnv

    pkgs.go
    pkgs.gomod2nix
    pkgs.sqlite

    pkgs.hivemind
  ];

  # shellHook = ''
  #   export TRUSTIX_BINARY_CACHE_PROXY="http://localhost:8080"
  #   export DB_URI="$(./tools/tool_attr PSQL_DB_URI)"
  # '';

}
