{ buildGoApplication, lib, gitignoreSource }:

buildGoApplication {
  pname = "trustix-nix-r13y";
  version = "dev";
  pwd = ./.;
  src = lib.cleanSourceWith {
    filter = name: type: ! lib.hasSuffix "tests" name;
    src = gitignoreSource ./.;
  };
  modules = ./gomod2nix.toml;
}
