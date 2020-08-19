{ pkgs ? import <nixpkgs> {} }:

pkgs.poetry2nix.mkPoetryApplication {
  projectDir = ./.;

  overrides = pkgs.poetry2nix.overrides.withDefaults(self: super: {
    pygit2 = super.pygit2.overridePythonAttrs(old: {
      buildInputs = old.buildInputs ++ [
        pkgs.libgit2
      ];
    });
  });
}
