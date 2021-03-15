{ pkgs ? import ../../nix }:

let
  inherit (pkgs) poetry2nix;

  pythonEnv = poetry2nix.mkPoetryEnv {
    projectDir = ./.;
    overrides = poetry2nix.overrides.withDefaults (
      import ./overrides.nix { inherit pkgs; }
    );
  };

in
pkgs.mkShell {

  buildInputs = [
    pkgs.nixpkgs-fmt
    pkgs.poetry
    pythonEnv

    pkgs.nix-nar-unpack

    pkgs.hydra-eval-jobs
    pkgs.sqlite

    pkgs.postgresql

    pkgs.diffoscope

    pkgs.hivemind

    pkgs.yajl
    pkgs.pkgconfig
  ];

  shellHook = ''
    export TRUSTIX_RPC="unix:../../sock"
    export TRUSTIX_BINARY_CACHE_PROXY="http://localhost:8080"
    export DB_URI="$(./tools/tool_attr PSQL_DB_URI)"
  '';

}
