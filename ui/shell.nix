let
  pkgs = import ../nix;
  inherit (pkgs) poetry2nix;

  pythonEnv = poetry2nix.mkPoetryEnv {
    projectDir = ./.;
    overrides = poetry2nix.overrides.withDefaults (
      import ./overrides.nix { inherit pkgs; }
    );
  };

in
pkgs.mkShell {

  buildInputs = [
    pkgs.nixpkgs-fmt
    pkgs.poetry
    pythonEnv
    pkgs.hydra-unstable  # For the hydra evaluator

    pkgs.yajl
    pkgs.pkgconfig
  ];

}
