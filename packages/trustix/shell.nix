{ pkgs ? import ../../pkgs.nix { } }:

let
  rootShell = import ../../shell.nix;

in
pkgs.mkShell {

  # Speed up compilation
  CGO_ENABLED = "0";

  TRUSTIX_STATE_DIR = rootShell.STATE_DIR + "/trustix";
  inherit (rootShell) TRUSTIX_RPC;

  buildInputs = [
    pkgs.hivemind # Process monitoring in development
    pkgs.go

    pkgs.reflex
    pkgs.entr

    pkgs.gomod2nix

    pkgs.protobuf

    pkgs.systemfd # Socket activation testing
  ];

}
