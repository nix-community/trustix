{ pkgs ? import ../../pkgs.nix { } }:

let
  inherit (pkgs) poetry2nix;

  pythonEnv = poetry2nix.mkPoetryEnv {
    projectDir = ./.;
    overrides = poetry2nix.overrides.withDefaults (
      import ./overrides.nix { inherit pkgs; }
    );
  };

  rootShell = import ../../shell.nix;

in
pkgs.mkShell {

  inherit (rootShell) TRUSTIX_RPC TRUSTIX_ROOT;

  NIX_REPROD_STATE_DIR = "${rootShell.STATE_DIR}/nix-reprod";

  buildInputs = [
    pkgs.nixpkgs-fmt
    pkgs.poetry
    pythonEnv

    pkgs.hydra-eval-jobs

    pkgs.redis

    pkgs.postgresql

    pkgs.hivemind

    pkgs.yajl
    pkgs.pkg-config
  ];

  shellHook = ''
    export TRUSTIX_BINARY_CACHE_PROXY="http://localhost:8080"
    export DB_URI="$(./tools/tool_attr PSQL_DB_URI)"
  '';

}
