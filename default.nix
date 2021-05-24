{ pkgs ? import ./nix }:

{
  inherit (pkgs) trustix trustix-doc trustix-nix trustix-nix-reprod;
}
