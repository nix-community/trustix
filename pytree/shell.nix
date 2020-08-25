{ pkgs ? import <nixpkgs> {
  overlays = [
    (import ../../../nix-community/poetry2nix/overlay.nix)

    (self: super: {
      libgit2 = super.libgit2.overrideAttrs(old: let
        version = "unstable-2020-09-05";
      in {
        name = "libgit2-${version}";
        inherit version;
        src = super.fetchFromGitHub {
          owner = "libgit2";
          repo = "libgit2";
          rev = "04d59466238e69c57d2a82d0693a77ecb05e1194";
          sha256 = "1ppbja0cmw6x8y92zsj4vyr0mwhiq1sinq2vkffryrzjmw4d6r7j";
        };
      });      # Bump libgit2 for partial clone support
    })

  ];
} }:

let
  pythonEnv = pkgs.poetry2nix.mkPoetryEnv {
    projectDir = ./.;
    overrides = pkgs.poetry2nix.overrides.withDefaults(self: super: {
      pygit2 = super.pygit2.overridePythonAttrs(old: {
        buildInputs = old.buildInputs ++ [
          pkgs.libgit2
        ];
      });
    });
  };

in pkgs.mkShell {
  buildInputs = [
    pkgs.poetry
    pythonEnv
  ];
}
