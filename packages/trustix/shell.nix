{ pkgs ? import ../../pkgs.nix { } }:

let
  goEnv = pkgs.mkGoEnv { pwd = ./.; };
in
pkgs.mkShell {
  buildInputs = [
    goEnv
    pkgs.systemfd # Socket activation testing
  ];

  shellHook = ''
    TRUSTIX_STATE_DIR="$STATE_DIR/trustix";
  '';

}
