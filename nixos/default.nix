flake:
{ config, lib, pkgs, system, ... }:

{
  nixpkgs.overlays = [
    (_: _: flake.packages.${system})
  ];

  imports = [
    ../packages/trustix/nixos
    ../packages/trustix-nix/nixos
    ../packages/trustix-nix-r13y/nixos
  ];
}
