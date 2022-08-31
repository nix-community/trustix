{ pkgs ? import ../../pkgs.nix { } }:
let
  goEnv = pkgs.mkGoEnv { pwd = ./.; };
in
pkgs.mkShell {
  buildInputs = [
    goEnv
  ];
}
