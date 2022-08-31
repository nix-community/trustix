{ buildGoApplication, lib }:

buildGoApplication {
  pname = "trustix-nix-reprod";
  version = "dev";
  pwd = ./.;
  src = lib.cleanSourceWith {
    filter = name: type: ! lib.hasSuffix "tests" name;
    src = lib.cleanSource ./.;
  };
  modules = ./gomod2nix.toml;
}
