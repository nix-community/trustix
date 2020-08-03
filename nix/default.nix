let
  src = let
    meta = builtins.fromJSON (builtins.readFile ./nixpkgs.json);
  in builtins.fetchTarball {
    url = "https://github.com/nixos/nixpkgs-channels/archive/${meta.rev}.tar.gz";
    sha256 = meta.sha256;
  };

  args = {
    overlays = [
      (import ./overlay.nix)
    ];
  };

  patches = [];

  pkgs = import src args;

  patched = import (pkgs.stdenv.mkDerivation {
    name = "nixpkgs";
    inherit src patches;
    dontBuild = true;
    preferLocalBuild = true;
    fixupPhase = ":";  # We dont need to patch nixpkgs shebangs
    installPhase = ''
      mkdir -p $out
      cp -a .version * $out/
    '';
  }) args;

in if patches == [] then pkgs else patched
