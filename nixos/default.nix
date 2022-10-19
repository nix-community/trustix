{ config, lib, pkgs, ... }:

{
  nixpkgs.overlays = [
    (final: prev: import ../default.nix { })
  ];

  imports = [
    ../packages/trustix/nixos
    ../packages/trustix-nix/nixos
    ../packages/trustix-nix-r13y/nixos
  ];
}
