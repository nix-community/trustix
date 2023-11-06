flake:
{ config, lib, pkgs, ... }:

{
  nixpkgs.overlays = [
    (_: prev: builtins.removeAttrs (flake.packages.${prev.stdenv.targetPlatform.system} or { }) [ "default" ])
  ];

  imports = [
    ../packages/trustix/nixos
    ../packages/trustix-nix/nixos
    ../packages/trustix-nix-r13y/nixos
  ];
}
