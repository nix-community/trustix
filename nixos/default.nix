{ config, lib, pkgs, ... }:

{
  config = {
    nixpkgs.overlays = ../nix/overlays.nix;

    imports = [
      ../packages/trustix/nixos
      ../packages/trustix-nix/nixos
    ];
  };
}
