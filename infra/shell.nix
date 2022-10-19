{ pkgs ? import ../pkgs.nix { } }:

pkgs.mkShell {
  buildInputs = [
    pkgs.colmena
  ];
}
