let
  pkgs = import ../../nix;
  inherit (pkgs) poetry2nix;

  pythonEnv = poetry2nix.mkPoetryEnv { projectDir = ./.; };

in
pkgs.mkShell {
  buildInputs = [ pkgs.poetry pythonEnv ];
}
