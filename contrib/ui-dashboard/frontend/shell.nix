let
  pkgs = import <nixpkgs> { overlays = import ../../../nix/overlays.nix; };

in pkgs.mkShell {
  buildInputs = [
    pkgs.nodejs
  ];
}
