{ pkgs ? import ../../nix }:

pkgs.mkShell {
  buildInputs = [
    pkgs.rustc
    pkgs.cargo
  ];
}
