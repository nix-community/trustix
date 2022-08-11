{ pkgs ? import ../../pkgs.nix { } }:

let
  rootShell = import ../../shell.nix;

in
pkgs.mkShell {
  # Speed up compilation
  CGO_ENABLED = "0";

  inherit (rootShell) TRUSTIX_RPC;

  buildInputs = [
    pkgs.go
    pkgs.gomod2nix
    pkgs.hivemind
    pkgs.entr
    pkgs.reflex
  ];

}
