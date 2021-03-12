let
  pkgs = import <nixpkgs> { overlays = import ../../../../nix/overlays.nix; };

  shellDrv = pkgs.npmlock2nix.shell {
    src = import ./src.nix { inherit (pkgs) lib; };
  };

in shellDrv.overrideAttrs(old: {
  buildInputs = old.buildInputs ++ [
    pkgs.hivemind
  ];
})
