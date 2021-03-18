{ buildGoApplication, lib, pkgconfig }:

buildGoApplication {
  pname = "trustix";
  version = "dev";

  pwd = ./.;

  src = lib.cleanSourceWith {
    filter = name: type: ! lib.hasSuffix "tests" name;
    src = lib.cleanSource ./.;
  };

  modules = ./gomod2nix.toml;

  subPackages = [ "." ];

  nativeBuildInputs = [ pkgconfig ];

  CGO_ENABLED = "1";

}
