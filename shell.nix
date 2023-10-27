let
  pkgs = import ./pkgs.nix { };
  inherit (pkgs) lib;

in
pkgs.mkShell {

  buildInputs = [
    # Procfile process runner
    pkgs.hivemind

    # Nix go modules code generator
    pkgs.gomod2nix

    # Protobuf
    pkgs.protobuf
    pkgs.grpcurl # gRPC CLI

    # Go linters
    pkgs.golangci-lint # Multi purpose linter

    # File system watchers
    pkgs.reflex
    pkgs.entr

    # Docs
    pkgs.mdbook

    # License management and compliance
    pkgs.reuse

    # Socket activation testing
    pkgs.systemfd

    # Dev
    pkgs.go
    pkgs.nix-eval-jobs
    pkgs.sqlite
    pkgs.diffoscope
    pkgs.sqlc
    pkgs.goose
    pkgs.protoc-gen-go
    pkgs.protoc-gen-doc
    pkgs.protoc-gen-connect-go
    pkgs.nodejs
  ];

  # Write token used for log submission
  TRUSTIX_TOKEN = "${builtins.toString ./packages/trustix/dev/token-priv}";

  FLAKE_ROOT = builtins.toString ./.;

  shellHook = ''
    export TRUSTIX_RPC="unix://$FLAKE_ROOT/state/trustix.sock"
    export TRUSTIX_NIX_REPROD_STATE_DIR="$STATE_DIR/nix-reprod"
    export PATH=${builtins.toString ./packages/trustix-nix-r13y-web}/node_modules/.bin:$PATH

    export TRUSTIX_STATE_DIR="$STATE_DIR/trustix";
  '';
}
