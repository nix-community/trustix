{ buildGoApplication, lib, pkg-config, gitignoreSource }:

buildGoApplication {
  pname = "trustix_nix";
  version = "dev";

  pwd = ./.;

  src = lib.cleanSourceWith {
    filter = name: type: ! lib.hasSuffix "tests" name;
    src = gitignoreSource ./.;
  };

  modules = ./gomod2nix.toml;

  subPackages = [ "." ];

  nativeBuildInputs = [ pkg-config ];

  CGO_ENABLED = "1";

}
