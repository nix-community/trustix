{ config, lib, pkgs, ... }:

{
  imports = [
    ./binarycache.nix
    ./post-build-hook.nix
  ];
}
