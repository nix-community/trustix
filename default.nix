{ pkgs ? import <nixpkgs> {
  overlays = [
    (import ./nix/overlay.nix)
    (import "${(builtins.fetchGit {
      url = "git@github.com:tweag/gomod2nix.git";
      rev = "1d342d55cd9476af8e165b60a55dbec0d8c10977";
    })}/overlay.nix")
  ];
} }:

let
  inherit (pkgs) lib;

in pkgs.buildGoApplication {
  pname = "trustix";
  version = "dev";

  src = lib.cleanSourceWith {
    filter = name: type: ! lib.hasSuffix "tests" name;
    src = lib.cleanSource ./.;
  };

  modules = ./gomod2nix.toml;

  nativeBuildInputs = [
    pkgs.pkgconfig
  ];

  CGO_ENABLED = "1";

}
