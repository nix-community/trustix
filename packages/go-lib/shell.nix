{ pkgs ? import ../../pkgs.nix { } }:

let
  rootShell = import ../../shell.nix;

  goEnv = pkgs.mkGoEnv {
    pwd = ./.;
  };

in
pkgs.mkShell {

  inherit (rootShell) TRUSTIX_RPC TRUSTIX_ROOT;

  NIX_REPROD_STATE_DIR = "${rootShell.STATE_DIR}/nix-reprod";

  CGO_ENABLED = false;

  buildInputs = [
    goEnv
    pkgs.go
    pkgs.gomod2nix
  ];

}
