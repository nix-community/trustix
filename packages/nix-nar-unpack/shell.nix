{ pkgs ? import ../../nix }:

pkgs.mkShell {
  buildInputs = [
    pkgs.rustc
    pkgs.cargo
    pkgs.cargo-watch
    pkgs.hivemind
  ];
}
