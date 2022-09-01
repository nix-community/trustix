let
  pkgs = import ./pkgs.nix { };

  STATE_DIR = "${builtins.toString ./.}/state";
  TRUSTIX_RPC = "unix://${STATE_DIR}/trustix.sock";
  TRUSTIX_ROOT = builtins.toString ./.;

in
pkgs.mkShell {

  # Speed up compilation, guarantee static linking
  CGO_ENABLED = "0";

  buildInputs = [
    pkgs.nixpkgs-fmt
    pkgs.hivemind
    pkgs.gomod2nix
    pkgs.protobuf

    pkgs.golangci-lint

    # File system watchers
    pkgs.reflex
    pkgs.entr

    # Docs
    pkgs.mdbook
  ];

  inherit STATE_DIR TRUSTIX_RPC TRUSTIX_ROOT;

}
