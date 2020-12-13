{ pkgs ? import <nixpkgs> {
    overlays = [
      (import ./nix/overlay.nix)
      (import "${(builtins.fetchGit {
      url = "git@github.com:tweag/gomod2nix.git";
      rev = "929d740884811b388acc0f037efba7b5bc5745e8";
    })}/overlay.nix")
    ];
  }
}:
let
  inherit (pkgs) lib;

in
pkgs.buildGoApplication {
  pname = "trustix";
  version = "dev";

  src = lib.cleanSourceWith {
    filter = name: type: ! lib.hasSuffix "tests" name;
    src = lib.cleanSource ./.;
  };

  modules = ./gomod2nix.toml;

  subPackages = [ "." ];

  nativeBuildInputs = [ pkgs.pkgconfig ];

  CGO_ENABLED = "1";

}
