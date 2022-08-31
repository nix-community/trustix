{ pkgs ? import ../../pkgs.nix { } }:

let
  inherit (pkgs) poetry2nix;

  pythonEnv = poetry2nix.mkPoetryEnv {
    projectDir = ./.;
    python = pkgs.python39;
    overrides = poetry2nix.overrides.withDefaults (
      import ./overrides.nix { inherit pkgs; }
    );
  };

in
pkgs.mkShell {

  buildInputs = [
    pkgs.poetry
    pythonEnv

    pkgs.hydra-eval-jobs

    pkgs.redis

    pkgs.postgresql

    pkgs.yajl
    pkgs.pkg-config
  ];

  shellHook = ''
    export TRUSTIX_BINARY_CACHE_PROXY="http://localhost:8080"
    export DB_URI="$(./tools/tool_attr PSQL_DB_URI)"
    export NIX_REPROD_STATE_DIR="$STATE_DIR/nix-reprod";
  '';

}
