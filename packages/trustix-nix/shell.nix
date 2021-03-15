{ pkgs ? import ../../nix }:

pkgs.mkShell {
  # Speed up compilation
  CGO_ENABLED = "0";

  buildInputs = [
    pkgs.go
    pkgs.gomod2nix
  ];

}
