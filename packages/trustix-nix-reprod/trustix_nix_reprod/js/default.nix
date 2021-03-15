let
  pkgs = import <nixpkgs> { overlays = import ../../../../nix/overlays.nix; };
in
pkgs.npmlock2nix.build {
  src = import ./src.nix { inherit (pkgs) lib; };
  buildCommands = [ "npm run build" ];
  installPhase = "cp -r dist $out";
}
