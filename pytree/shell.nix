{ pkgs ? import <nixpkgs> {
  overlays = [
    (import ../../../nix-community/poetry2nix/overlay.nix)
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
