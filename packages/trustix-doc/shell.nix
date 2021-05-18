{ pkgs ? import ../../nix }:

let
  rootShell = import ../../shell.nix;

in
pkgs.mkShell {
  buildInputs = [
    pkgs.mdbook
  ];
}
