let
  pkgs = import ./pkgs.nix { };

  STATE_DIR = "${builtins.toString ./.}/state";
  TRUSTIX_RPC = "unix://${STATE_DIR}/trustix.sock";
  TRUSTIX_ROOT = builtins.toString ./.;

  # Wrap treefmt with a Go compiler so it can do gofmt without recursively loading subprojects
  treefmt = pkgs.writeShellScriptBin "treefmt" ''
    export PATH=${pkgs.go}/bin:$PATH
    exec ${pkgs.treefmt}/bin/treefmt "$@"
  '';

in
pkgs.mkShell {

  # Speed up compilation, guarantee static linking
  CGO_ENABLED = "0";

  buildInputs = [
    # Meta code formatter
    treefmt

    # Only build job if it's not in the binary cache
    pkgs.nix-build-uncached

    # Protobuf formatter (clang-format)
    pkgs.clang-tools

    # Format Nix expressions
    pkgs.nixpkgs-fmt

    # Procfile process runner
    pkgs.hivemind

    # Nix go modules code generator
    pkgs.gomod2nix

    # Protobuf
    pkgs.protobuf
    pkgs.grpcurl # gRPC CLI

    # Go linter
    pkgs.golangci-lint

    # File system watchers
    pkgs.reflex
    pkgs.entr

    # Docs
    pkgs.mdbook

    # License
    pkgs.reuse
  ];

  inherit STATE_DIR TRUSTIX_RPC TRUSTIX_ROOT;

}
