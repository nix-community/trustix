{ pkgs ? import ../../nix }:

let
  rootShell = import ../../shell.nix { };

in
pkgs.mkShell rec {
  inherit (rootShell) TRUSTIX_RPC STATE_DIR;
  # Speed up compilation
  CGO_ENABLED = "0";

  TRUSTIX_STATE_DIR = STATE_DIR + "trustix";


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
