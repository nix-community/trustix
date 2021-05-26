{ config, lib, pkgs, ... }:

{
  config = {
    imports = [
      ./binarycache.nix
      ./post-build-hook.nix
    ];
  };
}
