{ pkgs ? import <nixpkgs> {
    overlays = [
      (import ./nix/overlay.nix)
      (import "${(builtins.fetchGit {
      url = "git@github.com:tweag/gomod2nix.git";
      rev = "d43300e22e1d379b80e4736c6583d5b9078b3c45";
    })}/overlay.nix")
    ];
  }
}:
let
  inherit (pkgs) lib;

in
pkgs.buildGoApplication {
  pname = "trustiqx";
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
