let
  pkgs = import ./pkgs.nix { };
  inherit (pkgs) lib;

  STATE_DIR = "${builtins.toString ./.}/state";
  TRUSTIX_RPC = "unix://${STATE_DIR}/trustix.sock";
  TRUSTIX_ROOT = builtins.toString ./.;

  python = pkgs.python3.override {
    self = python;
    packageOverrides = self: super: {

      pretty-errors =
        let
          version = "1.2.25";
        in
        self.buildPythonPackage {
          pname = "pretty-errors";
          inherit version;

          src = self.fetchPypi {
            pname = "pretty_errors";
            inherit version;
            hash = "sha256-oWulx1LIfCY7+S+LS1hiTjseKScak5H1ZPErhuk8Z1U=";
          };

          # Work around interactive installer
          postPatch = "rm ./pretty_errors/__main__.py";

          propagatedBuildInputs = [
            self.colorama
          ];
        };

    };
  };

  # Some development tools (like the license file generator) is written in Python
  pythonEnv = python.withPackages (ps: [
    ps.mypy
    ps.black
    ps.pretty-errors
  ]);

in
pkgs.mkShell {

  buildInputs = [
    # Development scripts
    pythonEnv

    # Only build job if it's not in the binary cache
    pkgs.nix-build-uncached

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

  inherit STATE_DIR TRUSTIX_RPC TRUSTIX_ROOT;

  # Write token used for log submission
  TRUSTIX_TOKEN = "${builtins.toString ./packages/trustix/dev/token-priv}";


  shellHook = ''
    export TRUSTIX_NIX_REPROD_STATE_DIR="$STATE_DIR/nix-reprod"
    export PATH=${builtins.toString ./packages/trustix-nix-r13y-web}/node_modules/.bin:$PATH
    export TRUSTIX_STATE_DIR="$STATE_DIR/trustix";
  '';
}
