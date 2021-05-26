{ config, lib, pkgs, ... }:

{
  nixpkgs.overlays = import ../nix/overlays.nix;

  imports = [
    ../packages/trustix/nixos
    ../packages/trustix-nix/nixos
  ];

}
