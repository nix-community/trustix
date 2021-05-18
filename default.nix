{ pkgs ? import ./nix }:

{
  inherit (pkgs) nix-nar-unpack trustix trustix-doc trustix-nix trustix-nix-reprod;
}
