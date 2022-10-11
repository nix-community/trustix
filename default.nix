{ pkgs ? import ./pkgs.nix { } }:
{
  inherit (pkgs) trustix trustix-doc trustix-nix trustix-nix-r13y trustix-nix-r13y-web;
}
