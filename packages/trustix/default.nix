{ buildGoApplication, lib, pkg-config, gitignoreSource, callPackage, pkgs }:

lib.fix (self: buildGoApplication {
  pname = "trustix";
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

  passthru.tests = import ./tests {
    inherit pkgs;
    trustix = self;
  };
})
