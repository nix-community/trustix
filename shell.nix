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
    # Format Nix expressions
    pkgs.nixpkgs-fmt

    # Procfile process runner
    pkgs.hivemind

    # Nix go modules code generator
    pkgs.gomod2nix

    # Protobuf
    pkgs.protobuf
    pkgs.grpcurl  # gRPC CLI

    pkgs.golangci-lint

    # File system watchers
    pkgs.reflex
    pkgs.entr

    # Docs
    pkgs.mdbook
  ];

  inherit STATE_DIR TRUSTIX_RPC TRUSTIX_ROOT;

}
