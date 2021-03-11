let
  pkgs = import <nixpkgs> { overlays = import ../../../../nix/overlays.nix; };

in
pkgs.npmlock2nix.shell {
  src = ./.;
}
