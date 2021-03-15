{ pkgs ? import ../../nix }:

pkgs.mkShell {

  # Speed up compilation
  CGO_ENABLED = "0";

  buildInputs = [
    pkgs.hivemind # Process monitoring in development
    pkgs.go

    pkgs.gomod2nix

    pkgs.protobuf

    pkgs.systemfd # Socket activation testing
  ];

}
